package service_create_reward_order

import (
	"context"
	"errors"
	"fmt"
	"time"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_reward "github.com/Roum1212/hahaton/internal/domain/aggregate/reward"
	aggregate_reward_order "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order"
	aggregate_reward_order_command "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order/command"
	aggregate_reward_command "github.com/Roum1212/hahaton/internal/domain/aggregate/reward/command"
	aggregate_wallet "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet"
	aggregate_wallet_command "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet/command"
	model_cooperation_status "github.com/Roum1212/hahaton/internal/domain/model/cooperation-status"
	model_operation_id "github.com/Roum1212/hahaton/internal/domain/model/operation-id"
	model_reward_order_id "github.com/Roum1212/hahaton/internal/domain/model/reward-order-id"
)

var (
	ErrAlreadyExists      = errors.New("already exists")
	ErrFailedPrecondition = errors.New("failed precondition")
	ErrInvalidRequest     = errors.New("invalid request")
)

type (
	Service interface {
		CreateRewardOrder(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		rewards      aggregate_reward.Repository
		cooperations aggregate_cooperation.Repository
		orders       aggregate_reward_order.Repository
		wallets      aggregate_wallet.Repository
	}
)

func (x service) CreateRewardOrder(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	var created *aggregate_reward_order.RewardOrder

	if err := x.rewards.Modify(ctx, r.GetRewardID(), func(
		ctx context.Context,
		reward *aggregate_reward.Reward,
	) error {
		orders, err := x.orders.List(ctx)
		if err != nil {
			return err
		}

		for _, order := range orders {
			if order.GetVolunteerID() != r.GetUserID() || order.GetRequestID() != r.GetRequestID() {
				continue
			}

			if order.GetRewardID() != r.GetRewardID() || order.GetPriceCoins() != r.GetExpectedPriceCoins() {
				return ErrAlreadyExists
			}

			created = order

			return nil
		}

		order, err := x.purchase(ctx, r, reward)
		if err != nil {
			return err
		}

		created = order

		return nil
	}); err != nil {
		return Response{}, err
	}

	return NewResponse(created), nil
}

func (x service) purchase(
	ctx context.Context,
	r Request,
	reward *aggregate_reward.Reward,
) (*aggregate_reward_order.RewardOrder, error) {
	if reward.GetCoordinatorID() == r.GetUserID() {
		return nil, ErrInvalidRequest
	}

	if reward.GetDeleted() {
		return nil, aggregate_reward.ErrNotFound
	}

	cooperation, err := x.cooperations.Load(ctx, reward.GetCoordinatorID(), r.GetUserID())
	if err != nil {
		return nil, err
	}

	if cooperation.GetStatus() != model_cooperation_status.Active {
		return nil, ErrFailedPrecondition
	}

	now := time.Now().UTC()
	if err := reward.Reserve(aggregate_reward_command.NewReserve(r.GetExpectedPriceCoins(), now)); err != nil {
		return nil, err
	}

	id, err := model_reward_order_id.GenerateRewardOrderID()
	if err != nil {
		return nil, err
	}

	order := aggregate_reward_order.NewRewardOrder(id, reward.GetRewardID(), reward.GetCoordinatorID(), r.GetUserID())
	if err := order.Create(aggregate_reward_order_command.NewCreate(
		reward.GetTitle(),
		reward.GetPriceCoins(),
		r.GetRequestID(),
		now,
	)); err != nil {
		return nil, err
	}

	operationID, err := model_operation_id.GenerateOperationID()
	if err != nil {
		return nil, err
	}

	if err := x.wallets.Modify(ctx, r.GetUserID(), func(
		ctx context.Context,
		wallet *aggregate_wallet.Wallet,
	) error {
		operation, err := wallet.Purchase(aggregate_wallet_command.NewPurchase(
			operationID,
			reward.GetCoordinatorID(),
			id,
			reward.GetTitle(),
			order.GetPriceCoins(),
			now,
		))
		if err != nil {
			return err
		}

		return x.wallets.AppendOperation(ctx, operation)
	}); err != nil {
		return nil, err
	}

	if err := x.orders.Create(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func NewService(
	rewards aggregate_reward.Repository,
	cooperations aggregate_cooperation.Repository,
	orders aggregate_reward_order.Repository,
	wallets aggregate_wallet.Repository,
) Service {
	return service{
		rewards:      rewards,
		cooperations: cooperations,
		orders:       orders,
		wallets:      wallets,
	}
}
