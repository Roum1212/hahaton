//go:build integration

package integration

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/go-jose/go-jose/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/proto"

	"github.com/Roum1212/hahaton/cmd/hahaton"
	"github.com/Roum1212/hahaton/internal/config"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_event_schedule "github.com/Roum1212/hahaton/internal/domain/model/event-schedule"
	client_postgresql "github.com/Roum1212/hahaton/internal/infra/client/postgresql"
	repository_event "github.com/Roum1212/hahaton/internal/infra/repository/event/postgresql"
)

type fixture struct {
	seen map[string]bool
	t    *testing.T

	pool         *pgxpool.Pool
	handler      http.Handler
	admin        http.Handler
	provider     *httptest.Server
	signer       jose.Signer
	dadataStatus atomic.Int32
	dadataBody   string
}

type actor struct {
	id     string
	cookie *http.Cookie
}

type transport func(r *http.Request) (*http.Response, error)

func (x transport) RoundTrip(r *http.Request) (*http.Response, error) {
	return x(r)
}

func (x *fixture) call(a actor, method string, body map[string]any) map[string]any {
	x.t.Helper()

	return x.expect(a, method, body, http.StatusOK, "")
}

func (x *fixture) expect(a actor, method string, body map[string]any, status int, code string) map[string]any {
	x.t.Helper()

	response := x.request(a, method, body, false)
	if response.Code != status {
		x.t.Fatalf("%s: status %d, expected %d: %s", method, response.Code, status, response.Body.String())
	}

	result := document(x.t, response.Body.Bytes())
	if code != "" && result["code"] != code {
		x.t.Fatalf("%s: expected %s: %s", method, code, response.Body.String())
	}

	if code == "not_found" {
		resourceInfo(x.t, result)
	}

	return result
}

func (x *fixture) login(subject string, overrides ...map[string]any) actor {
	x.t.Helper()

	opened := x.request(actor{}, "session.v1.SessionService/OpenSession", nil, false)
	if opened.Code != http.StatusOK {
		x.t.Fatal(opened.Body.String())
	}

	anonymous := opened.Result().Cookies()[0]
	claims := map[string]any{
		"iss":           x.provider.URL,
		"aud":           "demo",
		"sub":           subject,
		"iat":           time.Now().Unix(),
		"exp":           time.Now().Add(time.Hour).Unix(),
		"given_name":    subject,
		"family_name":   "Демо",
		"middle_name":   "",
		"email":         subject + "@example.com",
		"phone_number":  "+79990000001",
		"date_of_birth": "2001-01-01",
	}

	for _, values := range overrides {
		for key, value := range values {
			claims[key] = value
		}
	}

	response := x.request(actor{
		cookie: anonymous,
	}, "session.v1.SessionService/AuthenticateSession", map[string]any{"id_token": x.token(claims)}, false)
	if response.Code != http.StatusOK {
		x.t.Fatal(response.Body.String())
	}

	a := actor{cookie: response.Result().Cookies()[0]}
	profile := x.call(a, "session.v1.SessionService/GetSession", nil)
	a.id = text(x.t, object(x.t, profile, "session")["user_id"])

	return a
}

func (x *fixture) coordinator(subject string) actor {
	x.t.Helper()
	a := x.login(subject)

	response := x.request(actor{}, "admin.v1.AdminService/CreateCoordinator", map[string]any{
		"volunteer_id":                         a.id,
		"cooperation_request_cooldown_seconds": "60",
	}, true)
	if response.Code != http.StatusOK {
		x.t.Fatal(response.Body.String())
	}

	return a
}

func (x *fixture) cooperate(coordinator, volunteer actor) {
	x.t.Helper()
	x.call(volunteer, "cooperation.v1.CooperationService/RequestCooperation", map[string]any{
		"coordinator_id": coordinator.id,
	})
	x.call(coordinator, "cooperation.v1.CooperationService/ApproveCooperation", map[string]any{
		"volunteer_id": volunteer.id,
	})
}

func (x *fixture) event(coordinator actor, title string, limit int, overrides ...map[string]any) string {
	x.t.Helper()

	start := time.Now().Add(time.Hour).Truncate(time.Second)
	body := map[string]any{
		"title": title,
		"point": map[string]any{
			"latitude":  55.1,
			"longitude": 36.6,
		},
		"starts_at":       start.Format(time.RFC3339Nano),
		"ends_at":         start.Add(time.Hour).Format(time.RFC3339Nano),
		"volunteer_limit": limit,
	}

	for _, values := range overrides {
		for key, value := range values {
			body[key] = value
		}
	}

	response := x.call(coordinator, "event.v1.EventService/CreateEvent", body)

	return text(x.t, response["id"])
}

func (x *fixture) finishTime(id string) {
	x.t.Helper()

	eventID, err := model_event_id.NewEventID(id)
	must(x.t, err)
	events := repository_event.NewEvents(client_postgresql.NewClient(x.pool))
	value, err := events.Load(x.t.Context(), eventID)
	must(x.t, err)

	start := time.Now().Add(-4 * time.Hour).Truncate(time.Second)
	schedule, err := model_event_schedule.NewEventSchedule(start, start.Add(3*time.Hour))
	must(x.t, err)

	value = aggregate_event.NewEventFromFields(
		value.GetEventID(),
		value.GetCoordinatorID(),
		value.GetTitle(),
		value.GetDescription(),
		value.GetImagePath(),
		value.GetLocation(),
		value.GetPoint(),
		schedule,
		value.GetRequirements(),
		value.GetVolunteerLimit(),
		value.GetStatus(),
		value.GetCompletionRequestID(),
		value.GetResultFingerprint(),
		value.GetCreatedAt(),
		value.GetUpdatedAt(),
	)
	must(x.t, events.Save(x.t.Context(), value))
}

func (x *fixture) participate(coordinator, volunteer actor, eventID string) string {
	x.t.Helper()
	response := x.call(volunteer, "participation.v1.ParticipationService/RequestParticipation", map[string]any{
		"event_id": eventID,
	})
	id := text(x.t, response["id"])
	x.call(coordinator, "participation.v1.ParticipationService/ApproveParticipation", map[string]any{"id": id})

	return id
}

func (x *fixture) request(a actor, method string, body map[string]any, admin bool) *httptest.ResponseRecorder {
	x.t.Helper()

	if body == nil {
		body = map[string]any{}
	}

	data, err := json.Marshal(body)
	must(x.t, err)

	r := httptest.NewRequest(http.MethodPost, "https://demo.test/"+method, bytes.NewReader(data))
	r.Header.Set("Content-Type", "application/json")

	if a.cookie != nil {
		r.AddCookie(a.cookie)
	}

	handler := x.handler
	if admin {
		handler = x.admin

		r.Header.Set("Authorization", "Bearer integration-admin")
	}

	if x.seen != nil {
		x.seen[method] = true
	}

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, r)

	return response
}

func (x *fixture) token(claims map[string]any) string {
	x.t.Helper()

	data, err := json.Marshal(claims)
	must(x.t, err)
	signed, err := x.signer.Sign(data)
	must(x.t, err)
	value, err := signed.CompactSerialize()
	must(x.t, err)

	return value
}

func document(t *testing.T, data []byte) map[string]any {
	t.Helper()

	result := make(map[string]any)
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	must(t, decoder.Decode(&result))

	return result
}

func list(t *testing.T, value map[string]any, key string) []any {
	t.Helper()

	if value[key] == nil {
		return nil
	}

	result, ok := value[key].([]any)
	if !ok {
		t.Fatalf("%s is not a list: %v", key, value)
	}

	return result
}

func must(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatal(err)
	}
}

func object(t *testing.T, value map[string]any, key string) map[string]any {
	t.Helper()

	result, ok := value[key].(map[string]any)
	if !ok {
		t.Fatalf("%s is not an object: %v", key, value)
	}

	return result
}

func require(t *testing.T, condition bool, message string) {
	t.Helper()

	if !condition {
		t.Fatal(message)
	}
}

func newFixture(t *testing.T) *fixture {
	t.Helper()

	databaseURL := os.Getenv("INTEGRATION_POSTGRESQL_DSN")
	if databaseURL == "" {
		t.Skip("set INTEGRATION_POSTGRESQL_DSN to run PostgreSQL integration scenarios")
	}

	ctx := t.Context()
	adminPool, err := pgxpool.New(ctx, databaseURL)
	must(t, err)

	schema := fmt.Sprintf("integration_%d", time.Now().UnixNano())
	schema = strings.ToLower(schema)
	schema = pgx.Identifier{schema}.Sanitize()
	_, err = adminPool.Exec(ctx, "CREATE SCHEMA "+schema)
	must(t, err)
	t.Cleanup(func() {
		_, err := adminPool.Exec(context.Background(), "DROP SCHEMA "+schema+" CASCADE")
		must(t, err)
		adminPool.Close()
	})

	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	must(t, err)

	poolConfig.ConnConfig.RuntimeParams["search_path"] = strings.Trim(schema, "\"")
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	must(t, err)
	t.Cleanup(pool.Close)
	must(t, applyMigrations(ctx, pool))
	must(t, applyMigrations(ctx, pool))
	x := &fixture{
		t:    t,
		pool: pool,
	}
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	must(t, err)
	signer, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.RS256,
		Key:       privateKey,
	}, (&jose.SignerOptions{}).WithHeader("kid", "demo"))
	must(t, err)

	x.signer = signer
	x.dadataStatus.Store(http.StatusOK)
	x.dadataBody = `{
  "suggestions": [
    {
      "value": "Обнинск",
      "unrestricted_value": "Калужская обл, г Обнинск",
      "data": {
        "fias_id": "city",
        "city_fias_id": "city",
        "city_with_type": "г Обнинск",
        "region_with_type": "Калужская обл",
        "geo_lat": "55.1",
        "geo_lon": "36.6",
        "qc_geo": "4"
      }
    }
  ]
}`
	x.provider = httptest.NewTLSServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			must(t, json.NewEncoder(w).Encode(map[string]any{
				"issuer":   x.provider.URL,
				"jwks_uri": x.provider.URL + "/keys",
				"id_token_signing_alg_values_supported": []string{
					"RS256",
				},
			}))
		case "/keys":
			must(t, json.NewEncoder(w).Encode(jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{
				Key:       &privateKey.PublicKey,
				KeyID:     "demo",
				Algorithm: "RS256",
				Use:       "sig",
			}}}))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(x.provider.Close)

	original := http.DefaultTransport
	http.DefaultTransport = transport(func(r *http.Request) (*http.Response, error) {
		if r.URL.Host != "suggestions.dadata.ru" {
			return x.provider.Client().Transport.RoundTrip(r)
		}

		require(t, r.Header.Get("Authorization") == "Token demo-dadata", "dadata key missing")

		body := make(map[string]any)
		must(t, json.NewDecoder(r.Body).Decode(&body))
		require(t, body["query"] == "Обнинск", "query was not normalized")
		require(t, object(t, body, "from_bound")["value"] == "city", "from_bound missing")
		require(t, object(t, body, "to_bound")["value"] == "settlement", "to_bound missing")
		locations := list(t, body, "locations")
		require(t, len(locations) == 1, "Russia filter missing")
		require(t, dictionary(t, locations[0])["country_iso_code"] == "RU", "wrong DaData country filter")
		require(t, integer(t, body, "count") == 10, "default DaData count is not 10")

		w := httptest.NewRecorder()
		w.WriteHeader(int(x.dadataStatus.Load()))
		_, err := w.WriteString(x.dadataBody)
		must(t, err)

		return w.Result(), nil
	})

	t.Cleanup(func() {
		http.DefaultTransport = original
	})
	directory := t.TempDir()
	csvPath := filepath.Join(directory, "cities.csv")
	must(t, os.WriteFile(csvPath, []byte(
		"fias_id,city,region,geo_lat,geo_lon\ncity,Обнинск,Калужская область,55.1,36.6\nno-boundary,Без границы,Регион,,\n",
	), 0o600))

	boundaryPath := filepath.Join(directory, "boundaries.json")
	must(t, os.WriteFile(boundaryPath, []byte(
		`{
  "city": [
    [
      [
        [
          36,
          55
        ],
        [
          37,
          55
        ],
        [
          37,
          56
        ],
        [
          36,
          56
        ],
        [
          36,
          55
        ]
      ],
      [
        [
          36.8,
          55.8
        ],
        [
          36.9,
          55.8
        ],
        [
          36.9,
          55.9
        ],
        [
          36.8,
          55.9
        ],
        [
          36.8,
          55.8
        ]
      ]
    ]
  ]
}`,
	), 0o600))

	for name, value := range map[string]string{
		"POSTGRESQL_DSN":                 databaseURL,
		"ADMIN_API_KEY":                  "integration-admin",
		"FRONTEND_URL":                   "https://demo.test",
		"FLIDA_ISSUER":                   x.provider.URL,
		"FLIDA_CLIENT_ID":                "demo",
		"SESSION_TOKEN_V4_SYMMETRIC_KEY": strings.Repeat("01", 32),
		"DATA_ENCRYPTION_KEY":            strings.Repeat("02", 32),
		"STORAGE_DIRECTORY":              directory,
		"CITY_CATALOG_CSV_PATH":          csvPath,
		"CITY_CATALOG_VERSION":           "integration",
		"CITY_CATALOG_BOUNDARIES_PATH":   boundaryPath,
		"DADATA_API_KEY":                 "demo-dadata",
	} {
		t.Setenv(name, value)
	}

	var cfg config.Config
	must(t, env.ParseWithOptions(&cfg, config.ParseOptions()))
	must(t, cfg.Validate())
	x.handler, x.admin, err = hahaton.NewHandlers(ctx, cfg, pool)
	must(t, err)

	return x
}

func applyMigrations(ctx context.Context, pool *pgxpool.Pool) (err error) {
	database := stdlib.OpenDBFromPool(pool)

	defer func() {
		err = errors.Join(err, database.Close())
	}()

	provider, err := goose.NewProvider(goose.DialectPostgres, database, os.DirFS("../migrations"))
	if err != nil {
		return err
	}

	_, err = provider.Up(ctx)

	return err
}

func (x *fixture) fail(a actor, method string, body map[string]any, code string) map[string]any {
	x.t.Helper()

	statuses := map[string]int{
		"invalid_argument":    400,
		"unauthenticated":     401,
		"permission_denied":   403,
		"not_found":           404,
		"already_exists":      409,
		"aborted":             409,
		"failed_precondition": 400,
		"resource_exhausted":  429,
		"unavailable":         503,
	}

	return x.expect(a, method, body, statuses[code], code)
}

func integer(t *testing.T, value map[string]any, key string) int64 {
	t.Helper()

	if value[key] == nil {
		return 0
	}

	result, err := strconv.ParseInt(fmt.Sprint(value[key]), 10, 64)
	must(t, err)

	return result
}

func dictionary(t *testing.T, value any) map[string]any {
	t.Helper()

	result, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("not an object: %v", value)
	}

	return result
}

func text(t *testing.T, value any) string {
	t.Helper()

	result, ok := value.(string)
	if !ok {
		t.Fatalf("not a string: %v", value)
	}

	return result
}

func resourceInfo(t *testing.T, result map[string]any) *errdetails.ResourceInfo {
	t.Helper()

	details := list(t, result, "details")
	require(t, len(details) == 1, "not found must include resource details")
	detail := dictionary(t, details[0])
	require(t, detail["type"] == "google.rpc.ResourceInfo", "unexpected not found detail type")
	data, err := base64.RawStdEncoding.DecodeString(text(t, detail["value"]))
	must(t, err)

	info := &errdetails.ResourceInfo{}
	must(t, proto.Unmarshal(data, info))
	require(t, info.GetResourceType() != "", "resource type missing")
	require(t, info.GetDescription() != "", "resource description missing")

	return info
}
