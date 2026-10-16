package client_postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	pool *pgxpool.Pool
}

func (x Client) Exec(ctx context.Context, query string, arguments ...any) (pgconn.CommandTag, error) {
	return x.pickPool(ctx).Exec(ctx, query, arguments...)
}

func (x Client) InTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if err := pgx.BeginFunc(ctx, x.pool, func(transaction pgx.Tx) error {
		return fn(context.WithValue(ctx, transactionKey{}, transaction))
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func (x Client) Query(ctx context.Context, query string, arguments ...any) (pgx.Rows, error) {
	return x.pickPool(ctx).Query(ctx, query, arguments...)
}

func (x Client) QueryRow(ctx context.Context, query string, arguments ...any) pgx.Row {
	return x.pickPool(ctx).QueryRow(ctx, query, arguments...)
}

type sqlClient interface {
	Exec(ctx context.Context, query string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, query string, arguments ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, arguments ...any) pgx.Row
}

type transactionKey struct{}

func (x Client) pickPool(ctx context.Context) sqlClient {
	if transaction, ok := ctx.Value(transactionKey{}).(pgx.Tx); ok {
		return transaction
	}

	return x.pool
}

func NewClient(pool *pgxpool.Pool) Client {
	return Client{pool: pool}
}
