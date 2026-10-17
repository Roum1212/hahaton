package service_get_update_recipients

import (
	"context"
	"errors"
	"fmt"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	aggregate_reward "github.com/Roum1212/hahaton/internal/domain/aggregate/reward"
	aggregate_reward_order "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order"
	model_cooperation_status "github.com/Roum1212/hahaton/internal/domain/model/cooperation-status"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

var ErrInvalidRequest = errors.New("invalid request")

type (
	Service interface {
		GetUpdateRecipients(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		cooperations   aggregate_cooperation.Repository
		events         aggregate_event.Repository
		participations aggregate_participation.Repository
		rewards        aggregate_reward.Repository
		orders         aggregate_reward_order.Repository
	}
)

func (x service) GetUpdateRecipients(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	users := make(map[model_user_id.UserID]struct{})
	add := func(userID model_user_id.UserID) {
		if !userID.IsZero() {
			users[userID] = struct{}{}
		}
	}
	add(r.GetUserID())
	add(r.GetCounterpartID())
	coordinatorID := r.GetUserID()
	eventID := r.GetEventID()

	if !r.GetParticipationID().IsZero() {
		value, err := x.participations.Load(ctx, r.GetParticipationID())
		if err != nil {
			return Response{}, err
		}

		coordinatorID, eventID = value.GetCoordinatorID(), value.GetEventID()
		add(value.GetVolunteerID())
	}

	if !r.GetRewardOrderID().IsZero() {
		value, err := x.orders.Load(ctx, r.GetRewardOrderID())
		if err != nil {
			return Response{}, err
		}

		coordinatorID = value.GetCoordinatorID()
		add(value.GetVolunteerID())
	}

	if !r.GetRewardID().IsZero() {
		value, err := x.rewards.Load(ctx, r.GetRewardID())
		if err != nil {
			return Response{}, err
		}

		coordinatorID = value.GetCoordinatorID()
	}

	if !eventID.IsZero() {
		owner, participants, err := x.eventRecipients(ctx, eventID)
		if err != nil {
			return Response{}, err
		}

		coordinatorID = owner

		for _, participant := range participants {
			add(participant)
		}
	}

	add(coordinatorID)

	cooperations, err := x.cooperations.List(ctx)
	if err != nil {
		return Response{}, err
	}

	for _, cooperation := range cooperations {
		if cooperation.GetStatus() != model_cooperation_status.Active {
			continue
		}

		if cooperation.GetCoordinatorID() == coordinatorID {
			add(cooperation.GetVolunteerID())
		}

		if cooperation.GetVolunteerID() == r.GetUserID() {
			add(cooperation.GetCoordinatorID())
		}
	}

	result := make([]model_user_id.UserID, 0, len(users))
	for userID := range users {
		result = append(result, userID)
	}

	return NewResponse(result), nil
}

func (x service) eventRecipients(
	ctx context.Context,
	eventID model_event_id.EventID,
) (model_user_id.UserID, []model_user_id.UserID, error) {
	event, err := x.events.Load(ctx, eventID)
	if err != nil {
		return model_user_id.UserID{}, nil, err
	}

	participations, err := x.participations.List(ctx)
	if err != nil {
		return model_user_id.UserID{}, nil, err
	}

	users := make([]model_user_id.UserID, 0)

	for _, participation := range participations {
		if participation.GetEventID() == eventID {
			users = append(users, participation.GetVolunteerID())
		}
	}

	return event.GetCoordinatorID(), users, nil
}

func NewService(
	cooperations aggregate_cooperation.Repository,
	events aggregate_event.Repository,
	participations aggregate_participation.Repository,
	rewards aggregate_reward.Repository,
	orders aggregate_reward_order.Repository,
) Service {
	return service{
		cooperations:   cooperations,
		events:         events,
		participations: participations,
		rewards:        rewards,
		orders:         orders,
	}
}
