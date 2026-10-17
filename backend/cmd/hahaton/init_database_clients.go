package hahaton

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Roum1212/hahaton/internal/config"
	client_postgresql "github.com/Roum1212/hahaton/internal/infra/client/postgresql"
)

type databaseClientSet struct {
	pool       *pgxpool.Pool
	postgresql client_postgresql.Client
}

func initDatabaseClients(ctx context.Context, cfg config.Config) (databaseClientSet, error) {
	pool, err := pgxpool.New(ctx, cfg.PostgreSQL.DSN)
	if err != nil {
		return databaseClientSet{}, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()

		return databaseClientSet{}, err
	}

	return databaseClientSet{
		pool:       pool,
		postgresql: client_postgresql.NewClient(pool),
	}, nil
}
