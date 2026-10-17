package aggregate_event_chat

import (
	"errors"

	aggregate_event_chat_command "github.com/Roum1212/hahaton/internal/domain/aggregate/event-chat/command"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
)

var ErrInvalidSequence = errors.New("invalid chat sequence")

type EventChat struct {
	eventID        model_event_id.EventID
	latestSequence uint64
}

func (x *EventChat) GetEventID() model_event_id.EventID {
	return x.eventID
}

func (x *EventChat) GetLatestSequence() uint64 {
	return x.latestSequence
}

func (x *EventChat) Append(command aggregate_event_chat_command.Append) error {
	if command.GetSequence() != x.latestSequence+1 {
		return ErrInvalidSequence
	}

	x.latestSequence = command.GetSequence()

	return nil
}

func NewEventChat(eventID model_event_id.EventID) *EventChat {
	return &EventChat{
		eventID: eventID,
	}
}

func NewEventChatFromFields(eventID model_event_id.EventID, latestSequence uint64) *EventChat {
	return &EventChat{
		eventID:        eventID,
		latestSequence: latestSequence,
	}
}
