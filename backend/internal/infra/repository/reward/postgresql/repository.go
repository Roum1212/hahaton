package repository_reward_postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	aggregate_reward "github.com/Roum1212/hahaton/internal/domain/aggregate/reward"
	model_reward_id "github.com/Roum1212/hahaton/internal/domain/model/reward-id"
	repository_locking "github.com/Roum1212/hahaton/internal/infra/repository/locking"
)

const (
	uniqueViolationCode = "23505"

	idColumn            = "id"
	coordinatorIDColumn = "coordinator_id"

	deletedColumn          = "deleted"
	descriptionColumn      = "description"
	imagePathColumn        = "image_path"
	issuedQuantityColumn   = "issued_quantity"
	priceCoinsColumn       = "price_coins"
	quantityLimitColumn    = "quantity_limit"
	reservedQuantityColumn = "reserved_quantity"
	titleColumn            = "title"

	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type Rewards struct {
	database Database
	queries  squirrel.StatementBuilderType
}

func (x Rewards) Create(ctx context.Context, value *aggregate_reward.Reward) error {
	_, err := x.execute(ctx, x.queries.Insert("rewards").SetMap(aggregateValues(value)))

	return err
}

func (x Rewards) List(ctx context.Context) ([]*aggregate_reward.Reward, error) {
	sql, arguments, err := x.queries.Select(
		idColumn,
		coordinatorIDColumn,
		titleColumn,
		descriptionColumn,
		imagePathColumn,
		priceCoinsColumn,
		quantityLimitColumn,
		reservedQuantityColumn,
		issuedQuantityColumn,
		deletedColumn,
		createdAtColumn,
		updatedAtColumn,
	).From("rewards").ToSql()
	if err != nil {
		return nil, fmt.Errorf("build reward list: %w", err)
	}

	rows, err := x.database.Query(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := make([]*aggregate_reward.Reward, 0)

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

func (x Rewards) Load(ctx context.Context, rewardID model_reward_id.RewardID) (*aggregate_reward.Reward, error) {
	return x.load(ctx, rewardID, false)
}

func (x Rewards) Modify(
	ctx context.Context,
	rewardID model_reward_id.RewardID,
	change func(
		ctx context.Context,
		value *aggregate_reward.Reward,
	) error,
) error {
	return x.database.InTx(ctx, func(ctx context.Context) error {
		value, err := x.Load(ctx, rewardID)
		if err != nil {
			return err
		}

		if err := repository_locking.Coordinator(ctx, x.database, value.GetCoordinatorID(), false); err != nil {
			return err
		}

		value, err = x.load(ctx, rewardID, true)
		if err != nil {
			return err
		}

		if err := change(ctx, value); err != nil {
			return err
		}

		return x.Save(ctx, value)
	})
}

func (x Rewards) Save(ctx context.Context, value *aggregate_reward.Reward) error {
	values := aggregateValues(value)
	delete(values, idColumn)

	result, err := x.execute(
		ctx,
		x.queries.Update("rewards").SetMap(values).Where(squirrel.Eq{idColumn: value.GetRewardID().String()}),
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return aggregate_reward.ErrNotFound
	}

	return nil
}

func aggregateValues(value *aggregate_reward.Reward) map[string]any {
	imagePath := value.GetImagePath().String()

	return map[string]any{
		idColumn:               value.GetRewardID().String(),
		coordinatorIDColumn:    value.GetCoordinatorID().String(),
		titleColumn:            value.GetTitle(),
		descriptionColumn:      value.GetDescription(),
		imagePathColumn:        imagePath,
		priceCoinsColumn:       value.GetPriceCoins(),
		quantityLimitColumn:    value.GetQuantityLimit(),
		reservedQuantityColumn: value.GetReservedQuantity(),
		issuedQuantityColumn:   value.GetIssuedQuantity(),
		deletedColumn:          value.GetDeleted(),
		createdAtColumn:        value.GetCreatedAt(),
		updatedAtColumn:        value.GetUpdatedAt(),
	}
}

func (x Rewards) execute(ctx context.Context, query squirrel.Sqlizer) (pgconn.CommandTag, error) {
	sql, arguments, err := query.ToSql()
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("build reward query: %w", err)
	}

	result, err := x.database.Exec(ctx, sql, arguments...)
	if err != nil {
		var pgError *pgconn.PgError
		switch {
		case errors.As(err, &pgError) && pgError.Code == uniqueViolationCode:
			return pgconn.CommandTag{}, aggregate_reward.ErrAlreadyExists
		default:
			return pgconn.CommandTag{}, fmt.Errorf("execute reward query: %w", err)
		}
	}

	return result, nil
}

func (x Rewards) load(
	ctx context.Context,
	rewardID model_reward_id.RewardID,
	lock bool,
) (*aggregate_reward.Reward, error) {
	query := x.queries.Select(
		idColumn,
		coordinatorIDColumn,
		titleColumn,
		descriptionColumn,
		imagePathColumn,
		priceCoinsColumn,
		quantityLimitColumn,
		reservedQuantityColumn,
		issuedQuantityColumn,
		deletedColumn,
		createdAtColumn,
		updatedAtColumn,
	).From("rewards").Where(squirrel.Eq{idColumn: rewardID.String()})
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
			return nil, aggregate_reward.ErrNotFound
		default:
			return nil, err
		}
	}

	return record.toAggregate()
}

func NewRewards(database Database) Rewards {
	return Rewards{
		database: database,
		queries:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
