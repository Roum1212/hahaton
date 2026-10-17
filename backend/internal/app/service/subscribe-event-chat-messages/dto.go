package service_subscribe_event_chat_messages

import (
	"fmt"

	model_chat_message "github.com/Roum1212/hahaton/internal/domain/model/chat-message"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_session_id "github.com/Roum1212/hahaton/internal/domain/model/session-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Request struct {
	userID        model_user_id.UserID
	eventID       model_event_id.EventID
	sessionID     model_session_id.SessionID
	afterSequence uint64
}

func (x Request) GetAfterSequence() uint64 {
	return x.afterSequence
}

func (x Request) GetEventID() model_event_id.EventID {
	return x.eventID
}

func (x Request) GetSessionID() model_session_id.SessionID {
	return x.sessionID
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

	if err := x.sessionID.Validate(); err != nil {
		return fmt.Errorf("session id: %w", err)
	}

	return nil
}

func NewRequest(
	userID model_user_id.UserID,
	eventID model_event_id.EventID,
	sessionID model_session_id.SessionID,
	afterSequence uint64,
) Request {
	return Request{
		userID:        userID,
		eventID:       eventID,
		sessionID:     sessionID,
		afterSequence: afterSequence,
	}
}

type Response struct {
	messages             []model_chat_message.ChatMessage
	readOnly             bool
	visibleAfterSequence uint64
}

func (x Response) GetMessages() []model_chat_message.ChatMessage {
	return x.messages
}

func (x Response) GetReadOnly() bool {
	return x.readOnly
}

func (x Response) GetVisibleAfterSequence() uint64 {
	return x.visibleAfterSequence
}

func NewResponse(messages []model_chat_message.ChatMessage, readOnly bool, visibleAfterSequence uint64) Response {
	return Response{
		messages:             messages,
		readOnly:             readOnly,
		visibleAfterSequence: visibleAfterSequence,
	}
}
