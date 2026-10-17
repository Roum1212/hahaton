package model_chat_message

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	model_actor_kind "github.com/Roum1212/hahaton/internal/domain/model/actor-kind"
	model_chat_message_id "github.com/Roum1212/hahaton/internal/domain/model/chat-message-id"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

const MaxTextLength = 2000

type ChatMessage struct {
	messageID         model_chat_message_id.ChatMessageID
	eventID           model_event_id.EventID
	authorID          model_user_id.UserID
	authorKind        model_actor_kind.ActorKind
	authorDisplayName string
	text              string
	requestID         string
	sequence          uint64
	createdAt         time.Time
}

func (x ChatMessage) GetAuthorDisplayName() string {
	return x.authorDisplayName
}

func (x ChatMessage) GetAuthorID() model_user_id.UserID {
	return x.authorID
}

func (x ChatMessage) GetAuthorKind() model_actor_kind.ActorKind {
	return x.authorKind
}

func (x ChatMessage) GetCreatedAt() time.Time {
	return x.createdAt
}

func (x ChatMessage) GetEventID() model_event_id.EventID {
	return x.eventID
}

func (x ChatMessage) GetMessageID() model_chat_message_id.ChatMessageID {
	return x.messageID
}

func (x ChatMessage) GetRequestID() string {
	return x.requestID
}

func (x ChatMessage) GetSequence() uint64 {
	return x.sequence
}

func (x ChatMessage) GetText() string {
	return x.text
}

func (x ChatMessage) Validate() error {
	if err := x.messageID.Validate(); err != nil {
		return err
	}

	if err := x.eventID.Validate(); err != nil {
		return err
	}

	if err := x.authorID.Validate(); err != nil {
		return err
	}

	if x.authorKind != model_actor_kind.Coordinator && x.authorKind != model_actor_kind.Volunteer {
		return errors.New("invalid chat author kind")
	}

	if len(strings.TrimSpace(x.text)) == 0 || utf8.RuneCountInString(x.text) > MaxTextLength {
		return errors.New("message text must contain 1 to 2000 characters")
	}

	if len(strings.TrimSpace(x.requestID)) == 0 || x.sequence == 0 {
		return errors.New("message request id and sequence are required")
	}

	return nil
}

func NewChatMessage(
	messageID model_chat_message_id.ChatMessageID,
	eventID model_event_id.EventID,
	authorID model_user_id.UserID,
	authorKind model_actor_kind.ActorKind,
	authorDisplayName string,
	text string,
	requestID string,
	sequence uint64,
	createdAt time.Time,
) (ChatMessage, error) {
	x := ChatMessage{
		messageID:         messageID,
		eventID:           eventID,
		authorID:          authorID,
		authorKind:        authorKind,
		authorDisplayName: authorDisplayName,
		text:              text,
		requestID:         requestID,
		sequence:          sequence,
		createdAt:         createdAt,
	}
	if err := x.Validate(); err != nil {
		return ChatMessage{}, err
	}

	return x, nil
}
