package repository_session_postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	aggregate_session "github.com/Roum1212/hahaton/internal/domain/aggregate/session"
	model_session_id "github.com/Roum1212/hahaton/internal/domain/model/session-id"
)

type Sessions struct {
	database Database
	queries  squirrel.StatementBuilderType
}

func (x Sessions) Create(ctx context.Context, session *aggregate_session.Session) error {
	query := x.queries.Insert("sessions").SetMap(sessionValues(session))

	_, err := execute(ctx, x.database, query)

	return err
}

func (x Sessions) Delete(ctx context.Context, sessionID model_session_id.SessionID) error {
	query := x.queries.Delete("sessions").Where(squirrel.Eq{"id": sessionID.String()})
	_, err := execute(ctx, x.database, query)

	return err
}

func (x Sessions) Load(ctx context.Context, sessionID model_session_id.SessionID) (*aggregate_session.Session, error) {
	query := x.queries.Select(
		"COALESCE(user_id::text, '')", "created_at", "updated_at", "expires_at",
	).From("sessions").Where(squirrel.Eq{"id": sessionID.String()})

	sql, arguments, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build session query: %w", err)
	}

	var row sessionRow

	if err := x.database.QueryRow(ctx, sql, arguments...).Scan(
		&row.userID, &row.createdAt, &row.updatedAt, &row.expiresAt,
	); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, aggregate_session.ErrNotFound
		default:
			return nil, fmt.Errorf("load session: %w", err)
		}
	}

	return row.toSession(sessionID)
}

func (x Sessions) Save(ctx context.Context, session *aggregate_session.Session) error {
	values := sessionValues(session)
	delete(values, "id")

	query := x.queries.Update("sessions").SetMap(values).Where(squirrel.Eq{"id": session.GetSessionID().String()})

	result, err := execute(ctx, x.database, query)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return aggregate_session.ErrNotFound
	}

	return nil
}

func sessionValues(session *aggregate_session.Session) map[string]any {
	values := map[string]any{
		"id":         session.GetSessionID().String(),
		"user_id":    nil,
		"created_at": session.GetCreatedAt(),
		"updated_at": session.GetUpdatedAt(),
		"expires_at": session.GetExpiresAt(),
	}

	if !session.GetUserID().IsZero() {
		values["user_id"] = session.GetUserID().String()
	}

	return values
}

func NewSessions(database Database) Sessions {
	return Sessions{
		database: database,
		queries:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
