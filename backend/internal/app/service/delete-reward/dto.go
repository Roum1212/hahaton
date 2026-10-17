package service_delete_reward

import (
	"fmt"

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
