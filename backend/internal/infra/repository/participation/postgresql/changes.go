package repository_participation_postgresql

import (
	"context"
	"fmt"
	"time"

	model_actor_kind "github.com/Roum1212/hahaton/internal/domain/model/actor-kind"
	model_operation_id "github.com/Roum1212/hahaton/internal/domain/model/operation-id"
	model_participation_change "github.com/Roum1212/hahaton/internal/domain/model/participation-change"
	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
	model_participation_status "github.com/Roum1212/hahaton/internal/domain/model/participation-status"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

func (x Participations) AppendChange(ctx context.Context, value model_participation_change.ParticipationChange) error {
	values := map[string]any{
		idColumn:                 value.GetChangeID().String(),
		participationIDColumn:    value.GetParticipationID().String(),
		previousStatusColumn:     "",
		statusColumn:             value.GetStatus().String(),
		actorIDColumn:            "",
		actorKindColumn:          value.GetActorKind().String(),
		cancellationReasonColumn: value.GetCancellationReason(),
		rejectionReasonColumn:    value.GetRejectionReason(),
		createdAtColumn:          value.GetCreatedAt(),
	}
	if previous, ok := value.GetPreviousStatus().Get(); ok {
		values[previousStatusColumn] = previous.String()
	}

	if !value.GetActorID().IsZero() {
		values[actorIDColumn] = value.GetActorID().String()
	}

	sql, arguments, err := x.queries.Insert("participation_changes").SetMap(values).ToSql()
	if err != nil {
		return err
	}

	_, err = x.database.Exec(ctx, sql, arguments...)

	return err
}

func (x Participations) ListChanges(ctx context.Context) ([]model_participation_change.ParticipationChange, error) {
	sql, arguments, err := x.queries.Select(
		idColumn,
		participationIDColumn,
		previousStatusColumn,
		statusColumn,
		actorIDColumn,
		actorKindColumn,
		cancellationReasonColumn,
		rejectionReasonColumn,
		sequenceColumn,
		createdAtColumn,
	).From("participation_changes").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := x.database.Query(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := make([]model_participation_change.ParticipationChange, 0)

	for rows.Next() {
		var record struct {
			changeID           string
			participationID    string
			previousStatus     string
			status             string
			actorID            string
			actorKind          string
			cancellationReason string
			rejectionReason    string
			sequence           uint64
			createdAt          time.Time
		}
		if err := rows.Scan(
			&record.changeID,
			&record.participationID,
			&record.previousStatus,
			&record.status,
			&record.actorID,
			&record.actorKind,
			&record.cancellationReason,
			&record.rejectionReason,
			&record.sequence,
			&record.createdAt,
		); err != nil {
			return nil, err
		}

		changeID, err := model_operation_id.NewOperationID(record.changeID)
		if err != nil {
			return nil, fmt.Errorf("restore change id: %w", err)
		}

		participationID, err := model_participation_id.NewParticipationID(record.participationID)
		if err != nil {
			return nil, fmt.Errorf("restore participation id: %w", err)
		}

		var previousStatus utils_value.Optional[model_participation_status.ParticipationStatus]
		if len(record.previousStatus) != 0 {
			previous, err := model_participation_status.NewParticipationStatus(record.previousStatus)
			if err != nil {
				return nil, fmt.Errorf("restore previous: %w", err)
			}

			previousStatus = utils_value.NewOptional(previous)
		}

		status, err := model_participation_status.NewParticipationStatus(record.status)
		if err != nil {
			return nil, fmt.Errorf("restore status: %w", err)
		}

		var actorID model_user_id.UserID
		if len(record.actorID) != 0 {
			actorID, err = model_user_id.NewUserID(record.actorID)
			if err != nil {
				return nil, fmt.Errorf("restore actor id: %w", err)
			}
		}

		actorKind, err := model_actor_kind.NewActorKind(record.actorKind)
		if err != nil {
			return nil, fmt.Errorf("restore actor kind: %w", err)
		}

		value, err := model_participation_change.NewParticipationChange(
			changeID,
			participationID,
			previousStatus,
			status,
			actorID,
			actorKind,
			record.cancellationReason,
			record.rejectionReason,
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
