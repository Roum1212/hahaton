package aggregate_participation_command

import (
	"time"

	model_event_result "github.com/Roum1212/hahaton/internal/domain/model/event-result"
	model_operation_id "github.com/Roum1212/hahaton/internal/domain/model/operation-id"
)

type Complete struct {
	result      model_event_result.EventResult
	coins       int64
	operationID model_operation_id.OperationID
	force       bool
	now         time.Time
}

func (x Complete) GetCoins() int64 {
	return x.coins
}

func (x Complete) GetNow() time.Time {
	return x.now
}

func (x Complete) GetOperationID() model_operation_id.OperationID {
	return x.operationID
}

func (x Complete) GetResult() model_event_result.EventResult {
	return x.result
}

func (x Complete) IsForced() bool {
	return x.force
}

func NewComplete(
	result model_event_result.EventResult,
	coins int64,
	operationID model_operation_id.OperationID,
	now time.Time,
) Complete {
	return Complete{
		result:      result,
		coins:       coins,
		operationID: operationID,
		now:         now,
	}
}

func NewForceComplete(
	result model_event_result.EventResult,
	coins int64,
	operationID model_operation_id.OperationID,
	now time.Time,
) Complete {
	x := NewComplete(result, coins, operationID, now)
	x.force = true

	return x
}
