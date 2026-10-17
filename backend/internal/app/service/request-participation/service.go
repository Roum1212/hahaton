package service_request_participation

import (
	"context"
	"errors"
	"fmt"
	"time"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	aggregate_participation_command "github.com/Roum1212/hahaton/internal/domain/aggregate/participation/command"
	model_actor_kind "github.com/Roum1212/hahaton/internal/domain/model/actor-kind"
	model_cooperation_status "github.com/Roum1212/hahaton/internal/domain/model/cooperation-status"
	model_event_status "github.com/Roum1212/hahaton/internal/domain/model/event-status"
	model_operation_id "github.com/Roum1212/hahaton/internal/domain/model/operation-id"
	model_participation_change "github.com/Roum1212/hahaton/internal/domain/model/participation-change"
	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
	model_participation_status "github.com/Roum1212/hahaton/internal/domain/model/participation-status"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

var (
	ErrFailedPrecondition = errors.New("failed precondition")
	ErrInvalidRequest     = errors.New("invalid request")
)

type (
	Service interface {
		RequestParticipation(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		events         aggregate_event.Repository
		participations aggregate_participation.Repository
		cooperations   aggregate_cooperation.Repository
	}
)

func (x service) RequestParticipation(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	var participationID model_participation_id.ParticipationID

	if err := x.events.Modify(ctx, r.GetEventID(), func(
		ctx context.Context,
		event *aggregate_event.Event,
	) error {
		if event.GetCoordinatorID() == r.GetUserID() {
			return ErrInvalidRequest
		}

		if event.GetStatus() != model_event_status.Published ||
			!time.Now().UTC().Before(event.GetSchedule().GetStartsAt()) {
			return ErrFailedPrecondition
		}

		cooperation, err := x.cooperations.Load(ctx, event.GetCoordinatorID(), r.GetUserID())
		if err != nil {
			switch {
			case errors.Is(err, aggregate_cooperation.ErrNotFound):
				return ErrFailedPrecondition
			default:
				return err
			}
		}

		if cooperation.GetStatus() != model_cooperation_status.Active {
			return ErrFailedPrecondition
		}

		values, err := x.participations.List(ctx)
		if err != nil {
			return err
		}

		for _, value := range values {
			if value.GetEventID() == event.GetEventID() && value.GetVolunteerID() == r.GetUserID() {
				return aggregate_participation.ErrAlreadyExists
			}
		}

		participationID, err = model_participation_id.GenerateParticipationID()
		if err != nil {
			return err
		}

		value := aggregate_participation.NewParticipation(
			participationID,
			event.GetEventID(),
			event.GetCoordinatorID(),
			r.GetUserID(),
		)
		if err := value.Request(aggregate_participation_command.NewRequest(time.Now().UTC())); err != nil {
			return err
		}

		if err := x.participations.Create(ctx, value); err != nil {
			return err
		}

		return x.recordChange(
			ctx,
			value,
			utils_value.Optional[model_participation_status.ParticipationStatus]{},
			r.GetUserID(),
			model_actor_kind.Volunteer,
		)
	}); err != nil {
		return Response{}, err
	}

	return NewResponse(participationID), nil
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
	participations aggregate_participation.Repository,
	cooperations aggregate_cooperation.Repository,
) Service {
	return service{
		events:         events,
		participations: participations,
		cooperations:   cooperations,
	}
}
