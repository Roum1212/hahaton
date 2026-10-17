package model_wallet_operation_source

import (
	"errors"

	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
	model_reward_order_id "github.com/Roum1212/hahaton/internal/domain/model/reward-order-id"
)

type OperationSource struct {
	participationID model_participation_id.ParticipationID
	rewardOrderID   model_reward_order_id.RewardOrderID
}

func (x OperationSource) GetParticipationID() model_participation_id.ParticipationID {
	return x.participationID
}

func (x OperationSource) GetRewardOrderID() model_reward_order_id.RewardOrderID {
	return x.rewardOrderID
}

func (x OperationSource) String() string {
	if !x.participationID.IsZero() {
		return x.participationID.String()
	}

	return x.rewardOrderID.String()
}

func (x OperationSource) Validate() error {
	switch {
	case !x.participationID.IsZero() && x.rewardOrderID.IsZero():
		return x.participationID.Validate()
	case x.participationID.IsZero() && !x.rewardOrderID.IsZero():
		return x.rewardOrderID.Validate()
	default:
		return errors.New("operation must reference a participation or a reward order")
	}
}

func NewOperationSourceFromParticipationID(value model_participation_id.ParticipationID) (OperationSource, error) {
	x := OperationSource{participationID: value}
	if err := x.Validate(); err != nil {
		return OperationSource{}, err
	}

	return x, nil
}

func NewOperationSourceFromRewardOrderID(value model_reward_order_id.RewardOrderID) (OperationSource, error) {
	x := OperationSource{rewardOrderID: value}
	if err := x.Validate(); err != nil {
		return OperationSource{}, err
	}

	return x, nil
}
