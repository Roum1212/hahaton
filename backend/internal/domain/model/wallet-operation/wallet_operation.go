package model_wallet_operation

import (
	"errors"
	"time"

	model_operation_id "github.com/Roum1212/hahaton/internal/domain/model/operation-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	model_wallet_operation_source "github.com/Roum1212/hahaton/internal/domain/model/wallet-operation-source"
	model_wallet_operation_type "github.com/Roum1212/hahaton/internal/domain/model/wallet-operation-type"
)

type Operation struct {
	operationID           model_operation_id.OperationID
	coordinatorID         model_user_id.UserID
	volunteerID           model_user_id.UserID
	description           string
	operationType         model_wallet_operation_type.OperationType
	source                model_wallet_operation_source.OperationSource
	coinDelta             int64
	minuteDelta           int64
	coinBalanceAfter      int64
	confirmedMinutesAfter int64
	sequence              uint64
	createdAt             time.Time
}

func (x Operation) GetCoinBalanceAfter() int64 {
	return x.coinBalanceAfter
}

func (x Operation) GetCoinDelta() int64 {
	return x.coinDelta
}

func (x Operation) GetConfirmedMinutesAfter() int64 {
	return x.confirmedMinutesAfter
}

func (x Operation) GetCoordinatorID() model_user_id.UserID {
	return x.coordinatorID
}

func (x Operation) GetCreatedAt() time.Time {
	return x.createdAt
}

func (x Operation) GetDescription() string {
	return x.description
}

func (x Operation) GetMinuteDelta() int64 {
	return x.minuteDelta
}

func (x Operation) GetOperationID() model_operation_id.OperationID {
	return x.operationID
}

func (x Operation) GetOperationType() model_wallet_operation_type.OperationType {
	return x.operationType
}

func (x Operation) GetSequence() uint64 {
	return x.sequence
}

func (x Operation) GetSource() model_wallet_operation_source.OperationSource {
	return x.source
}

func (x Operation) GetVolunteerID() model_user_id.UserID {
	return x.volunteerID
}

func (x Operation) Validate() error {
	if err := x.operationID.Validate(); err != nil {
		return err
	}

	if err := x.volunteerID.Validate(); err != nil {
		return err
	}

	if err := x.coordinatorID.Validate(); err != nil {
		return err
	}

	if x.sequence == 0 || x.coinBalanceAfter < 0 || x.confirmedMinutesAfter < 0 {
		return errors.New("invalid wallet operation")
	}

	if err := x.source.Validate(); err != nil {
		return err
	}

	if err := x.operationType.Validate(); err != nil {
		return err
	}

	if (x.operationType == model_wallet_operation_type.ParticipationAccrual) != !x.source.GetParticipationID().IsZero() {
		return errors.New("operation type does not match its source")
	}

	switch x.operationType {
	case model_wallet_operation_type.ParticipationAccrual:
		if x.coinDelta <= 0 || x.minuteDelta <= 0 {
			return errors.New("invalid accrual deltas")
		}
	case model_wallet_operation_type.RewardPurchase:
		if x.coinDelta >= 0 || x.minuteDelta != 0 {
			return errors.New("invalid purchase deltas")
		}
	case model_wallet_operation_type.RewardRefund:
		if x.coinDelta <= 0 || x.minuteDelta != 0 {
			return errors.New("invalid refund deltas")
		}
	default:
		return errors.New("invalid operation type")
	}

	return nil
}

func NewOperation(
	operationID model_operation_id.OperationID,
	coordinatorID model_user_id.UserID,
	volunteerID model_user_id.UserID,
	description string,
	operationType model_wallet_operation_type.OperationType,
	source model_wallet_operation_source.OperationSource,
	coinDelta int64,
	minuteDelta int64,
	coinBalanceAfter int64,
	confirmedMinutesAfter int64,
	sequence uint64,
	createdAt time.Time,
) (Operation, error) {
	x := Operation{
		operationID:           operationID,
		coordinatorID:         coordinatorID,
		volunteerID:           volunteerID,
		description:           description,
		operationType:         operationType,
		source:                source,
		coinDelta:             coinDelta,
		minuteDelta:           minuteDelta,
		coinBalanceAfter:      coinBalanceAfter,
		confirmedMinutesAfter: confirmedMinutesAfter,
		sequence:              sequence,
		createdAt:             createdAt,
	}
	if err := x.Validate(); err != nil {
		return Operation{}, err
	}

	return x, nil
}
