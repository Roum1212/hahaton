package service_cancel_event

import (
	"context"
	"errors"
	"fmt"
	"time"

	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_event_command "github.com/Roum1212/hahaton/internal/domain/aggregate/event/command"
	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	aggregate_participation_command "github.com/Roum1212/hahaton/internal/domain/aggregate/participation/command"
	model_actor_kind "github.com/Roum1212/hahaton/internal/domain/model/actor-kind"
	model_operation_id "github.com/Roum1212/hahaton/internal/domain/model/operation-id"
	model_participation_change "github.com/Roum1212/hahaton/internal/domain/model/participation-change"
	model_participation_status "github.com/Roum1212/hahaton/internal/domain/model/participation-status"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrPermissionDenied = errors.New("permission denied")
)

type (
	Service interface {
		CancelEvent(ctx context.Context, r Request) error
	}
	service struct {
		events         aggregate_event.Repository
		coordinators   aggregate_coordinator.Repository
		participations aggregate_participation.Repository
	}
)

func (x service) CancelEvent(ctx context.Context, r Request) error {
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

	return x.events.Modify(ctx, r.GetEventID(), func(
		ctx context.Context,
		event *aggregate_event.Event,
	) error {
		if event.GetCoordinatorID() != r.GetUserID() {
			return ErrPermissionDenied
		}

		now := time.Now().UTC()
		if err := event.Cancel(aggregate_event_command.NewCancel(now)); err != nil {
			return err
		}

		participations, err := x.participations.List(ctx)
		if err != nil {
			return err
		}

		for _, value := range participations {
			if value.GetEventID() != event.GetEventID() {
				continue
			}

			if value.GetStatus() != model_participation_status.Pending &&
				value.GetStatus() != model_participation_status.Approved {
				continue
			}

			previous := utils_value.NewOptional(value.GetStatus())
			if err := value.Cancel(aggregate_participation_command.NewCancel("event_canceled", now)); err != nil {
				return err
			}

			if err := x.participations.Save(ctx, value); err != nil {
				return err
			}

			if err := x.recordChange(ctx, value, previous, model_user_id.UserID{}, model_actor_kind.System); err != nil {
				return err
			}
		}

		return nil
	})
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

func NewService(
	events aggregate_event.Repository,
	coordinators aggregate_coordinator.Repository,
	participations aggregate_participation.Repository,
) Service {
	return service{
		events:         events,
		coordinators:   coordinators,
		participations: participations,
	}
}
