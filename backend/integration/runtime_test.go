//go:build integration

package integration

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v10"

	"github.com/Roum1212/hahaton/cmd/hahaton"
	"github.com/Roum1212/hahaton/internal/config"
	client_flida "github.com/Roum1212/hahaton/internal/infra/client/flida"
)

func checkCORS(x *fixture) {
	origins := []string{
		"http://localhost:3000",
		"https://another-frontend.example.com",
	}

	handlers := []struct {
		handler http.Handler
		method  string
	}{
		{
			x.handler,
			"volunteer.v1.VolunteerService/GetVolunteer",
		},
		{
			x.handler,
			"admin.v1.AdminService/FindVolunteers",
		},
	}
	for _, origin := range origins {
		for _, target := range handlers {
			request := httptest.NewRequest(http.MethodOptions, "https://demo.test/"+target.method, nil)
			request.Header.Set("Origin", origin)
			request.Header.Set("Access-Control-Request-Method", http.MethodPost)
			request.Header.Set(
				"Access-Control-Request-Headers",
				"authorization,connect-protocol-version,content-type,grpc-timeout",
			)

			response := httptest.NewRecorder()
			target.handler.ServeHTTP(response, request)
			require(x.t, response.Code == http.StatusNoContent, "preflight requires authentication")
			require(
				x.t,
				response.Header().Get("Access-Control-Allow-Origin") == origin,
				"cors rejected arbitrary origin",
			)
			require(x.t, response.Header().Get("Access-Control-Allow-Credentials") == "true", "cors credentials missing")

			allowed := strings.ToLower(strings.Join(response.Header().Values("Access-Control-Allow-Headers"), ","))
			for _, header := range []string{
				"content-type",
				"connect-protocol-version",
				"grpc-timeout",
				"authorization",
			} {
				require(x.t, strings.Contains(allowed, header), "cors preflight omitted "+header)
			}

			request = httptest.NewRequest(http.MethodPost, "https://demo.test/"+target.method, strings.NewReader("{}"))
			request.Header.Set("Origin", origin)
			request.Header.Set("Content-Type", "application/json")

			response = httptest.NewRecorder()
			target.handler.ServeHTTP(response, request)
			require(x.t, response.Code == http.StatusUnauthorized, "cors bypassed authentication")
			require(x.t, response.Header().Get("Access-Control-Allow-Origin") == origin, "cors headers missing on error")
			exposed := strings.ToLower(strings.Join(response.Header().Values("Access-Control-Expose-Headers"), ","))
			require(x.t, strings.Contains(exposed, "grpc-status-details-bin"), "grpc error detail header not exposed")
		}

		request := httptest.NewRequest(
			http.MethodPost,
			"https://demo.test/session.v1.SessionService/OpenSession",
			strings.NewReader("{}"),
		)
		request.Header.Set("Origin", origin)
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()
		x.handler.ServeHTTP(response, request)
		require(x.t, response.Code == http.StatusOK, "cross-origin session could not be opened")
		require(x.t, response.Header().Get("Access-Control-Allow-Origin") == origin, "cors headers missing on success")
		cookies := response.Result().Cookies()
		require(x.t, len(cookies) == 1, "session cookie missing")
		require(x.t, cookies[0].Secure && cookies[0].HttpOnly && cookies[0].SameSite == http.SameSiteNoneMode,
			"session cookie does not support cross-site requests")

		request = httptest.NewRequest(
			http.MethodPost,
			"https://demo.test/session.v1.SessionService/CloseSession",
			strings.NewReader("{}"),
		)
		request.Header.Set("Origin", origin)
		request.Header.Set("Content-Type", "application/json")
		request.AddCookie(cookies[0])

		response = httptest.NewRecorder()
		x.handler.ServeHTTP(response, request)
		require(x.t, response.Code == http.StatusOK, "cross-origin session could not be closed")
		cookies = response.Result().Cookies()
		require(x.t, len(cookies) == 1 && cookies[0].MaxAge < 0 && cookies[0].SameSite == http.SameSiteNoneMode,
			"cross-site session cookie was not removed")
	}
}

func checkRuntime(x *fixture) {
	checkCORS(x)

	providerClient := x.provider.Client()
	providerClient.Timeout = time.Second
	_, err := client_flida.NewClient(
		x.t.Context(), x.provider.URL, "demo", x.provider.URL+"/invalid-issuer-discovery", providerClient,
	)
	require(
		x.t,
		err != nil && strings.Contains(err.Error(), "issuer does not match"),
		"discovery issuer mismatch accepted",
	)

	var cfg config.Config

	originalKey := os.Getenv("SESSION_TOKEN_V4_SYMMETRIC_KEY")

	x.t.Setenv("SESSION_TOKEN_V4_SYMMETRIC_KEY", "invalid")

	err = env.ParseWithOptions(&cfg, config.ParseOptions())
	require(x.t, err != nil, "invalid session encryption key was accepted")
	x.t.Setenv("SESSION_TOKEN_V4_SYMMETRIC_KEY", originalKey)
	x.t.Setenv("SESSION_TOKEN_TTL", "0s")

	must(x.t, env.ParseWithOptions(&cfg, config.ParseOptions()))
	require(x.t, cfg.Validate() != nil, "zero session ttl was accepted")
	x.t.Setenv("SESSION_TOKEN_TTL", "24h")
	clientAddress := unusedAddress(x)
	x.t.Setenv("HTTP_SERVER_ADDRESS", clientAddress)

	databaseURL, err := url.Parse(os.Getenv("POSTGRESQL_DSN"))
	must(x.t, err)

	values := databaseURL.Query()
	values.Set("search_path", x.pool.Config().ConnConfig.RuntimeParams["search_path"])
	databaseURL.RawQuery = values.Encode()
	x.t.Setenv("POSTGRESQL_DSN", databaseURL.String())

	ctx, cancel := context.WithCancel(x.t.Context())
	defer cancel()

	result := make(chan error, 1)

	go func() {
		result <- hahaton.Run(ctx)
	}()

	client := &http.Client{Timeout: time.Second}
	ready := false

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+clientAddress+"/healthz", nil)
		must(x.t, err)

		response, err := client.Do(request)
		if err == nil {
			must(x.t, response.Body.Close())

			ready = response.StatusCode == http.StatusNoContent
			if ready {
				break
			}
		}

		select {
		case err := <-result:
			x.t.Fatalf("server stopped before readiness: %v", err)
		case <-time.After(20 * time.Millisecond):
		}
	}

	require(x.t, ready, "server did not become ready")

	request, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"http://"+clientAddress+"/admin.v1.AdminService/FindVolunteers", strings.NewReader(
			`{"email_address":"missing@example.com"}`,
		))
	must(x.t, err)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer integration-admin")
	response, err := client.Do(request)
	must(x.t, err)
	data, err := io.ReadAll(response.Body)
	must(x.t, err)
	must(x.t, response.Body.Close())
	require(x.t, response.StatusCode == http.StatusOK, "admin route on shared listener failed: "+string(data))
	cancel()

	select {
	case err := <-result:
		must(x.t, err)
	case <-time.After(6 * time.Second):
		x.t.Fatal("server did not shut down")
	}
}

func unusedAddress(x *fixture) string {
	var cfg net.ListenConfig

	first, err := cfg.Listen(x.t.Context(), "tcp", "127.0.0.1:0")
	must(x.t, err)

	defer func() {
		must(x.t, first.Close())
	}()

	return first.Addr().String()
}
