package model_reward_order_id

import (
	"fmt"

	types_uuid_v7 "github.com/Roum1212/hahaton/internal/pkg/types/uuid-v7"
)

type RewardOrderID struct {
	value types_uuid_v7.UUIDv7
}

func (x RewardOrderID) IsZero() bool {
	return x.value.IsZero()
}

func (x RewardOrderID) String() string {
	return x.value.String()
}

func (x RewardOrderID) Validate() error {
	return x.value.Validate()
}

func GenerateRewardOrderID() (RewardOrderID, error) {
	id, err := types_uuid_v7.GenerateUUIDv7()
	if err != nil {
		return RewardOrderID{}, fmt.Errorf("generate reward order id: %w", err)
	}

	return NewRewardOrderID(id.String())
}

func NewRewardOrderID(value string) (RewardOrderID, error) {
	id, err := types_uuid_v7.NewUUIDv7(value)
	if err != nil {
		return RewardOrderID{}, fmt.Errorf("parse reward order id: %w", err)
	}

	x := RewardOrderID{value: id}
	if err := x.Validate(); err != nil {
		return RewardOrderID{}, err
	}

	return x, nil
}
