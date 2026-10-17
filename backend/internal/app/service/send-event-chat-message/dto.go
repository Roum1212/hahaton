package service_send_event_chat_message

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	model_chat_message "github.com/Roum1212/hahaton/internal/domain/model/chat-message"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Request struct {
	userID    model_user_id.UserID
	eventID   model_event_id.EventID
	text      string
	requestID string
}

func (x Request) GetEventID() model_event_id.EventID {
	return x.eventID
}

func (x Request) GetRequestID() string {
	return x.requestID
}

func (x Request) GetText() string {
	return x.text
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

	if len(strings.TrimSpace(x.text)) == 0 || utf8.RuneCountInString(x.text) > 2000 {
		return errors.New("text must contain 1 to 2000 characters")
	}

	if len(strings.TrimSpace(x.requestID)) == 0 {
		return errors.New("request id is required")
	}

	return nil
}

func NewRequest(userID model_user_id.UserID, eventID model_event_id.EventID, text, requestID string) Request {
	return Request{
		userID:    userID,
		eventID:   eventID,
		text:      text,
		requestID: requestID,
	}
}

type Response struct {
	message model_chat_message.ChatMessage
}

func (x Response) GetMessage() model_chat_message.ChatMessage {
	return x.message
}

func NewResponse(message model_chat_message.ChatMessage) Response {
	return Response{
		message: message,
	}
}
