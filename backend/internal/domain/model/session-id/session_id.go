package model_session_id

import (
	"fmt"

	types_uuid_v7 "github.com/Roum1212/hahaton/internal/pkg/types/uuid-v7"
)

type SessionID struct {
	value types_uuid_v7.UUIDv7
}

func (x SessionID) IsZero() bool {
	return x.value.IsZero()
}

func (x SessionID) String() string {
	return x.value.String()
}

func (x SessionID) Validate() error {
	return x.value.Validate()
}

func GenerateSessionID() (SessionID, error) {
	value, err := types_uuid_v7.GenerateUUIDv7()
	if err != nil {
		return SessionID{}, fmt.Errorf("generate session id: %w", err)
	}

	return NewSessionID(value.String())
}

func NewSessionID(value string) (SessionID, error) {
	id, err := types_uuid_v7.NewUUIDv7(value)
	if err != nil {
		return SessionID{}, fmt.Errorf("parse session id: %w", err)
	}

	x := SessionID{value: id}
	if err = x.Validate(); err != nil {
		return SessionID{}, fmt.Errorf("invalid session id: %w", err)
	}

	return x, nil
}
