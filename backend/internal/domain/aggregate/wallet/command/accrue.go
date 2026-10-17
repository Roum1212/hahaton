package aggregate_wallet_command

import (
	"time"

	model_event_result "github.com/Roum1212/hahaton/internal/domain/model/event-result"
	model_operation_id "github.com/Roum1212/hahaton/internal/domain/model/operation-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Accrue struct {
	operationID   model_operation_id.OperationID
	coordinatorID model_user_id.UserID
	result        model_event_result.EventResult
	description   string
	now           time.Time
}

func (x Accrue) GetCoordinatorID() model_user_id.UserID {
	return x.coordinatorID
}

func (x Accrue) GetDescription() string {
	return x.description
}

func (x Accrue) GetNow() time.Time {
	return x.now
}

func (x Accrue) GetOperationID() model_operation_id.OperationID {
	return x.operationID
}

func (x Accrue) GetResult() model_event_result.EventResult {
	return x.result
}

func NewAccrue(
	operationID model_operation_id.OperationID,
	coordinatorID model_user_id.UserID,
	result model_event_result.EventResult,
	description string,
	now time.Time,
) Accrue {
	return Accrue{
		operationID:   operationID,
		coordinatorID: coordinatorID,
		result:        result,
		description:   description,
		now:           now,
	}
}
