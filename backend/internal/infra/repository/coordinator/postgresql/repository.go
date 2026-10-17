package repository_coordinator_postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

const (
	uniqueViolationCode = "23505"

	idColumn = "id"

	avatarPathColumn      = "avatar_path"
	cooldownSecondsColumn = "cooldown_seconds"
	descriptionColumn     = "description"

	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type Coordinators struct {
	database Database
	queries  squirrel.StatementBuilderType
}

func (x Coordinators) Create(ctx context.Context, value *aggregate_coordinator.Coordinator) error {
	_, err := x.execute(ctx, x.queries.Insert("coordinators").SetMap(aggregateValues(value)))

	return err
}

func (x Coordinators) List(ctx context.Context) ([]*aggregate_coordinator.Coordinator, error) {
	sql, arguments, err := x.queries.Select(
		idColumn,
		descriptionColumn,
		avatarPathColumn,
		cooldownSecondsColumn,
		createdAtColumn,
		updatedAtColumn,
	).From("coordinators").ToSql()
	if err != nil {
		return nil, fmt.Errorf("build coordinator list: %w", err)
	}

	rows, err := x.database.Query(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := make([]*aggregate_coordinator.Coordinator, 0)

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

func (x Coordinators) Load(
	ctx context.Context,
	userID model_user_id.UserID,
) (*aggregate_coordinator.Coordinator, error) {
	return x.load(ctx, userID, false)
}

func (x Coordinators) Modify(
	ctx context.Context,
	userID model_user_id.UserID,
	change func(
		ctx context.Context,
		value *aggregate_coordinator.Coordinator,
	) error,
) error {
	return x.database.InTx(ctx, func(ctx context.Context) error {
		value, err := x.load(ctx, userID, true)
		if err != nil {
			return err
		}

		if err := change(ctx, value); err != nil {
			return err
		}

		return x.Save(ctx, value)
	})
}

func (x Coordinators) Save(ctx context.Context, value *aggregate_coordinator.Coordinator) error {
	values := aggregateValues(value)
	delete(values, idColumn)

	result, err := x.execute(
		ctx,
		x.queries.Update("coordinators").SetMap(values).Where(squirrel.Eq{idColumn: value.GetUserID().String()}),
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return aggregate_coordinator.ErrNotFound
	}

	return nil
}

func aggregateValues(value *aggregate_coordinator.Coordinator) map[string]any {
	imagePath := value.GetAvatarPath().String()

	return map[string]any{
		idColumn:              value.GetUserID().String(),
		descriptionColumn:     value.GetDescription(),
		avatarPathColumn:      imagePath,
		cooldownSecondsColumn: value.GetCooldownSeconds(),
		createdAtColumn:       value.GetCreatedAt(),
		updatedAtColumn:       value.GetUpdatedAt(),
	}
}

func (x Coordinators) execute(ctx context.Context, query squirrel.Sqlizer) (pgconn.CommandTag, error) {
	sql, arguments, err := query.ToSql()
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("build coordinator query: %w", err)
	}

	result, err := x.database.Exec(ctx, sql, arguments...)
	if err != nil {
		var pgError *pgconn.PgError
		switch {
		case errors.As(err, &pgError) && pgError.Code == uniqueViolationCode:
			return pgconn.CommandTag{}, aggregate_coordinator.ErrAlreadyExists
		default:
			return pgconn.CommandTag{}, fmt.Errorf("execute coordinator query: %w", err)
		}
	}

	return result, nil
}

func (x Coordinators) load(
	ctx context.Context,
	userID model_user_id.UserID,
	lock bool,
) (*aggregate_coordinator.Coordinator, error) {
	query := x.queries.Select(
		idColumn,
		descriptionColumn,
		avatarPathColumn,
		cooldownSecondsColumn,
		createdAtColumn,
		updatedAtColumn,
	).From("coordinators").Where(squirrel.Eq{idColumn: userID.String()})
	if lock {
		query = query.Suffix("FOR UPDATE")
	}

	sql, arguments, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var record row
	if err := record.scan(x.database.QueryRow(ctx, sql, arguments...)); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, aggregate_coordinator.ErrNotFound
		default:
			return nil, err
		}
	}

	return record.toAggregate()
}

func NewCoordinators(database Database) Coordinators {
	return Coordinators{
		database: database,
		queries:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
