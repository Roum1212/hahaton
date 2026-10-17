package repository_participation_postgresql

import (
	"fmt"
	"time"

	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_event_result "github.com/Roum1212/hahaton/internal/domain/model/event-result"
	model_operation_id "github.com/Roum1212/hahaton/internal/domain/model/operation-id"
	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
	model_participation_status "github.com/Roum1212/hahaton/internal/domain/model/participation-status"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

type row struct {
	participationID      string
	eventID              string
	coordinatorID        string
	volunteerID          string
	status               string
	cancellationReason   string
	rejectionReason      string
	awardedMinutes       int64
	awardedCoins         int64
	accrualOperationID   string
	hasChatAccess        bool
	visibleAfterSequence uint64
	createdAt            time.Time
	updatedAt            time.Time
}

func (x *row) scan(source interface{ Scan(dest ...any) error }) error {
	return source.Scan(
		&x.participationID,
		&x.eventID,
		&x.coordinatorID,
		&x.volunteerID,
		&x.status,
		&x.cancellationReason,
		&x.rejectionReason,
		&x.awardedMinutes,
		&x.awardedCoins,
		&x.accrualOperationID,
		&x.hasChatAccess,
		&x.visibleAfterSequence,
		&x.createdAt,
		&x.updatedAt,
	)
}

func (x *row) toAggregate() (*aggregate_participation.Participation, error) {
	participationID, err := model_participation_id.NewParticipationID(x.participationID)
	if err != nil {
		return nil, fmt.Errorf("restore participation id: %w", err)
	}

	eventID, err := model_event_id.NewEventID(x.eventID)
	if err != nil {
		return nil, fmt.Errorf("restore event id: %w", err)
	}

	coordinatorID, err := model_user_id.NewUserID(x.coordinatorID)
	if err != nil {
		return nil, fmt.Errorf("restore coordinator id: %w", err)
	}

	volunteerID, err := model_user_id.NewUserID(x.volunteerID)
	if err != nil {
		return nil, fmt.Errorf("restore volunteer id: %w", err)
	}

	status, err := model_participation_status.NewParticipationStatus(x.status)
	if err != nil {
		return nil, fmt.Errorf("restore status: %w", err)
	}

	var result utils_value.Optional[model_event_result.EventResult]

	if status == model_participation_status.Confirmed || status == model_participation_status.NotAttended {
		value, err := model_event_result.NewEventResult(
			participationID, x.awardedMinutes, status == model_participation_status.NotAttended,
		)
		if err != nil {
			return nil, fmt.Errorf("restore value: %w", err)
		}

		result = utils_value.NewOptional(value)
	}

	var accrualOperationID model_operation_id.OperationID
	if len(x.accrualOperationID) != 0 {
		value, err := model_operation_id.NewOperationID(x.accrualOperationID)
		if err != nil {
			return nil, fmt.Errorf("restore value: %w", err)
		}

		accrualOperationID = value
	}

	return aggregate_participation.NewParticipationFromFields(
		participationID,
		eventID,
		coordinatorID,
		volunteerID,
		status,
		x.cancellationReason,
		x.rejectionReason,
		result,
		x.awardedCoins,
		accrualOperationID,
		x.hasChatAccess,
		x.visibleAfterSequence,
		x.createdAt,
		x.updatedAt,
	), nil
}
