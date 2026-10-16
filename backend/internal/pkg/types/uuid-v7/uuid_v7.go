package types_uuid_v7

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const version = 7

type UUIDv7 struct {
	value uuid.UUID
}

func (x UUIDv7) IsZero() bool {
	return x.value == uuid.Nil
}

func (x UUIDv7) String() string {
	if x.IsZero() {
		return ""
	}

	return x.value.String()
}

func (x UUIDv7) Validate() error {
	switch {
	case x.IsZero():
		return errors.New("uuid is empty")
	case x.value.Variant() != uuid.RFC4122:
		return fmt.Errorf("invalid uuid variant: %v", x.value.Variant())
	case x.value.Version() != version:
		return fmt.Errorf("invalid uuid version: %v", x.value.Version())
	default:
		return nil
	}
}

func GenerateUUIDv7() (UUIDv7, error) {
	value, err := uuid.NewV7()
	if err != nil {
		return UUIDv7{}, fmt.Errorf("generate uuid v7: %w", err)
	}

	return NewUUIDv7(value.String())
}

func NewUUIDv7(value string) (UUIDv7, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return UUIDv7{}, fmt.Errorf("parse uuid: %w", err)
	}

	x := UUIDv7{value: id}
	if err = x.Validate(); err != nil {
		return UUIDv7{}, fmt.Errorf("invalid uuid v7: %w", err)
	}

	return x, nil
}
