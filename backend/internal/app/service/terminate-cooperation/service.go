package service_terminate_cooperation

import (
	"context"
	"errors"
	"fmt"
	"time"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_cooperation_command "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation/command"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	aggregate_participation_command "github.com/Roum1212/hahaton/internal/domain/aggregate/participation/command"
	aggregate_reward "github.com/Roum1212/hahaton/internal/domain/aggregate/reward"
	aggregate_reward_order "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order"
	aggregate_reward_order_command "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order/command"
	aggregate_reward_command "github.com/Roum1212/hahaton/internal/domain/aggregate/reward/command"
	aggregate_wallet "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet"
	aggregate_wallet_command "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet/command"
	model_actor_kind "github.com/Roum1212/hahaton/internal/domain/model/actor-kind"
	model_cooperation_change "github.com/Roum1212/hahaton/internal/domain/model/cooperation-change"
	model_operation_id "github.com/Roum1212/hahaton/internal/domain/model/operation-id"
	model_participation_change "github.com/Roum1212/hahaton/internal/domain/model/participation-change"
	model_participation_status "github.com/Roum1212/hahaton/internal/domain/model/participation-status"
	model_reward_order_status "github.com/Roum1212/hahaton/internal/domain/model/reward-order-status"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrPermissionDenied = errors.New("permission denied")
)

type (
	Service interface {
		TerminateCooperation(ctx context.Context, r Request) error
	}
	service struct {
		cooperations   aggregate_cooperation.Repository
		coordinators   aggregate_coordinator.Repository
		participations aggregate_participation.Repository

		orders  aggregate_reward_order.Repository
		rewards aggregate_reward.Repository
		wallets aggregate_wallet.Repository
	}
)

func (x service) TerminateCooperation(ctx context.Context, r Request) error {
	if err := r.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	coordinatorID, volunteerID := r.GetCounterpartID(), r.GetUserID()
	if r.GetAsCoordinator() {
		if _, err := x.coordinators.Load(ctx, r.GetUserID()); err != nil {
			switch {
			case errors.Is(err, aggregate_coordinator.ErrNotFound):
				return ErrPermissionDenied
			default:
				return err
			}
		}

		coordinatorID, volunteerID = r.GetUserID(), r.GetCounterpartID()
	}

	return x.cooperations.Modify(ctx, coordinatorID, volunteerID, func(
		ctx context.Context,
		cooperation *aggregate_cooperation.Cooperation,
	) error {
		if err := x.terminate(ctx, cooperation); err != nil {
			return err
		}

		return x.recordCooperationChange(ctx, cooperation, r.GetUserID())
	})
}

func (x service) cancelParticipation(
	ctx context.Context,
	value *aggregate_participation.Participation,
	now time.Time,
) error {
	if value.GetStatus() != model_participation_status.Pending &&
		value.GetStatus() != model_participation_status.Approved {
		value.RevokeChat(aggregate_participation_command.NewRevokeChat(now))

		return x.participations.Save(ctx, value)
	}

	previous := utils_value.NewOptional(value.GetStatus())
	if err := value.Cancel(aggregate_participation_command.NewCancel("cooperation_terminated", now)); err != nil {
		return err
	}

	if err := x.participations.Save(ctx, value); err != nil {
		return err
	}

	return x.recordChange(ctx, value, previous, model_user_id.UserID{}, model_actor_kind.System)
}

func (x service) recordChange(
	ctx context.Context,
	value *aggregate_participation.Participation,
	previous utils_value.Optional[model_participation_status.ParticipationStatus],
	actorID model_user_id.UserID,
	actorKind model_actor_kind.ActorKind,
) error {
	id, err := model_operation_id.GenerateOperationID()
	if err != nil {
		return err
	}

	change, err := model_participation_change.NewParticipationChange(
		id,
		value.GetParticipationID(),
		previous,
		value.GetStatus(),
		actorID,
		actorKind,
		value.GetCancellationReason(),
		value.GetRejectionReason(),
		0,
		value.GetUpdatedAt(),
	)
	if err != nil {
		return err
	}

	return x.participations.AppendChange(ctx, change)
}

func (x service) recordCooperationChange(
	ctx context.Context,
	value *aggregate_cooperation.Cooperation,
	actorID model_user_id.UserID,
) error {
	change, err := model_cooperation_change.NewCooperationChange(value.GetCoordinatorID(), value.GetVolunteerID(),
		value.GetStatus(), actorID, value.GetRetryAfter(), value.GetUpdatedAt())
	if err != nil {
		return err
	}

	return x.cooperations.AppendChange(ctx, change)
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

func (x service) terminate(ctx context.Context, cooperation *aggregate_cooperation.Cooperation) error {
	now := time.Now().UTC()
	if err := cooperation.Terminate(aggregate_cooperation_command.NewTerminate(now)); err != nil {
		return err
	}

	participations, err := x.participations.List(ctx)
	if err != nil {
		return err
	}

	for _, value := range participations {
		if value.GetCoordinatorID() != cooperation.GetCoordinatorID() ||
			value.GetVolunteerID() != cooperation.GetVolunteerID() {
			continue
		}

		if err := x.cancelParticipation(ctx, value, now); err != nil {
			return err
		}
	}

	orders, err := x.orders.List(ctx)
	if err != nil {
		return err
	}

	for _, order := range orders {
		if order.GetCoordinatorID() != cooperation.GetCoordinatorID() ||
			order.GetVolunteerID() != cooperation.GetVolunteerID() || order.GetStatus() != model_reward_order_status.Pending {
			continue
		}

		if err := x.rewards.Modify(ctx, order.GetRewardID(), func(
			ctx context.Context,
			reward *aggregate_reward.Reward,
		) error {
			return x.refundOrder(ctx, order, reward, "cooperation_terminated")
		}); err != nil {
			return err
		}
	}

	return nil
}

func NewService(
	cooperations aggregate_cooperation.Repository,
	coordinators aggregate_coordinator.Repository,
	participations aggregate_participation.Repository,

	orders aggregate_reward_order.Repository,
	rewards aggregate_reward.Repository,
	wallets aggregate_wallet.Repository,
) Service {
	return service{
		cooperations:   cooperations,
		coordinators:   coordinators,
		participations: participations,

		orders:  orders,
		rewards: rewards,
		wallets: wallets,
	}
}
