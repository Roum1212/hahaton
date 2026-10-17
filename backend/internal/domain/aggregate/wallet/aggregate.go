package aggregate_wallet

import (
	"errors"
	"math"
	"time"

	aggregate_wallet_command "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet/command"
	model_event_result "github.com/Roum1212/hahaton/internal/domain/model/event-result"
	model_operation_id "github.com/Roum1212/hahaton/internal/domain/model/operation-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	model_wallet_operation "github.com/Roum1212/hahaton/internal/domain/model/wallet-operation"
	model_wallet_operation_source "github.com/Roum1212/hahaton/internal/domain/model/wallet-operation-source"
	model_wallet_operation_type "github.com/Roum1212/hahaton/internal/domain/model/wallet-operation-type"
)

var (
	ErrAlreadyOpened       = errors.New("wallet is already opened")
	ErrInsufficientBalance = errors.New("insufficient coin balance")
	ErrOverflow            = errors.New("wallet totals overflow")
)

type Wallet struct {
	userID           model_user_id.UserID
	coinBalance      int64
	confirmedMinutes int64
	lastSequence     uint64
	createdAt        time.Time
	updatedAt        time.Time
}

func (x *Wallet) GetCoinBalance() int64 {
	return x.coinBalance
}

func (x *Wallet) GetConfirmedMinutes() int64 {
	return x.confirmedMinutes
}

func (x *Wallet) GetCreatedAt() time.Time {
	return x.createdAt
}

func (x *Wallet) GetLastSequence() uint64 {
	return x.lastSequence
}

func (x *Wallet) GetUpdatedAt() time.Time {
	return x.updatedAt
}

func (x *Wallet) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x *Wallet) Accrue(command aggregate_wallet_command.Accrue) (model_wallet_operation.Operation, error) {
	result := command.GetResult()

	minutes := result.GetMinutes()
	if minutes <= 0 || minutes%model_event_result.MinutesPerStep != 0 ||
		minutes/model_event_result.MinutesPerStep > math.MaxInt64/model_event_result.CoinsPerStep {
		return model_wallet_operation.Operation{}, ErrOverflow
	}

	source, err := model_wallet_operation_source.NewOperationSourceFromParticipationID(result.GetParticipationID())
	if err != nil {
		return model_wallet_operation.Operation{}, err
	}

	return x.apply(
		command.GetOperationID(),
		command.GetCoordinatorID(),
		source,
		command.GetDescription(),
		model_wallet_operation_type.ParticipationAccrual,
		minutes/model_event_result.MinutesPerStep*model_event_result.CoinsPerStep,
		minutes,
		command.GetNow(),
	)
}

func (x *Wallet) Open(command aggregate_wallet_command.Open) error {
	if !x.createdAt.IsZero() {
		return ErrAlreadyOpened
	}

	x.createdAt = command.GetNow()
	x.updatedAt = command.GetNow()

	return nil
}

func (x *Wallet) Purchase(command aggregate_wallet_command.Purchase) (model_wallet_operation.Operation, error) {
	if command.GetCoins() <= 0 || x.coinBalance < command.GetCoins() {
		return model_wallet_operation.Operation{}, ErrInsufficientBalance
	}

	source, err := model_wallet_operation_source.NewOperationSourceFromRewardOrderID(command.GetSourceID())
	if err != nil {
		return model_wallet_operation.Operation{}, err
	}

	return x.apply(
		command.GetOperationID(),
		command.GetCoordinatorID(),
		source,
		command.GetDescription(),
		model_wallet_operation_type.RewardPurchase,
		-command.GetCoins(),
		0,
		command.GetNow(),
	)
}

func (x *Wallet) Refund(command aggregate_wallet_command.Refund) (model_wallet_operation.Operation, error) {
	source, err := model_wallet_operation_source.NewOperationSourceFromRewardOrderID(command.GetSourceID())
	if err != nil {
		return model_wallet_operation.Operation{}, err
	}

	return x.apply(
		command.GetOperationID(),
		command.GetCoordinatorID(),
		source,
		command.GetDescription(),
		model_wallet_operation_type.RewardRefund,
		command.GetCoins(),
		0,
		command.GetNow(),
	)
}

func (x *Wallet) apply(
	operationID model_operation_id.OperationID,
	coordinatorID model_user_id.UserID,
	source model_wallet_operation_source.OperationSource,
	description string,
	operationType model_wallet_operation_type.OperationType,
	coins int64,
	minutes int64,
	now time.Time,
) (model_wallet_operation.Operation, error) {
	if coins > 0 && x.coinBalance > math.MaxInt64-coins {
		return model_wallet_operation.Operation{}, ErrOverflow
	}

	if minutes > 0 && x.confirmedMinutes > math.MaxInt64-minutes {
		return model_wallet_operation.Operation{}, ErrOverflow
	}

	if x.lastSequence >= math.MaxInt64 {
		return model_wallet_operation.Operation{}, ErrOverflow
	}

	operation, err := model_wallet_operation.NewOperation(
		operationID,
		coordinatorID,
		x.userID,
		description,
		operationType,
		source,
		coins,
		minutes,
		x.coinBalance+coins,
		x.confirmedMinutes+minutes,
		x.lastSequence+1,
		now,
	)
	if err != nil {
		return model_wallet_operation.Operation{}, err
	}

	x.coinBalance = operation.GetCoinBalanceAfter()
	x.confirmedMinutes = operation.GetConfirmedMinutesAfter()
	x.lastSequence = operation.GetSequence()
	x.updatedAt = now

	return operation, nil
}

func NewWallet(userID model_user_id.UserID) *Wallet {
	return &Wallet{userID: userID}
}

func NewWalletFromFields(
	userID model_user_id.UserID,
	coinBalance int64,
	confirmedMinutes int64,
	lastSequence uint64,
	createdAt time.Time,
	updatedAt time.Time,
) *Wallet {
	return &Wallet{
		userID:           userID,
		coinBalance:      coinBalance,
		confirmedMinutes: confirmedMinutes,
		lastSequence:     lastSequence,
		createdAt:        createdAt,
		updatedAt:        updatedAt,
	}
}
