package service_list_my_event_chats

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_event_chat "github.com/Roum1212/hahaton/internal/domain/aggregate/event-chat"
	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	model_cooperation_status "github.com/Roum1212/hahaton/internal/domain/model/cooperation-status"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_participation_status "github.com/Roum1212/hahaton/internal/domain/model/participation-status"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	"github.com/Roum1212/hahaton/internal/pkg/page"
)

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrPermissionDenied = errors.New("permission denied")
)

type (
	Service interface {
		ListMyEventChats(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		events         aggregate_event.Repository
		coordinators   aggregate_coordinator.Repository
		cooperations   aggregate_cooperation.Repository
		participations aggregate_participation.Repository
		chats          aggregate_event_chat.Repository
	}
)

func (x service) ListMyEventChats(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	events, err := x.events.List(ctx)
	if err != nil {
		return Response{}, err
	}

	filtered := make([]*aggregate_event.Event, 0)
	chats := make(map[model_event_id.EventID]*aggregate_event_chat.EventChat)
	boundaries := make(map[model_event_id.EventID]uint64)

	for _, event := range events {
		_, boundary, err := x.chatAccess(ctx, event.GetEventID(), r.GetUserID())
		if err != nil {
			if errors.Is(err, ErrPermissionDenied) {
				continue
			}

			return Response{}, err
		}

		chat, err := x.chats.Load(ctx, event.GetEventID())
		if err != nil {
			return Response{}, err
		}

		filtered = append(filtered, event)
		chats[event.GetEventID()] = chat
		boundaries[event.GetEventID()] = boundary
	}

	slices.SortFunc(filtered, func(
		a *aggregate_event.Event,
		b *aggregate_event.Event,
	) int {
		if value := b.GetSchedule().GetStartsAt().Compare(a.GetSchedule().GetStartsAt()); value != 0 {
			return value
		}

		return strings.Compare(b.GetEventID().String(), a.GetEventID().String())
	})
	values, total := page.Slice(filtered, r.GetPage())

	return NewResponse(values, chats, boundaries, total), nil
}

func (x service) chatAccess(
	ctx context.Context,
	eventID model_event_id.EventID,
	userID model_user_id.UserID,
) (*aggregate_event.Event, uint64, error) {
	event, err := x.events.Load(ctx, eventID)
	if err != nil {
		return nil, 0, err
	}

	if event.GetCoordinatorID() == userID {
		if _, err := x.coordinators.Load(ctx, userID); err != nil {
			return nil, 0, err
		}

		return event, 0, nil
	}

	cooperation, err := x.cooperations.Load(ctx, event.GetCoordinatorID(), userID)
	if err != nil {
		switch {
		case errors.Is(err, aggregate_cooperation.ErrNotFound):
			return nil, 0, ErrPermissionDenied
		default:
			return nil, 0, err
		}
	}

	if cooperation.GetStatus() != model_cooperation_status.Active {
		return nil, 0, ErrPermissionDenied
	}

	participations, err := x.participations.List(ctx)
	if err != nil {
		return nil, 0, err
	}

	for _, value := range participations {
		if value.GetEventID() != eventID || value.GetVolunteerID() != userID || !value.GetHasChatAccess() {
			continue
		}

		switch value.GetStatus() {
		case model_participation_status.Approved:
			return event, value.GetVisibleAfterSequence(), nil
		case model_participation_status.Confirmed:
			return event, value.GetVisibleAfterSequence(), nil
		case model_participation_status.NotAttended:
			return event, value.GetVisibleAfterSequence(), nil
		case model_participation_status.Pending:
		case model_participation_status.Rejected:
		case model_participation_status.Canceled:
		case model_participation_status.Closed:
		default:
		}
	}

	return nil, 0, ErrPermissionDenied
}

func NewService(
	events aggregate_event.Repository,
	coordinators aggregate_coordinator.Repository,
	cooperations aggregate_cooperation.Repository,
	participations aggregate_participation.Repository,
	chats aggregate_event_chat.Repository,
) Service {
	return service{
		events:         events,
		coordinators:   coordinators,
		cooperations:   cooperations,
		participations: participations,
		chats:          chats,
	}
}
