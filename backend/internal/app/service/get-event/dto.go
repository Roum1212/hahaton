package service_get_event

import (
	"fmt"

	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
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
	event         *aggregate_event.Event
	approvedCount int32
	pendingCount  int32
}

func (x Response) GetApprovedCount() int32 {
	return x.approvedCount
}

func (x Response) GetEvent() *aggregate_event.Event {
	return x.event
}

func (x Response) GetPendingCount() int32 {
	return x.pendingCount
}

func NewResponse(event *aggregate_event.Event, approvedCount, pendingCount int32) Response {
	return Response{
		event:         event,
		approvedCount: approvedCount,
		pendingCount:  pendingCount,
	}
}
