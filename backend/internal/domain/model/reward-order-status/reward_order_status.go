package model_reward_order_status

import (
	"errors"
)

type RewardOrderStatus string

const (
	Pending  RewardOrderStatus = "pending"
	Issued   RewardOrderStatus = "issued"
	Canceled RewardOrderStatus = "canceled"
)

func (x RewardOrderStatus) String() string {
	return string(x)
}

func (x RewardOrderStatus) Validate() error {
	switch x {
	case Pending:
		return nil
	case Issued:
		return nil
	case Canceled:
		return nil
	default:
		return errors.New("invalid reward order status")
	}
}

func NewRewardOrderStatus(value string) (RewardOrderStatus, error) {
	x := RewardOrderStatus(value)
	if err := x.Validate(); err != nil {
		return RewardOrderStatus(""), err
	}

	return x, nil
}
