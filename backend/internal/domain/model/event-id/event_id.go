package model_event_id

import (
	"fmt"

	types_uuid_v7 "github.com/Roum1212/hahaton/internal/pkg/types/uuid-v7"
)

type EventID struct {
	value types_uuid_v7.UUIDv7
}

func (x EventID) IsZero() bool {
	return x.value.IsZero()
}

func (x EventID) String() string {
	return x.value.String()
}

func (x EventID) Validate() error {
	return x.value.Validate()
}

func GenerateEventID() (EventID, error) {
	id, err := types_uuid_v7.GenerateUUIDv7()
	if err != nil {
		return EventID{}, fmt.Errorf("generate event id: %w", err)
	}

	return NewEventID(id.String())
}

func NewEventID(value string) (EventID, error) {
	id, err := types_uuid_v7.NewUUIDv7(value)
	if err != nil {
		return EventID{}, fmt.Errorf("parse event id: %w", err)
	}

	x := EventID{value: id}
	if err := x.Validate(); err != nil {
		return EventID{}, fmt.Errorf("invalid event id: %w", err)
	}

	return x, nil
}
