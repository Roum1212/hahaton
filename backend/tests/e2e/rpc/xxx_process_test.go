//go:build e2e

package rpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Roum1212/hahaton/cmd/hahaton"
)

func TestBackendProcess(t *testing.T) {
	t.Parallel()

	if os.Getenv("E2E_BACKEND_PROCESS") != "1" {
		t.Skip("started by the RPC suite as a child process")
	}

	certificate, err := os.ReadFile(os.Getenv("SSL_CERT_FILE"))
	require.NoError(t, err)

	roots := x509.NewCertPool()
	require.True(t, roots.AppendCertsFromPEM(certificate))
	// The child trusts only the local provider; TLS verification stays enabled.
	http.DefaultTransport = &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:    roots,
			MinVersion: tls.VersionTLS12,
		},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	require.NoError(t, hahaton.Run(ctx))
}
