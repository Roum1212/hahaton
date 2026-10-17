package aggregate_reward_order

import (
	"errors"
	"time"

	aggregate_reward_order_command "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order/command"
	model_reward_id "github.com/Roum1212/hahaton/internal/domain/model/reward-id"
	model_reward_order_id "github.com/Roum1212/hahaton/internal/domain/model/reward-order-id"
	model_reward_order_status "github.com/Roum1212/hahaton/internal/domain/model/reward-order-status"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

var ErrInvalidState = errors.New("invalid reward order state")

type RewardOrder struct {
	rewardOrderID      model_reward_order_id.RewardOrderID
	rewardID           model_reward_id.RewardID
	coordinatorID      model_user_id.UserID
	volunteerID        model_user_id.UserID
	rewardTitle        string
	priceCoins         int64
	requestID          string
	status             model_reward_order_status.RewardOrderStatus
	cancellationReason string
	createdAt          time.Time
	updatedAt          time.Time
}

func (x *RewardOrder) GetCancellationReason() string {
	return x.cancellationReason
}

func (x *RewardOrder) GetCoordinatorID() model_user_id.UserID {
	return x.coordinatorID
}

func (x *RewardOrder) GetCreatedAt() time.Time {
	return x.createdAt
}

func (x *RewardOrder) GetPriceCoins() int64 {
	return x.priceCoins
}

func (x *RewardOrder) GetRequestID() string {
	return x.requestID
}

func (x *RewardOrder) GetRewardID() model_reward_id.RewardID {
	return x.rewardID
}

func (x *RewardOrder) GetRewardOrderID() model_reward_order_id.RewardOrderID {
	return x.rewardOrderID
}

func (x *RewardOrder) GetRewardTitle() string {
	return x.rewardTitle
}

func (x *RewardOrder) GetStatus() model_reward_order_status.RewardOrderStatus {
	return x.status
}

func (x *RewardOrder) GetUpdatedAt() time.Time {
	return x.updatedAt
}

func (x *RewardOrder) GetVolunteerID() model_user_id.UserID {
	return x.volunteerID
}

func (x *RewardOrder) Cancel(command aggregate_reward_order_command.Cancel) error {
	if x.status == model_reward_order_status.Canceled {
		return nil
	}

	if x.status != model_reward_order_status.Pending {
		return ErrInvalidState
	}

	x.status = model_reward_order_status.Canceled
	x.cancellationReason = command.GetReason()
	x.updatedAt = command.GetNow()

	return nil
}

func (x *RewardOrder) Create(command aggregate_reward_order_command.Create) error {
	if !x.createdAt.IsZero() {
		return ErrInvalidState
	}

	x.rewardTitle = command.GetRewardTitle()
	x.priceCoins = command.GetPriceCoins()
	x.requestID = command.GetRequestID()
	x.status = model_reward_order_status.Pending
	x.createdAt = command.GetNow()
	x.updatedAt = command.GetNow()

	return nil
}

func (x *RewardOrder) Issue(command aggregate_reward_order_command.Issue) error {
	if x.status == model_reward_order_status.Issued {
		return nil
	}

	if x.status != model_reward_order_status.Pending {
		return ErrInvalidState
	}

	x.status = model_reward_order_status.Issued
	x.updatedAt = command.GetNow()

	return nil
}

func NewRewardOrder(
	rewardOrderID model_reward_order_id.RewardOrderID,
	rewardID model_reward_id.RewardID,
	coordinatorID model_user_id.UserID,
	volunteerID model_user_id.UserID,
) *RewardOrder {
	return &RewardOrder{
		rewardOrderID: rewardOrderID,
		rewardID:      rewardID,
		coordinatorID: coordinatorID,
		volunteerID:   volunteerID,
	}
}

func NewRewardOrderFromFields(
	rewardOrderID model_reward_order_id.RewardOrderID,
	rewardID model_reward_id.RewardID,
	coordinatorID model_user_id.UserID,
	volunteerID model_user_id.UserID,
	rewardTitle string,
	priceCoins int64,
	requestID string,
	status model_reward_order_status.RewardOrderStatus,
	cancellationReason string,
	createdAt time.Time,
	updatedAt time.Time,
) *RewardOrder {
	return &RewardOrder{
		rewardOrderID:      rewardOrderID,
		rewardID:           rewardID,
		coordinatorID:      coordinatorID,
		volunteerID:        volunteerID,
		rewardTitle:        rewardTitle,
		priceCoins:         priceCoins,
		requestID:          requestID,
		status:             status,
		cancellationReason: cancellationReason,
		createdAt:          createdAt,
		updatedAt:          updatedAt,
	}
}
