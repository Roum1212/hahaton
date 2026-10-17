package service_delete_reward

import (
	"context"
	"errors"
	"fmt"
	"sort"
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
		DeleteReward(ctx context.Context, r Request) error
	}
	service struct {
		rewards      aggregate_reward.Repository
		coordinators aggregate_coordinator.Repository
		orders       aggregate_reward_order.Repository
		wallets      aggregate_wallet.Repository
	}
)

func (x service) DeleteReward(ctx context.Context, r Request) error {
	if err := r.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	if _, err := x.coordinators.Load(ctx, r.GetUserID()); err != nil {
		switch {
		case errors.Is(err, aggregate_coordinator.ErrNotFound):
			return ErrPermissionDenied
		default:
			return err
		}
	}

	return x.rewards.Modify(ctx, r.GetRewardID(), func(
		ctx context.Context,
		reward *aggregate_reward.Reward,
	) error {
		if reward.GetCoordinatorID() != r.GetUserID() {
			return ErrPermissionDenied
		}

		if reward.GetDeleted() {
			return nil
		}

		orders, err := x.orders.List(ctx)
		if err != nil {
			return err
		}

		// Match the wallet lock order used by event completion.
		sort.Slice(orders, func(i, j int) bool {
			return orders[i].GetVolunteerID().String() < orders[j].GetVolunteerID().String()
		})

		for _, order := range orders {
			if order.GetRewardID() != reward.GetRewardID() || order.GetStatus() != model_reward_order_status.Pending {
				continue
			}

			if err := x.refundOrder(ctx, order, reward, "reward_deleted"); err != nil {
				return err
			}
		}

		return reward.Delete(aggregate_reward_command.NewDelete(time.Now().UTC()))
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
	rewards aggregate_reward.Repository,
	coordinators aggregate_coordinator.Repository,
	orders aggregate_reward_order.Repository,
	wallets aggregate_wallet.Repository,
) Service {
	return service{
		rewards:      rewards,
		coordinators: coordinators,
		orders:       orders,
		wallets:      wallets,
	}
}
