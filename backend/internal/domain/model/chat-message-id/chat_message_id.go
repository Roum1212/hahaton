package model_chat_message_id

import (
	"fmt"

	types_uuid_v7 "github.com/Roum1212/hahaton/internal/pkg/types/uuid-v7"
)

type ChatMessageID struct {
	value types_uuid_v7.UUIDv7
}

func (x ChatMessageID) IsZero() bool {
	return x.value.IsZero()
}

func (x ChatMessageID) String() string {
	return x.value.String()
}

func (x ChatMessageID) Validate() error {
	return x.value.Validate()
}

func GenerateChatMessageID() (ChatMessageID, error) {
	id, err := types_uuid_v7.GenerateUUIDv7()
	if err != nil {
		return ChatMessageID{}, fmt.Errorf("generate chat message id: %w", err)
	}

	return NewChatMessageID(id.String())
}

func NewChatMessageID(value string) (ChatMessageID, error) {
	id, err := types_uuid_v7.NewUUIDv7(value)
	if err != nil {
		return ChatMessageID{}, fmt.Errorf("parse chat message id: %w", err)
	}

	x := ChatMessageID{value: id}
	if err := x.Validate(); err != nil {
		return ChatMessageID{}, err
	}

	return x, nil
}
