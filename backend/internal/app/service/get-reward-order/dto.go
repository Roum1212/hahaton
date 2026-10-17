package service_get_reward_order

import (
	"fmt"

	aggregate_reward_order "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order"
	model_reward_order_id "github.com/Roum1212/hahaton/internal/domain/model/reward-order-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Request struct {
	userID        model_user_id.UserID
	rewardOrderID model_reward_order_id.RewardOrderID
}

func (x Request) GetRewardOrderID() model_reward_order_id.RewardOrderID {
	return x.rewardOrderID
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if err := x.rewardOrderID.Validate(); err != nil {
		return fmt.Errorf("reward order id: %w", err)
	}

	return nil
}

func NewRequest(userID model_user_id.UserID, rewardOrderID model_reward_order_id.RewardOrderID) Request {
	return Request{
		userID:        userID,
		rewardOrderID: rewardOrderID,
	}
}

type Response struct {
	order *aggregate_reward_order.RewardOrder
}

func (x Response) GetOrder() *aggregate_reward_order.RewardOrder {
	return x.order
}

func NewResponse(order *aggregate_reward_order.RewardOrder) Response {
	return Response{
		order: order,
	}
}
