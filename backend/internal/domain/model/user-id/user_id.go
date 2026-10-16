package model_user_id

import (
	"fmt"

	types_uuid_v7 "github.com/Roum1212/hahaton/internal/pkg/types/uuid-v7"
)

type UserID struct {
	value types_uuid_v7.UUIDv7
}

func (x UserID) IsZero() bool {
	return x.value.IsZero()
}

func (x UserID) String() string {
	return x.value.String()
}

func (x UserID) Validate() error {
	return x.value.Validate()
}

func NewUserID(value string) (UserID, error) {
	id, err := types_uuid_v7.NewUUIDv7(value)
	if err != nil {
		return UserID{}, fmt.Errorf("parse user id: %w", err)
	}

	x := UserID{value: id}
	if err = x.Validate(); err != nil {
		return UserID{}, fmt.Errorf("invalid user id: %w", err)
	}

	return x, nil
}
