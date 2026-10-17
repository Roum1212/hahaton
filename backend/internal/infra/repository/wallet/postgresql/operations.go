package repository_wallet_postgresql

import (
	"context"
	"fmt"
	"time"

	model_operation_id "github.com/Roum1212/hahaton/internal/domain/model/operation-id"
	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
	model_reward_order_id "github.com/Roum1212/hahaton/internal/domain/model/reward-order-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	model_wallet_operation "github.com/Roum1212/hahaton/internal/domain/model/wallet-operation"
	model_wallet_operation_source "github.com/Roum1212/hahaton/internal/domain/model/wallet-operation-source"
	model_wallet_operation_type "github.com/Roum1212/hahaton/internal/domain/model/wallet-operation-type"
)

func (x Wallets) AppendOperation(ctx context.Context, value model_wallet_operation.Operation) error {
	values := map[string]any{
		idColumn:                    value.GetOperationID().String(),
		coordinatorIDColumn:         value.GetCoordinatorID().String(),
		volunteerIDColumn:           value.GetVolunteerID().String(),
		descriptionColumn:           value.GetDescription(),
		operationTypeColumn:         value.GetOperationType().Int32(),
		sourceIDColumn:              value.GetSource().String(),
		coinDeltaColumn:             value.GetCoinDelta(),
		minuteDeltaColumn:           value.GetMinuteDelta(),
		coinBalanceAfterColumn:      value.GetCoinBalanceAfter(),
		confirmedMinutesAfterColumn: value.GetConfirmedMinutesAfter(),
		sequenceColumn:              value.GetSequence(),
		createdAtColumn:             value.GetCreatedAt(),
	}

	sql, arguments, err := x.queries.Insert("volunteer_operations").SetMap(values).ToSql()
	if err != nil {
		return err
	}

	_, err = x.database.Exec(ctx, sql, arguments...)

	return err
}

func (x Wallets) ListOperations(ctx context.Context) ([]model_wallet_operation.Operation, error) {
	sql, arguments, err := x.queries.Select(
		idColumn,
		coordinatorIDColumn,
		volunteerIDColumn,
		descriptionColumn,
		operationTypeColumn,
		sourceIDColumn,
		coinDeltaColumn,
		minuteDeltaColumn,
		coinBalanceAfterColumn,
		confirmedMinutesAfterColumn,
		sequenceColumn,
		createdAtColumn,
	).From("volunteer_operations").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := x.database.Query(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := make([]model_wallet_operation.Operation, 0)

	for rows.Next() {
		var record struct {
			operationID           string
			coordinatorID         string
			volunteerID           string
			description           string
			operationType         int32
			sourceID              string
			coinDelta             int64
			minuteDelta           int64
			coinBalanceAfter      int64
			confirmedMinutesAfter int64
			sequence              uint64
			createdAt             time.Time
		}
		if err := rows.Scan(
			&record.operationID,
			&record.coordinatorID,
			&record.volunteerID,
			&record.description,
			&record.operationType,
			&record.sourceID,
			&record.coinDelta,
			&record.minuteDelta,
			&record.coinBalanceAfter,
			&record.confirmedMinutesAfter,
			&record.sequence,
			&record.createdAt,
		); err != nil {
			return nil, err
		}

		operationID, err := model_operation_id.NewOperationID(record.operationID)
		if err != nil {
			return nil, fmt.Errorf("restore operation id: %w", err)
		}

		coordinatorID, err := model_user_id.NewUserID(record.coordinatorID)
		if err != nil {
			return nil, fmt.Errorf("restore coordinator id: %w", err)
		}

		volunteerID, err := model_user_id.NewUserID(record.volunteerID)
		if err != nil {
			return nil, fmt.Errorf("restore volunteer id: %w", err)
		}

		operationType, err := model_wallet_operation_type.NewOperationType(record.operationType)
		if err != nil {
			return nil, fmt.Errorf("restore operation type: %w", err)
		}

		source, err := operationSourceFromFields(operationType, record.sourceID)
		if err != nil {
			return nil, err
		}

		value, err := model_wallet_operation.NewOperation(
			operationID,
			coordinatorID,
			volunteerID,
			record.description,
			operationType,
			source,
			record.coinDelta,
			record.minuteDelta,
			record.coinBalanceAfter,
			record.confirmedMinutesAfter,
			record.sequence,
			record.createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("restore value: %w", err)
		}

		values = append(values, value)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return values, nil
}

func operationSourceFromFields(
	operationType model_wallet_operation_type.OperationType,
	sourceID string,
) (model_wallet_operation_source.OperationSource, error) {
	if operationType == model_wallet_operation_type.ParticipationAccrual {
		id, err := model_participation_id.NewParticipationID(sourceID)
		if err != nil {
			return model_wallet_operation_source.OperationSource{}, fmt.Errorf("restore id: %w", err)
		}

		return model_wallet_operation_source.NewOperationSourceFromParticipationID(id)
	}

	id, err := model_reward_order_id.NewRewardOrderID(sourceID)
	if err != nil {
		return model_wallet_operation_source.OperationSource{}, fmt.Errorf("restore id: %w", err)
	}

	return model_wallet_operation_source.NewOperationSourceFromRewardOrderID(id)
}
