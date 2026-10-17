package repository_cooperation_postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

const (
	uniqueViolationCode = "23505"

	coordinatorIDColumn = "coordinator_id"
	volunteerIDColumn   = "volunteer_id"

	actorIDColumn = "actor_id"
	statusColumn  = "status"

	createdAtColumn  = "created_at"
	updatedAtColumn  = "updated_at"
	retryAfterColumn = "retry_after"
)

type Cooperations struct {
	database Database
	queries  squirrel.StatementBuilderType
}

func (x Cooperations) Create(
	ctx context.Context,
	value *aggregate_cooperation.Cooperation,
	afterCreate func(ctx context.Context) error,
) error {
	return x.database.InTx(ctx, func(ctx context.Context) error {
		if _, err := x.execute(ctx, x.queries.Insert("cooperations").SetMap(aggregateValues(value))); err != nil {
			return err
		}

		return afterCreate(ctx)
	})
}

func (x Cooperations) List(ctx context.Context) ([]*aggregate_cooperation.Cooperation, error) {
	sql, arguments, err := x.queries.Select(
		coordinatorIDColumn,
		volunteerIDColumn,
		statusColumn,
		retryAfterColumn,
		createdAtColumn,
		updatedAtColumn,
	).From("cooperations").ToSql()
	if err != nil {
		return nil, fmt.Errorf("build cooperation list: %w", err)
	}

	rows, err := x.database.Query(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := make([]*aggregate_cooperation.Cooperation, 0)

	for rows.Next() {
		var record row
		if err := record.scan(rows); err != nil {
			return nil, err
		}

		value, err := record.toAggregate()
		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return values, nil
}

func (x Cooperations) Load(
	ctx context.Context,
	coordinatorID model_user_id.UserID,
	volunteerID model_user_id.UserID,
) (*aggregate_cooperation.Cooperation, error) {
	sql, arguments, err := x.queries.Select(
		coordinatorIDColumn,
		volunteerIDColumn,
		statusColumn,
		retryAfterColumn,
		createdAtColumn,
		updatedAtColumn,
	).From("cooperations").Where(squirrel.Eq{
		coordinatorIDColumn: coordinatorID.String(),
		volunteerIDColumn:   volunteerID.String(),
	}).ToSql()
	if err != nil {
		return nil, err
	}

	var record row
	if err := record.scan(x.database.QueryRow(ctx, sql, arguments...)); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, aggregate_cooperation.ErrNotFound
		default:
			return nil, err
		}
	}

	return record.toAggregate()
}

func (x Cooperations) Modify(
	ctx context.Context,
	coordinatorID model_user_id.UserID,
	volunteerID model_user_id.UserID,
	change func(
		ctx context.Context,
		value *aggregate_cooperation.Cooperation,
	) error,
) error {
	return x.database.InTx(ctx, func(ctx context.Context) error {
		value, err := x.Load(ctx, coordinatorID, volunteerID)
		if err != nil {
			return err
		}

		if err := change(ctx, value); err != nil {
			return err
		}

		return x.Save(ctx, value)
	})
}

func (x Cooperations) Save(ctx context.Context, value *aggregate_cooperation.Cooperation) error {
	values := aggregateValues(value)
	delete(values, coordinatorIDColumn)
	delete(values, volunteerIDColumn)

	result, err := x.execute(ctx, x.queries.Update("cooperations").SetMap(values).Where(squirrel.Eq{
		coordinatorIDColumn: value.GetCoordinatorID().String(),
		volunteerIDColumn:   value.GetVolunteerID().String(),
	}))
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return aggregate_cooperation.ErrNotFound
	}

	return nil
}

func aggregateValues(value *aggregate_cooperation.Cooperation) map[string]any {
	return map[string]any{
		coordinatorIDColumn: value.GetCoordinatorID().String(),
		volunteerIDColumn:   value.GetVolunteerID().String(),
		statusColumn:        value.GetStatus().String(),
		retryAfterColumn:    value.GetRetryAfter(),
		createdAtColumn:     value.GetCreatedAt(),
		updatedAtColumn:     value.GetUpdatedAt(),
	}
}

func (x Cooperations) execute(ctx context.Context, query squirrel.Sqlizer) (pgconn.CommandTag, error) {
	sql, arguments, err := query.ToSql()
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("build cooperation query: %w", err)
	}

	result, err := x.database.Exec(ctx, sql, arguments...)
	if err != nil {
		var pgError *pgconn.PgError
		switch {
		case errors.As(err, &pgError) && pgError.Code == uniqueViolationCode:
			return pgconn.CommandTag{}, aggregate_cooperation.ErrAlreadyExists
		default:
			return pgconn.CommandTag{}, fmt.Errorf("execute cooperation query: %w", err)
		}
	}

	return result, nil
}

func NewCooperations(database Database) Cooperations {
	return Cooperations{
		database: database,
		queries:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
