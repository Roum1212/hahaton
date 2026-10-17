package repository_reward_order_postgresql

import (
	"fmt"
	"time"

	aggregate_reward_order "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order"
	model_reward_id "github.com/Roum1212/hahaton/internal/domain/model/reward-id"
	model_reward_order_id "github.com/Roum1212/hahaton/internal/domain/model/reward-order-id"
	model_reward_order_status "github.com/Roum1212/hahaton/internal/domain/model/reward-order-status"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type row struct {
	rewardOrderID      string
	rewardID           string
	coordinatorID      string
	volunteerID        string
	rewardTitle        string
	priceCoins         int64
	requestID          string
	status             string
	cancellationReason string
	createdAt          time.Time
	updatedAt          time.Time
}

func (x *row) scan(source interface{ Scan(dest ...any) error }) error {
	return source.Scan(
		&x.rewardOrderID,
		&x.rewardID,
		&x.coordinatorID,
		&x.volunteerID,
		&x.rewardTitle,
		&x.priceCoins,
		&x.requestID,
		&x.status,
		&x.cancellationReason,
		&x.createdAt,
		&x.updatedAt,
	)
}

func (x *row) toAggregate() (*aggregate_reward_order.RewardOrder, error) {
	rewardOrderID, err := model_reward_order_id.NewRewardOrderID(x.rewardOrderID)
	if err != nil {
		return nil, fmt.Errorf("restore reward order id: %w", err)
	}

	rewardID, err := model_reward_id.NewRewardID(x.rewardID)
	if err != nil {
		return nil, fmt.Errorf("restore reward id: %w", err)
	}

	coordinatorID, err := model_user_id.NewUserID(x.coordinatorID)
	if err != nil {
		return nil, fmt.Errorf("restore coordinator id: %w", err)
	}

	volunteerID, err := model_user_id.NewUserID(x.volunteerID)
	if err != nil {
		return nil, fmt.Errorf("restore volunteer id: %w", err)
	}

	status, err := model_reward_order_status.NewRewardOrderStatus(x.status)
	if err != nil {
		return nil, fmt.Errorf("restore status: %w", err)
	}

	return aggregate_reward_order.NewRewardOrderFromFields(
		rewardOrderID,
		rewardID,
		coordinatorID,
		volunteerID,
		x.rewardTitle,
		x.priceCoins,
		x.requestID,
		status,
		x.cancellationReason,
		x.createdAt,
		x.updatedAt,
	), nil
}
