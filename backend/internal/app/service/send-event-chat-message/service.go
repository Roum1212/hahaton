package service_send_event_chat_message

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_event_chat "github.com/Roum1212/hahaton/internal/domain/aggregate/event-chat"
	aggregate_event_chat_command "github.com/Roum1212/hahaton/internal/domain/aggregate/event-chat/command"
	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
	model_actor_kind "github.com/Roum1212/hahaton/internal/domain/model/actor-kind"
	model_chat_message "github.com/Roum1212/hahaton/internal/domain/model/chat-message"
	model_chat_message_id "github.com/Roum1212/hahaton/internal/domain/model/chat-message-id"
	model_cooperation_status "github.com/Roum1212/hahaton/internal/domain/model/cooperation-status"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_event_status "github.com/Roum1212/hahaton/internal/domain/model/event-status"
	model_participation_status "github.com/Roum1212/hahaton/internal/domain/model/participation-status"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

var (
	ErrAlreadyExists      = errors.New("already exists")
	ErrFailedPrecondition = errors.New("failed precondition")
	ErrInvalidRequest     = errors.New("invalid request")
	ErrPermissionDenied   = errors.New("permission denied")
)

type (
	Service interface {
		SendEventChatMessage(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		events         aggregate_event.Repository
		coordinators   aggregate_coordinator.Repository
		cooperations   aggregate_cooperation.Repository
		participations aggregate_participation.Repository
		chats          aggregate_event_chat.Repository
		volunteers     aggregate_volunteer.Repository
	}
)

func (x service) SendEventChatMessage(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	var sent model_chat_message.ChatMessage

	if err := x.chats.Modify(ctx, r.GetEventID(), func(
		ctx context.Context,
		chat *aggregate_event_chat.EventChat,
	) error {
		message, err := x.send(ctx, r, chat)
		if err != nil {
			return err
		}

		sent = message

		return nil
	}); err != nil {
		return Response{}, err
	}

	return NewResponse(sent), nil
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

func (x service) send(
	ctx context.Context,
	r Request,
	chat *aggregate_event_chat.EventChat,
) (model_chat_message.ChatMessage, error) {
	event, boundary, err := x.chatAccess(ctx, r.GetEventID(), r.GetUserID())
	if err != nil {
		return model_chat_message.ChatMessage{}, err
	}

	if event.GetStatus() != model_event_status.Published {
		return model_chat_message.ChatMessage{}, ErrFailedPrecondition
	}

	authorKind := model_actor_kind.Volunteer
	if event.GetCoordinatorID() == r.GetUserID() {
		authorKind = model_actor_kind.Coordinator
	}

	messages, err := x.chats.ListMessages(ctx, r.GetEventID())
	if err != nil {
		return model_chat_message.ChatMessage{}, err
	}

	for _, message := range messages {
		if message.GetAuthorID() != r.GetUserID() ||
			message.GetAuthorKind() != authorKind || message.GetRequestID() != r.GetRequestID() {
			continue
		}

		if message.GetSequence() <= boundary {
			return model_chat_message.ChatMessage{}, ErrFailedPrecondition
		}

		if message.GetText() != r.GetText() {
			return model_chat_message.ChatMessage{}, ErrAlreadyExists
		}

		return message, nil
	}

	volunteer, err := x.volunteers.Load(ctx, r.GetUserID())
	if err != nil {
		return model_chat_message.ChatMessage{}, err
	}

	name := volunteer.GetFullName()
	displayName := strings.TrimSpace(name.GetLastName() + " " + name.GetFirstName() + " " + name.GetMiddleName())

	id, err := model_chat_message_id.GenerateChatMessageID()
	if err != nil {
		return model_chat_message.ChatMessage{}, err
	}

	message, err := model_chat_message.NewChatMessage(
		id,
		r.GetEventID(),
		r.GetUserID(),
		authorKind,
		displayName,
		r.GetText(),
		r.GetRequestID(),
		chat.GetLatestSequence()+1,
		time.Now().UTC(),
	)
	if err != nil {
		return model_chat_message.ChatMessage{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	if err := chat.Append(aggregate_event_chat_command.NewAppend(message.GetSequence())); err != nil {
		return model_chat_message.ChatMessage{}, err
	}

	if err := x.chats.AppendMessage(ctx, message); err != nil {
		return model_chat_message.ChatMessage{}, err
	}

	return message, nil
}

func NewService(
	events aggregate_event.Repository,
	coordinators aggregate_coordinator.Repository,
	cooperations aggregate_cooperation.Repository,
	participations aggregate_participation.Repository,
	chats aggregate_event_chat.Repository,
	volunteers aggregate_volunteer.Repository,
) Service {
	return service{
		events:         events,
		coordinators:   coordinators,
		cooperations:   cooperations,
		participations: participations,
		chats:          chats,
		volunteers:     volunteers,
	}
}
