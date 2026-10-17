package aggregate_event_chat

import (
	"context"
	"errors"

	model_chat_message "github.com/Roum1212/hahaton/internal/domain/model/chat-message"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
)

var ErrNotFound = errors.New("event chat not found")

type Repository interface {
	AppendMessage(ctx context.Context, message model_chat_message.ChatMessage) error
	Create(ctx context.Context, chat *EventChat) error
	ListMessages(ctx context.Context, eventID model_event_id.EventID) ([]model_chat_message.ChatMessage, error)
	Load(ctx context.Context, eventID model_event_id.EventID) (*EventChat, error)
	Modify(
		ctx context.Context,
		eventID model_event_id.EventID,
		change func(
			ctx context.Context,
			chat *EventChat,
		) error,
	) error
}
