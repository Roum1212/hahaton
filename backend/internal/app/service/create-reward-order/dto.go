package service_create_reward_order

import (
	"errors"
	"fmt"
	"strings"

	aggregate_reward_order "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order"
	model_reward_id "github.com/Roum1212/hahaton/internal/domain/model/reward-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Request struct {
	userID             model_user_id.UserID
	rewardID           model_reward_id.RewardID
	expectedPriceCoins int64
	requestID          string
}

func (x Request) GetExpectedPriceCoins() int64 {
	return x.expectedPriceCoins
}

func (x Request) GetRequestID() string {
	return x.requestID
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

	if x.expectedPriceCoins <= 0 {
		return errors.New("expected price must be positive")
	}

	if len(strings.TrimSpace(x.requestID)) == 0 {
		return errors.New("request id is required")
	}

	return nil
}

func NewRequest(
	userID model_user_id.UserID,
	rewardID model_reward_id.RewardID,
	expectedPriceCoins int64,
	requestID string,
) Request {
	return Request{
		userID:             userID,
		rewardID:           rewardID,
		expectedPriceCoins: expectedPriceCoins,
		requestID:          requestID,
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
