package aggregate_participation

import (
	"errors"
	"time"

	aggregate_participation_command "github.com/Roum1212/hahaton/internal/domain/aggregate/participation/command"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_event_result "github.com/Roum1212/hahaton/internal/domain/model/event-result"
	model_operation_id "github.com/Roum1212/hahaton/internal/domain/model/operation-id"
	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
	model_participation_status "github.com/Roum1212/hahaton/internal/domain/model/participation-status"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

var (
	ErrAlreadyExists = errors.New("participation already exists")
	ErrInvalidState  = errors.New("invalid participation state")
)

type Participation struct {
	participationID      model_participation_id.ParticipationID
	eventID              model_event_id.EventID
	coordinatorID        model_user_id.UserID
	volunteerID          model_user_id.UserID
	status               model_participation_status.ParticipationStatus
	cancellationReason   string
	rejectionReason      string
	result               utils_value.Optional[model_event_result.EventResult]
	awardedCoins         int64
	accrualOperationID   model_operation_id.OperationID
	hasChatAccess        bool
	visibleAfterSequence uint64
	createdAt            time.Time
	updatedAt            time.Time
}

func (x *Participation) GetAccrualOperationID() model_operation_id.OperationID {
	return x.accrualOperationID
}

func (x *Participation) GetAwardedCoins() int64 {
	return x.awardedCoins
}

func (x *Participation) GetAwardedMinutes() int64 {
	if result, ok := x.result.Get(); ok {
		return result.GetMinutes()
	}

	return 0
}

func (x *Participation) GetCancellationReason() string {
	return x.cancellationReason
}

func (x *Participation) GetCoordinatorID() model_user_id.UserID {
	return x.coordinatorID
}

func (x *Participation) GetCreatedAt() time.Time {
	return x.createdAt
}

func (x *Participation) GetEventID() model_event_id.EventID {
	return x.eventID
}

func (x *Participation) GetHasChatAccess() bool {
	return x.hasChatAccess
}

func (x *Participation) GetParticipationID() model_participation_id.ParticipationID {
	return x.participationID
}

func (x *Participation) GetRejectionReason() string {
	return x.rejectionReason
}

func (x *Participation) GetResult() utils_value.Optional[model_event_result.EventResult] {
	return x.result
}

func (x *Participation) GetStatus() model_participation_status.ParticipationStatus {
	return x.status
}

func (x *Participation) GetUpdatedAt() time.Time {
	return x.updatedAt
}

func (x *Participation) GetVisibleAfterSequence() uint64 {
	return x.visibleAfterSequence
}

func (x *Participation) GetVolunteerID() model_user_id.UserID {
	return x.volunteerID
}

func (x *Participation) Approve(command aggregate_participation_command.Approve) error {
	if x.status != model_participation_status.Pending {
		return ErrInvalidState
	}

	x.status = model_participation_status.Approved
	x.hasChatAccess = true
	x.visibleAfterSequence = command.GetVisibleAfterSequence()
	x.updatedAt = command.GetNow()

	return nil
}

func (x *Participation) Cancel(command aggregate_participation_command.Cancel) error {
	if x.status != model_participation_status.Pending && x.status != model_participation_status.Approved {
		return ErrInvalidState
	}

	x.status = model_participation_status.Canceled
	x.cancellationReason = command.GetReason()
	x.hasChatAccess = false
	x.updatedAt = command.GetNow()

	return nil
}

func (x *Participation) Close(command aggregate_participation_command.Close) error {
	if x.status != model_participation_status.Pending {
		return ErrInvalidState
	}

	x.status = model_participation_status.Closed
	x.updatedAt = command.GetNow()

	return nil
}

func (x *Participation) Complete(command aggregate_participation_command.Complete) error {
	result := command.GetResult()
	if x.status != model_participation_status.Approved &&
		(!command.IsForced() || x.status != model_participation_status.Pending) {
		return ErrInvalidState
	}

	if result.GetParticipationID() != x.participationID {
		return ErrInvalidState
	}

	x.result = utils_value.NewOptional(result)

	x.status = model_participation_status.NotAttended
	if !result.IsAbsent() {
		x.status = model_participation_status.Confirmed
		x.awardedCoins = command.GetCoins()
		x.accrualOperationID = command.GetOperationID()
	}

	x.updatedAt = command.GetNow()

	return nil
}

func (x *Participation) Reject(command aggregate_participation_command.Reject) error {
	if x.status != model_participation_status.Pending && x.status != model_participation_status.Approved {
		return ErrInvalidState
	}

	x.status = model_participation_status.Rejected
	x.rejectionReason = command.GetReason()
	x.hasChatAccess = false
	x.updatedAt = command.GetNow()

	return nil
}

func (x *Participation) Request(command aggregate_participation_command.Request) error {
	if !x.createdAt.IsZero() {
		return ErrAlreadyExists
	}

	x.status = model_participation_status.Pending
	x.createdAt = command.GetNow()
	x.updatedAt = command.GetNow()

	return nil
}

func (x *Participation) Restore(command aggregate_participation_command.Restore) error {
	if x.status != model_participation_status.Canceled {
		return ErrInvalidState
	}

	x.status = model_participation_status.Pending
	x.cancellationReason = ""
	x.updatedAt = command.GetNow()

	return nil
}

func (x *Participation) RevokeChat(command aggregate_participation_command.RevokeChat) {
	x.hasChatAccess = false
	x.updatedAt = command.GetNow()
}

func NewParticipation(
	participationID model_participation_id.ParticipationID,
	eventID model_event_id.EventID,
	coordinatorID model_user_id.UserID,
	volunteerID model_user_id.UserID,
) *Participation {
	return &Participation{
		participationID: participationID,
		eventID:         eventID,
		coordinatorID:   coordinatorID,
		volunteerID:     volunteerID,
	}
}

func NewParticipationFromFields(
	participationID model_participation_id.ParticipationID,
	eventID model_event_id.EventID,
	coordinatorID model_user_id.UserID,
	volunteerID model_user_id.UserID,
	status model_participation_status.ParticipationStatus,
	cancellationReason string,
	rejectionReason string,
	result utils_value.Optional[model_event_result.EventResult],
	awardedCoins int64,
	accrualOperationID model_operation_id.OperationID,
	hasChatAccess bool,
	visibleAfterSequence uint64,
	createdAt time.Time,
	updatedAt time.Time,
) *Participation {
	return &Participation{
		participationID:      participationID,
		eventID:              eventID,
		coordinatorID:        coordinatorID,
		volunteerID:          volunteerID,
		status:               status,
		cancellationReason:   cancellationReason,
		rejectionReason:      rejectionReason,
		result:               result,
		awardedCoins:         awardedCoins,
		accrualOperationID:   accrualOperationID,
		hasChatAccess:        hasChatAccess,
		visibleAfterSequence: visibleAfterSequence,
		createdAt:            createdAt,
		updatedAt:            updatedAt,
	}
}
