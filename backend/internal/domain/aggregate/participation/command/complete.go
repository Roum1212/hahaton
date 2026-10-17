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
