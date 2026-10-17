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

const (
	idColumn     = "id"
	userIDColumn = "user_id"

	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
	expiresAtColumn = "expires_at"
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
	query := x.queries.Delete("sessions").Where(squirrel.Eq{idColumn: sessionID.String()})
	_, err := execute(ctx, x.database, query)

	return err
}

func (x Sessions) Load(ctx context.Context, sessionID model_session_id.SessionID) (*aggregate_session.Session, error) {
	return x.load(ctx, sessionID, false)
}

func (x Sessions) Replace(
	ctx context.Context,
	sessionID model_session_id.SessionID,
	replace func(
		ctx context.Context,
		current *aggregate_session.Session,
	) (*aggregate_session.Session, error),
) error {
	return x.database.InTx(ctx, func(ctx context.Context) error {
		current, err := x.load(ctx, sessionID, true)
		if err != nil {
			return err
		}

		session, err := replace(ctx, current)
		if err != nil {
			return err
		}

		if err := x.Delete(ctx, sessionID); err != nil {
			return err
		}

		return x.Create(ctx, session)
	})
}

func (x Sessions) Save(ctx context.Context, session *aggregate_session.Session) error {
	values := sessionValues(session)
	delete(values, idColumn)

	query := x.queries.Update("sessions").SetMap(values).Where(squirrel.Eq{idColumn: session.GetSessionID().String()})

	result, err := execute(ctx, x.database, query)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return aggregate_session.ErrNotFound
	}

	return nil
}

func (x Sessions) load(
	ctx context.Context,
	sessionID model_session_id.SessionID,
	lock bool,
) (*aggregate_session.Session, error) {
	query := x.queries.Select(
		"COALESCE("+userIDColumn+"::text, '')",
		createdAtColumn,
		updatedAtColumn,
		expiresAtColumn,
	).From(
		"sessions",
	).Where(
		squirrel.Eq{idColumn: sessionID.String()},
	)

	if lock {
		query = query.Suffix("FOR UPDATE")
	}

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

func sessionValues(session *aggregate_session.Session) map[string]any {
	values := map[string]any{
		idColumn:        session.GetSessionID().String(),
		userIDColumn:    nil,
		createdAtColumn: session.GetCreatedAt(),
		updatedAtColumn: session.GetUpdatedAt(),
		expiresAtColumn: session.GetExpiresAt(),
	}

	if !session.GetUserID().IsZero() {
		values[userIDColumn] = session.GetUserID().String()
	}

	return values
}

func NewSessions(database Database) Sessions {
	return Sessions{
		database: database,
		queries:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
