package model_reward_id

import (
	"fmt"

	types_uuid_v7 "github.com/Roum1212/hahaton/internal/pkg/types/uuid-v7"
)

type RewardID struct {
	value types_uuid_v7.UUIDv7
}

func (x RewardID) IsZero() bool {
	return x.value.IsZero()
}

func (x RewardID) String() string {
	return x.value.String()
}

func (x RewardID) Validate() error {
	return x.value.Validate()
}

func GenerateRewardID() (RewardID, error) {
	id, err := types_uuid_v7.GenerateUUIDv7()
	if err != nil {
		return RewardID{}, fmt.Errorf("generate reward id: %w", err)
	}

	return NewRewardID(id.String())
}

func NewRewardID(value string) (RewardID, error) {
	id, err := types_uuid_v7.NewUUIDv7(value)
	if err != nil {
		return RewardID{}, fmt.Errorf("parse reward id: %w", err)
	}

	x := RewardID{value: id}
	if err := x.Validate(); err != nil {
		return RewardID{}, err
	}

	return x, nil
}
