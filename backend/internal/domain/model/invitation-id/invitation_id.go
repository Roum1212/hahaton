package model_invitation_id

import (
	"fmt"

	types_uuid_v7 "github.com/Roum1212/hahaton/internal/pkg/types/uuid-v7"
)

type InvitationID struct {
	value types_uuid_v7.UUIDv7
}

func (x InvitationID) IsZero() bool {
	return x.value.IsZero()
}

func (x InvitationID) String() string {
	return x.value.String()
}

func (x InvitationID) Validate() error {
	return x.value.Validate()
}

func GenerateInvitationID() (InvitationID, error) {
	id, err := types_uuid_v7.GenerateUUIDv7()
	if err != nil {
		return InvitationID{}, fmt.Errorf("generate invitation id: %w", err)
	}

	return NewInvitationID(id.String())
}

func NewInvitationID(value string) (InvitationID, error) {
	id, err := types_uuid_v7.NewUUIDv7(value)
	if err != nil {
		return InvitationID{}, fmt.Errorf("parse invitation id: %w", err)
	}

	x := InvitationID{value: id}
	if err := x.Validate(); err != nil {
		return InvitationID{}, err
	}

	return x, nil
}
