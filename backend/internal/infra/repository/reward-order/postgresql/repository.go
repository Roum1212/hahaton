package repository_reward_order_postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	aggregate_reward_order "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order"
	model_reward_order_id "github.com/Roum1212/hahaton/internal/domain/model/reward-order-id"
	repository_locking "github.com/Roum1212/hahaton/internal/infra/repository/locking"
)

const (
	uniqueViolationCode = "23505"

	idColumn            = "id"
	rewardIDColumn      = "reward_id"
	coordinatorIDColumn = "coordinator_id"
	volunteerIDColumn   = "volunteer_id"

	cancellationReasonColumn = "cancellation_reason"
	priceCoinsColumn         = "price_coins"
	requestIDColumn          = "request_id"
	rewardTitleColumn        = "reward_title"
	statusColumn             = "status"

	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type RewardOrders struct {
	database Database
	queries  squirrel.StatementBuilderType
}

func (x RewardOrders) Create(ctx context.Context, value *aggregate_reward_order.RewardOrder) error {
	_, err := x.execute(ctx, x.queries.Insert("reward_orders").SetMap(aggregateValues(value)))

	return err
}

func (x RewardOrders) List(ctx context.Context) ([]*aggregate_reward_order.RewardOrder, error) {
	sql, arguments, err := x.queries.Select(
		idColumn,
		rewardIDColumn,
		coordinatorIDColumn,
		volunteerIDColumn,
		rewardTitleColumn,
		priceCoinsColumn,
		requestIDColumn,
		statusColumn,
		cancellationReasonColumn,
		createdAtColumn,
		updatedAtColumn,
	).From("reward_orders").ToSql()
	if err != nil {
		return nil, fmt.Errorf("build reward-order list: %w", err)
	}

	rows, err := x.database.Query(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := make([]*aggregate_reward_order.RewardOrder, 0)

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

func (x RewardOrders) Load(
	ctx context.Context,
	rewardOrderID model_reward_order_id.RewardOrderID,
) (*aggregate_reward_order.RewardOrder, error) {
	return x.load(ctx, rewardOrderID, false)
}

func (x RewardOrders) Modify(
	ctx context.Context,
	rewardOrderID model_reward_order_id.RewardOrderID,
	change func(
		ctx context.Context,
		value *aggregate_reward_order.RewardOrder,
	) error,
) error {
	return x.database.InTx(ctx, func(ctx context.Context) error {
		value, err := x.Load(ctx, rewardOrderID)
		if err != nil {
			return err
		}

		if err := repository_locking.Coordinator(ctx, x.database, value.GetCoordinatorID(), false); err != nil {
			return err
		}

		if err := repository_locking.Reward(ctx, x.database, value.GetRewardID()); err != nil {
			return err
		}

		value, err = x.load(ctx, rewardOrderID, true)
		if err != nil {
			return err
		}

		if err := change(ctx, value); err != nil {
			return err
		}

		return x.Save(ctx, value)
	})
}

func (x RewardOrders) Save(ctx context.Context, value *aggregate_reward_order.RewardOrder) error {
	values := aggregateValues(value)
	delete(values, idColumn)

	result, err := x.execute(
		ctx,
		x.queries.Update("reward_orders").SetMap(values).Where(
			squirrel.Eq{idColumn: value.GetRewardOrderID().String()},
		),
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return aggregate_reward_order.ErrNotFound
	}

	return nil
}

func aggregateValues(value *aggregate_reward_order.RewardOrder) map[string]any {
	return map[string]any{
		idColumn:                 value.GetRewardOrderID().String(),
		rewardIDColumn:           value.GetRewardID().String(),
		coordinatorIDColumn:      value.GetCoordinatorID().String(),
		volunteerIDColumn:        value.GetVolunteerID().String(),
		rewardTitleColumn:        value.GetRewardTitle(),
		priceCoinsColumn:         value.GetPriceCoins(),
		requestIDColumn:          value.GetRequestID(),
		statusColumn:             value.GetStatus().String(),
		cancellationReasonColumn: value.GetCancellationReason(),
		createdAtColumn:          value.GetCreatedAt(),
		updatedAtColumn:          value.GetUpdatedAt(),
	}
}

func (x RewardOrders) execute(ctx context.Context, query squirrel.Sqlizer) (pgconn.CommandTag, error) {
	sql, arguments, err := query.ToSql()
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("build reward-order query: %w", err)
	}

	result, err := x.database.Exec(ctx, sql, arguments...)
	if err != nil {
		var pgError *pgconn.PgError
		switch {
		case errors.As(err, &pgError) && pgError.Code == uniqueViolationCode:
			return pgconn.CommandTag{}, aggregate_reward_order.ErrAlreadyExists
		default:
			return pgconn.CommandTag{}, fmt.Errorf("execute reward-order query: %w", err)
		}
	}

	return result, nil
}

func (x RewardOrders) load(
	ctx context.Context,
	rewardOrderID model_reward_order_id.RewardOrderID,
	lock bool,
) (*aggregate_reward_order.RewardOrder, error) {
	query := x.queries.Select(
		idColumn,
		rewardIDColumn,
		coordinatorIDColumn,
		volunteerIDColumn,
		rewardTitleColumn,
		priceCoinsColumn,
		requestIDColumn,
		statusColumn,
		cancellationReasonColumn,
		createdAtColumn,
		updatedAtColumn,
	).From("reward_orders").Where(squirrel.Eq{idColumn: rewardOrderID.String()})
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
			return nil, aggregate_reward_order.ErrNotFound
		default:
			return nil, err
		}
	}

	return record.toAggregate()
}

func NewRewardOrders(database Database) RewardOrders {
	return RewardOrders{
		database: database,
		queries:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
