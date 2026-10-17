package model_wallet_operation_type

import (
	"errors"
)

type OperationType int32

const (
	ParticipationAccrual OperationType = 1
	RewardPurchase       OperationType = 2
	RewardRefund         OperationType = 3
)

func (x OperationType) Int32() int32 {
	return int32(x)
}

func (x OperationType) Validate() error {
	switch x {
	case ParticipationAccrual:
		return nil
	case RewardPurchase:
		return nil
	case RewardRefund:
		return nil
	default:
		return errors.New("invalid operation type")
	}
}

func NewOperationType(value int32) (OperationType, error) {
	x := OperationType(value)
	if err := x.Validate(); err != nil {
		return 0, err
	}

	return x, nil
}
