package repository_session_postgresql

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Database interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
	QueryRow(ctx context.Context, sql string, arguments ...any) pgx.Row
}

func execute(ctx context.Context, database Database, query squirrel.Sqlizer) (pgconn.CommandTag, error) {
	sql, arguments, err := query.ToSql()
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("build sql: %w", err)
	}

	result, err := database.Exec(ctx, sql, arguments...)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("execute sql: %w", err)
	}

	return result, nil
}
