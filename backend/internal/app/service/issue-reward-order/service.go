package service_issue_reward_order

import (
	"context"
	"errors"
	"fmt"
	"time"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_reward "github.com/Roum1212/hahaton/internal/domain/aggregate/reward"
	aggregate_reward_order "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order"
	aggregate_reward_order_command "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order/command"
	aggregate_reward_command "github.com/Roum1212/hahaton/internal/domain/aggregate/reward/command"
	model_cooperation_status "github.com/Roum1212/hahaton/internal/domain/model/cooperation-status"
	model_reward_order_status "github.com/Roum1212/hahaton/internal/domain/model/reward-order-status"
)

var (
	ErrFailedPrecondition = errors.New("failed precondition")
	ErrInvalidRequest     = errors.New("invalid request")
	ErrPermissionDenied   = errors.New("permission denied")
)

type (
	Service interface {
		IssueRewardOrder(ctx context.Context, r Request) error
	}
	service struct {
		orders       aggregate_reward_order.Repository
		coordinators aggregate_coordinator.Repository
		rewards      aggregate_reward.Repository
		cooperations aggregate_cooperation.Repository
	}
)

func (x service) IssueRewardOrder(ctx context.Context, r Request) error {
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

	return x.orders.Modify(ctx, r.GetRewardOrderID(), func(
		ctx context.Context,
		order *aggregate_reward_order.RewardOrder,
	) error {
		if order.GetCoordinatorID() != r.GetUserID() {
			return ErrPermissionDenied
		}

		if order.GetStatus() == model_reward_order_status.Issued {
			return nil
		}

		if order.GetStatus() != model_reward_order_status.Pending {
			return aggregate_reward_order.ErrInvalidState
		}

		cooperation, err := x.cooperations.Load(ctx, r.GetUserID(), order.GetVolunteerID())
		if err != nil {
			return err
		}

		if cooperation.GetStatus() != model_cooperation_status.Active {
			return ErrFailedPrecondition
		}

		now := time.Now().UTC()

		if err := x.rewards.Modify(ctx, order.GetRewardID(), func(
			ctx context.Context,
			reward *aggregate_reward.Reward,
		) error {
			return reward.Issue(aggregate_reward_command.NewIssue(now))
		}); err != nil {
			return err
		}

		return order.Issue(aggregate_reward_order_command.NewIssue(now))
	})
}

func NewService(
	orders aggregate_reward_order.Repository,
	coordinators aggregate_coordinator.Repository,
	rewards aggregate_reward.Repository,
	cooperations aggregate_cooperation.Repository,
) Service {
	return service{
		orders:       orders,
		coordinators: coordinators,
		rewards:      rewards,
		cooperations: cooperations,
	}
}
