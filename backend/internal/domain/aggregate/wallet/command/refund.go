package aggregate_wallet_command

import (
	"time"

	model_operation_id "github.com/Roum1212/hahaton/internal/domain/model/operation-id"
	model_reward_order_id "github.com/Roum1212/hahaton/internal/domain/model/reward-order-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Refund struct {
	operationID   model_operation_id.OperationID
	coordinatorID model_user_id.UserID
	sourceID      model_reward_order_id.RewardOrderID
	description   string
	coins         int64
	now           time.Time
}

func (x Refund) GetCoins() int64 {
	return x.coins
}

func (x Refund) GetCoordinatorID() model_user_id.UserID {
	return x.coordinatorID
}

func (x Refund) GetDescription() string {
	return x.description
}

func (x Refund) GetNow() time.Time {
	return x.now
}

func (x Refund) GetOperationID() model_operation_id.OperationID {
	return x.operationID
}

func (x Refund) GetSourceID() model_reward_order_id.RewardOrderID {
	return x.sourceID
}

func NewRefund(
	operationID model_operation_id.OperationID,
	coordinatorID model_user_id.UserID,
	sourceID model_reward_order_id.RewardOrderID,
	description string,
	coins int64,
	now time.Time,
) Refund {
	return Refund{
		operationID:   operationID,
		coordinatorID: coordinatorID,
		sourceID:      sourceID,
		description:   description,
		coins:         coins,
		now:           now,
	}
}
