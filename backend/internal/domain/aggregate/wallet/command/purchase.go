package aggregate_wallet_command

import (
	"time"

	model_operation_id "github.com/Roum1212/hahaton/internal/domain/model/operation-id"
	model_reward_order_id "github.com/Roum1212/hahaton/internal/domain/model/reward-order-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Purchase struct {
	operationID   model_operation_id.OperationID
	coordinatorID model_user_id.UserID
	sourceID      model_reward_order_id.RewardOrderID
	description   string
	coins         int64
	now           time.Time
}

func (x Purchase) GetCoins() int64 {
	return x.coins
}

func (x Purchase) GetCoordinatorID() model_user_id.UserID {
	return x.coordinatorID
}

func (x Purchase) GetDescription() string {
	return x.description
}

func (x Purchase) GetNow() time.Time {
	return x.now
}

func (x Purchase) GetOperationID() model_operation_id.OperationID {
	return x.operationID
}

func (x Purchase) GetSourceID() model_reward_order_id.RewardOrderID {
	return x.sourceID
}

func NewPurchase(
	operationID model_operation_id.OperationID,
	coordinatorID model_user_id.UserID,
	sourceID model_reward_order_id.RewardOrderID,
	description string,
	coins int64,
	now time.Time,
) Purchase {
	return Purchase{
		operationID:   operationID,
		coordinatorID: coordinatorID,
		sourceID:      sourceID,
		description:   description,
		coins:         coins,
		now:           now,
	}
}
