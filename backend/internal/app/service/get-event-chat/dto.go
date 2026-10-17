package service_get_event_chat

import (
	"fmt"

	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_event_chat "github.com/Roum1212/hahaton/internal/domain/aggregate/event-chat"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Request struct {
	userID  model_user_id.UserID
	eventID model_event_id.EventID
}

func (x Request) GetEventID() model_event_id.EventID {
	return x.eventID
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if err := x.eventID.Validate(); err != nil {
		return fmt.Errorf("event id: %w", err)
	}

	return nil
}

func NewRequest(userID model_user_id.UserID, eventID model_event_id.EventID) Request {
	return Request{
		userID:  userID,
		eventID: eventID,
	}
}

type Response struct {
	event                *aggregate_event.Event
	chat                 *aggregate_event_chat.EventChat
	visibleAfterSequence uint64
}

func (x Response) GetChat() *aggregate_event_chat.EventChat {
	return x.chat
}

func (x Response) GetEvent() *aggregate_event.Event {
	return x.event
}

func (x Response) GetVisibleAfterSequence() uint64 {
	return x.visibleAfterSequence
}

func NewResponse(
	event *aggregate_event.Event,
	chat *aggregate_event_chat.EventChat,
	visibleAfterSequence uint64,
) Response {
	return Response{
		event:                event,
		chat:                 chat,
		visibleAfterSequence: visibleAfterSequence,
	}
}
