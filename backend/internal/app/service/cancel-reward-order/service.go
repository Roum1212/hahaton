package service_cancel_reward_order

import (
	"context"
	"errors"
	"fmt"
	"time"

	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_reward "github.com/Roum1212/hahaton/internal/domain/aggregate/reward"
	aggregate_reward_order "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order"
	aggregate_reward_order_command "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order/command"
	aggregate_reward_command "github.com/Roum1212/hahaton/internal/domain/aggregate/reward/command"
	aggregate_wallet "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet"
	aggregate_wallet_command "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet/command"
	model_operation_id "github.com/Roum1212/hahaton/internal/domain/model/operation-id"
	model_reward_order_status "github.com/Roum1212/hahaton/internal/domain/model/reward-order-status"
)

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrPermissionDenied = errors.New("permission denied")
)

type (
	Service interface {
		CancelRewardOrder(ctx context.Context, r Request) error
	}
	service struct {
		orders       aggregate_reward_order.Repository
		coordinators aggregate_coordinator.Repository
		rewards      aggregate_reward.Repository
		wallets      aggregate_wallet.Repository
	}
)

func (x service) CancelRewardOrder(ctx context.Context, r Request) error {
	if err := r.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	return x.orders.Modify(ctx, r.GetRewardOrderID(), func(
		ctx context.Context,
		order *aggregate_reward_order.RewardOrder,
	) error {
		if order.GetVolunteerID() != r.GetUserID() {
			if order.GetCoordinatorID() != r.GetUserID() {
				return ErrPermissionDenied
			}

			if _, err := x.coordinators.Load(ctx, r.GetUserID()); err != nil {
				switch {
				case errors.Is(err, aggregate_coordinator.ErrNotFound):
					return ErrPermissionDenied
				default:
					return err
				}
			}
		}

		if order.GetStatus() == model_reward_order_status.Canceled {
			return nil
		}

		if order.GetStatus() != model_reward_order_status.Pending {
			return aggregate_reward_order.ErrInvalidState
		}

		reason := "volunteer_canceled"
		if order.GetCoordinatorID() == r.GetUserID() {
			reason = "coordinator_canceled"
		}

		return x.rewards.Modify(ctx, order.GetRewardID(), func(
			ctx context.Context,
			reward *aggregate_reward.Reward,
		) error {
			return x.refundOrder(ctx, order, reward, reason)
		})
	})
}

func (x service) refundOrder(
	ctx context.Context,
	order *aggregate_reward_order.RewardOrder,
	reward *aggregate_reward.Reward,
	reason string,
) error {
	now := time.Now().UTC()

	operationID, err := model_operation_id.GenerateOperationID()
	if err != nil {
		return err
	}

	if err := order.Cancel(aggregate_reward_order_command.NewCancel(reason, now)); err != nil {
		return err
	}

	if err := x.wallets.Modify(ctx, order.GetVolunteerID(), func(
		ctx context.Context,
		wallet *aggregate_wallet.Wallet,
	) error {
		operation, err := wallet.Refund(aggregate_wallet_command.NewRefund(
			operationID,
			order.GetCoordinatorID(),
			order.GetRewardOrderID(),
			order.GetRewardTitle(),
			order.GetPriceCoins(),
			now,
		))
		if err != nil {
			return err
		}

		return x.wallets.AppendOperation(ctx, operation)
	}); err != nil {
		return err
	}

	if err := reward.Release(aggregate_reward_command.NewRelease(now)); err != nil {
		return err
	}

	return x.orders.Save(ctx, order)
}

func NewService(
	orders aggregate_reward_order.Repository,
	coordinators aggregate_coordinator.Repository,
	rewards aggregate_reward.Repository,
	wallets aggregate_wallet.Repository,
) Service {
	return service{
		orders:       orders,
		coordinators: coordinators,
		rewards:      rewards,
		wallets:      wallets,
	}
}
