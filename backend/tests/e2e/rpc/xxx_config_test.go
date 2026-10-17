//go:build e2e

package rpc

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	pkg_testcontainers_postgresql "github.com/Roum1212/hahaton/pkg/testcontainers/postgresql"
)

const (
	adminAPIKey    = "e2e-admin"
	clientID       = "e2e"
	requestTimeout = 10 * time.Second
	startupTimeout = 15 * time.Second
	pollInterval   = 20 * time.Millisecond
)

type RPCTestSuite struct {
	suite.Suite

	pool       *pgxpool.Pool
	provider   *identityProvider
	httpClient *http.Client
	clients    rpcClients
	baseURL    string
	adminURL   string
	coverage   *rpcCoverage
}

func (x *RPCTestSuite) SetupSuite() {
	t := x.T()
	x.pool = pkg_testcontainers_postgresql.NewPostgreSQL(t)
	x.provider = newIdentityProvider(t)
	directory := t.TempDir()
	certificate := filepath.Join(directory, "certificate.pem")
	key := filepath.Join(directory, "key.pem")

	x.Require().NoError(os.WriteFile(certificate, pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: x.provider.server.Certificate().Raw,
	}), 0o600))
	privateKey, err := x509.MarshalPKCS8PrivateKey(x.provider.server.TLS.Certificates[0].PrivateKey)
	x.Require().NoError(err)
	x.Require().NoError(
		os.WriteFile(key, pem.EncodeToMemory(&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: privateKey,
		}), 0o600),
	)
	x.baseURL = "https://" + freeAddress(t)
	x.adminURL = "http://" + freeAddress(t)

	_, source, _, ok := runtime.Caller(0)
	x.Require().True(ok)

	root := filepath.Clean(filepath.Join(filepath.Dir(source), "..", "..", ".."))
	binary, err := os.Executable()
	x.Require().NoError(err)

	csvPath := filepath.Join(directory, "cities.csv")
	x.Require().NoError(os.WriteFile(csvPath, []byte(
		"fias_id,city,region,geo_lat,geo_lon\nobninsk,Обнинск,Калужская область,55.1,36.6\n",
	), 0o600))

	environment := map[string]string{
		"E2E_BACKEND_PROCESS":            "1",
		"POSTGRESQL_DSN":                 x.pool.Config().ConnString(),
		"HTTP_SERVER_ADDRESS":            strings.TrimPrefix(x.baseURL, "https://"),
		"ADMIN_HTTP_SERVER_ADDRESS":      strings.TrimPrefix(x.adminURL, "http://"),
		"HTTP_SERVER_TLS_CERTIFICATE":    certificate,
		"HTTP_SERVER_TLS_KEY":            key,
		"SSL_CERT_FILE":                  certificate,
		"ADMIN_API_KEY":                  adminAPIKey,
		"FRONTEND_URL":                   "https://frontend.test",
		"FLIDA_ISSUER":                   x.provider.server.URL,
		"FLIDA_CLIENT_ID":                clientID,
		"SESSION_TOKEN_TTL":              "24h",
		"SESSION_TOKEN_V4_SYMMETRIC_KEY": strings.Repeat("01", 32),
		"DATA_ENCRYPTION_KEY":            strings.Repeat("02", 32),
		"STORAGE_DIRECTORY":              filepath.Join(directory, "storage"),
		"CITY_CATALOG_CSV_PATH":          csvPath,
		"CITY_CATALOG_BOUNDARIES_PATH":   "",
		"CITY_CATALOG_VERSION":           "e2e",
		"DADATA_API_KEY":                 "e2e-dadata",
		"DADATA_ENDPOINT":                x.provider.server.URL + "/addresses",
	}
	//nolint:gosec // os.Executable resolves the currently running test binary.
	command := exec.CommandContext(context.WithoutCancel(t.Context()), binary, "-test.run=^TestBackendProcess$")
	command.Env = appendEnvironment(os.Environ(), environment)
	command.Dir = root

	var logs bytes.Buffer

	command.Stdout, command.Stderr = &logs, &logs
	x.Require().NoError(command.Start())

	done := make(chan error, 1)

	go func() {
		done <- command.Wait()
	}()

	t.Cleanup(func() {
		select {
		case err := <-done:
			require.NoError(t, err, logs.String())

			return
		default:
		}

		require.NoError(t, command.Process.Signal(os.Interrupt))

		select {
		case err := <-done:
			require.NoError(t, err, logs.String())
		case <-time.After(startupTimeout):
			require.NoError(t, command.Process.Kill())
			<-done
			t.Error("backend did not stop in time")
		}
	})

	x.coverage = newRPCCoverage()
	transport := x.provider.server.Client().Transport
	x.httpClient = &http.Client{
		Transport: &observedTransport{
			next:     transport,
			coverage: x.coverage,
		},
		Timeout: requestTimeout,
	}
	x.clients = newRPCClients(x.httpClient, x.baseURL, x.adminURL)
	x.Require().Eventually(func() bool {
		return healthy(x.httpClient, x.baseURL)
	}, startupTimeout, pollInterval)
}

func (x *RPCTestSuite) TearDownSuite() {
	if strings.Contains(flag.Lookup("test.run").Value.String(), "/") ||
		flag.Lookup("testify.m").Value.String() != "" {
		return
	}

	x.coverage.assertComplete(x.T())
}

func TestRPC(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(RPCTestSuite))
}

func appendEnvironment(base []string, values map[string]string) []string {
	environment := make([]string, 0, len(base)+len(values))
	for _, entry := range base {
		name, _, _ := strings.Cut(entry, "=")
		if _, replaced := values[name]; !replaced {
			environment = append(environment, entry)
		}
	}

	for name, value := range values {
		environment = append(environment, name+"="+value)
	}

	return environment
}

func freeAddress(t *testing.T) string {
	t.Helper()

	var config net.ListenConfig

	listener, err := config.Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)

	address := listener.Addr().String()
	require.NoError(t, listener.Close())

	return address
}

func healthy(client *http.Client, address string) bool {
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, address+"/healthz", nil)
	if err != nil {
		return false
	}

	response, err := client.Do(request)
	if err != nil {
		return false
	}

	return response.Body.Close() == nil && response.StatusCode == http.StatusNoContent
}
