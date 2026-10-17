package service_subscribe_event_chat_messages

import (
	"context"
	"errors"
	"fmt"
	"time"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_event_chat "github.com/Roum1212/hahaton/internal/domain/aggregate/event-chat"
	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	aggregate_session "github.com/Roum1212/hahaton/internal/domain/aggregate/session"
	model_chat_message "github.com/Roum1212/hahaton/internal/domain/model/chat-message"
	model_cooperation_status "github.com/Roum1212/hahaton/internal/domain/model/cooperation-status"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_event_status "github.com/Roum1212/hahaton/internal/domain/model/event-status"
	model_participation_status "github.com/Roum1212/hahaton/internal/domain/model/participation-status"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrPermissionDenied = errors.New("permission denied")
)

type (
	Service interface {
		SubscribeEventChatMessages(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		events         aggregate_event.Repository
		coordinators   aggregate_coordinator.Repository
		cooperations   aggregate_cooperation.Repository
		participations aggregate_participation.Repository
		chats          aggregate_event_chat.Repository
		sessions       aggregate_session.Repository
	}
)

func (x service) SubscribeEventChatMessages(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	session, err := x.sessions.Load(ctx, r.GetSessionID())
	if err != nil {
		switch {
		case errors.Is(err, aggregate_session.ErrNotFound):
			return Response{}, aggregate_session.ErrInvalidCredentials
		default:
			return Response{}, err
		}
	}

	if session.GetUserID() != r.GetUserID() || !session.IsActive(time.Now().UTC()) {
		return Response{}, aggregate_session.ErrInvalidCredentials
	}

	event, boundary, err := x.chatAccess(ctx, r.GetEventID(), r.GetUserID())
	if err != nil {
		return Response{}, err
	}

	chat, err := x.chats.Load(ctx, r.GetEventID())
	if err != nil {
		return Response{}, err
	}

	if r.GetAfterSequence() > chat.GetLatestSequence() {
		return Response{}, ErrInvalidRequest
	}

	after := max(boundary, r.GetAfterSequence())

	messages, err := x.chats.ListMessages(ctx, r.GetEventID())
	if err != nil {
		return Response{}, err
	}

	filtered := make([]model_chat_message.ChatMessage, 0)

	for _, message := range messages {
		if message.GetSequence() > after {
			filtered = append(filtered, message)
		}
	}

	return NewResponse(filtered, event.GetStatus() != model_event_status.Published, boundary), nil
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
	sessions aggregate_session.Repository,
) Service {
	return service{
		events:         events,
		coordinators:   coordinators,
		cooperations:   cooperations,
		participations: participations,
		chats:          chats,
		sessions:       sessions,
	}
}
