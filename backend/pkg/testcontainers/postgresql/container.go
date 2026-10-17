package pkg_testcontainers_postgresql

import (
	"context"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	image               = "postgres:17"
	readyLogOccurrences = 2
)

func NewPostgreSQL(t *testing.T) *pgxpool.Pool {
	t.Helper()

	ctx := t.Context()
	container, err := testcontainers.Run(ctx, image,
		testcontainers.WithEnv(map[string]string{
			"POSTGRES_DB":       "test",
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
		}),
		testcontainers.WithExposedPorts("5432/tcp"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(readyLogOccurrences),
		),
	)
	require.NoError(t, err)
	testcontainers.CleanupContainer(t, container)

	host, err := container.Host(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err)

	dsn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("test", "test"),
		Host:     net.JoinHostPort(host, port.Port()),
		Path:     "test",
		RawQuery: "sslmode=disable",
	}
	pool, err := pgxpool.New(ctx, dsn.String())
	require.NoError(t, err)
	t.Cleanup(pool.Close)
	require.NoError(t, pool.Ping(ctx))
	migrate(t, ctx, pool)

	return pool
}

func migrate(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()

	_, source, _, ok := runtime.Caller(0)
	require.True(t, ok)

	migrations := filepath.Join(filepath.Dir(source), "..", "..", "..", "migrations")
	database := stdlib.OpenDBFromPool(pool)

	defer func() {
		require.NoError(t, database.Close())
	}()

	provider, err := goose.NewProvider(goose.DialectPostgres, database, os.DirFS(migrations))
	require.NoError(t, err)
	_, err = provider.Up(ctx)
	require.NoError(t, err)
}
