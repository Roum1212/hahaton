package rpc_service_reward_serde

import (
	reward_v1 "github.com/Roum1212/hahaton/genproto/reward/v1"
	aggregate_reward "github.com/Roum1212/hahaton/internal/domain/aggregate/reward"
)

func AggregateToPBReward(value *aggregate_reward.Reward) *reward_v1.Reward {
	return &reward_v1.Reward{
		Id:                value.GetRewardID().String(),
		CoordinatorId:     value.GetCoordinatorID().String(),
		Title:             value.GetTitle(),
		Description:       value.GetDescription(),
		ImageUrl:          value.GetImagePath().GetURL(),
		PriceCoins:        value.GetPriceCoins(),
		QuantityLimit:     value.GetQuantityLimit(),
		AvailableQuantity: value.GetAvailableQuantity(),
		Status:            reward_v1.RewardStatus_REWARD_STATUS_ACTIVE,
	}
}
