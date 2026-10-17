package model_participation_id

import (
	"fmt"

	types_uuid_v7 "github.com/Roum1212/hahaton/internal/pkg/types/uuid-v7"
)

type ParticipationID struct {
	value types_uuid_v7.UUIDv7
}

func (x ParticipationID) IsZero() bool {
	return x.value.IsZero()
}

func (x ParticipationID) String() string {
	return x.value.String()
}

func (x ParticipationID) Validate() error {
	return x.value.Validate()
}

func GenerateParticipationID() (ParticipationID, error) {
	id, err := types_uuid_v7.GenerateUUIDv7()
	if err != nil {
		return ParticipationID{}, fmt.Errorf("generate participation id: %w", err)
	}

	return NewParticipationID(id.String())
}

func NewParticipationID(value string) (ParticipationID, error) {
	id, err := types_uuid_v7.NewUUIDv7(value)
	if err != nil {
		return ParticipationID{}, fmt.Errorf("parse participation id: %w", err)
	}

	x := ParticipationID{value: id}
	if err := x.Validate(); err != nil {
		return ParticipationID{}, err
	}

	return x, nil
}
