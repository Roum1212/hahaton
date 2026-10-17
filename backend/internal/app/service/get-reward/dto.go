package service_get_reward

import (
	"fmt"

	aggregate_reward "github.com/Roum1212/hahaton/internal/domain/aggregate/reward"
	model_reward_id "github.com/Roum1212/hahaton/internal/domain/model/reward-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Request struct {
	userID   model_user_id.UserID
	rewardID model_reward_id.RewardID
}

func (x Request) GetRewardID() model_reward_id.RewardID {
	return x.rewardID
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if err := x.rewardID.Validate(); err != nil {
		return fmt.Errorf("reward id: %w", err)
	}

	return nil
}

func NewRequest(userID model_user_id.UserID, rewardID model_reward_id.RewardID) Request {
	return Request{
		userID:   userID,
		rewardID: rewardID,
	}
}

type Response struct {
	reward *aggregate_reward.Reward
}

func (x Response) GetReward() *aggregate_reward.Reward {
	return x.reward
}

func NewResponse(reward *aggregate_reward.Reward) Response {
	return Response{
		reward: reward,
	}
}
