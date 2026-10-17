package aggregate_cooperation

import (
	"errors"
	"time"

	aggregate_cooperation_command "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation/command"
	model_cooperation_status "github.com/Roum1212/hahaton/internal/domain/model/cooperation-status"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

var (
	ErrAlreadyExists   = errors.New("cooperation already exists")
	ErrCooldown        = errors.New("cooperation cooldown has not elapsed")
	ErrInvalidState    = errors.New("invalid cooperation state")
	ErrSelfCooperation = errors.New("cannot cooperate with yourself")
)

type Cooperation struct {
	coordinatorID model_user_id.UserID
	volunteerID   model_user_id.UserID
	status        model_cooperation_status.CooperationStatus
	retryAfter    time.Time
	createdAt     time.Time
	updatedAt     time.Time
}

func (x *Cooperation) GetCoordinatorID() model_user_id.UserID {
	return x.coordinatorID
}

func (x *Cooperation) GetCreatedAt() time.Time {
	return x.createdAt
}

func (x *Cooperation) GetRetryAfter() time.Time {
	return x.retryAfter
}

func (x *Cooperation) GetStatus() model_cooperation_status.CooperationStatus {
	return x.status
}

func (x *Cooperation) GetUpdatedAt() time.Time {
	return x.updatedAt
}

func (x *Cooperation) GetVolunteerID() model_user_id.UserID {
	return x.volunteerID
}

func (x *Cooperation) AcceptInvitation(command aggregate_cooperation_command.AcceptInvitation) error {
	if x.coordinatorID == x.volunteerID {
		return ErrSelfCooperation
	}

	if x.status == model_cooperation_status.Active {
		return nil
	}

	if x.createdAt.IsZero() {
		x.createdAt = command.GetNow()
	}

	x.status = model_cooperation_status.Active
	x.retryAfter = time.Time{}
	x.updatedAt = command.GetNow()

	return nil
}

func (x *Cooperation) Approve(command aggregate_cooperation_command.Approve) error {
	if x.status != model_cooperation_status.Pending {
		return ErrInvalidState
	}

	x.status = model_cooperation_status.Active
	x.updatedAt = command.GetNow()

	return nil
}

func (x *Cooperation) Reject(command aggregate_cooperation_command.Reject) error {
	if x.status != model_cooperation_status.Pending {
		return ErrInvalidState
	}

	x.status = model_cooperation_status.Rejected
	x.retryAfter = command.GetRetryAfter()
	x.updatedAt = command.GetNow()

	return nil
}

func (x *Cooperation) Request(command aggregate_cooperation_command.Request) error {
	if x.coordinatorID == x.volunteerID {
		return ErrSelfCooperation
	}

	switch x.status {
	case model_cooperation_status.Pending:
		return ErrAlreadyExists
	case model_cooperation_status.Active:
		return ErrAlreadyExists
	case model_cooperation_status.Terminated:
	case model_cooperation_status.Rejected:
		if command.GetNow().Before(x.retryAfter) {
			return ErrCooldown
		}
	}

	if x.createdAt.IsZero() {
		x.createdAt = command.GetNow()
	}

	x.status = model_cooperation_status.Pending
	x.retryAfter = time.Time{}
	x.updatedAt = command.GetNow()

	return nil
}

func (x *Cooperation) Terminate(command aggregate_cooperation_command.Terminate) error {
	if x.status != model_cooperation_status.Active {
		return ErrInvalidState
	}

	x.status = model_cooperation_status.Terminated
	x.updatedAt = command.GetNow()

	return nil
}

func NewCooperation(coordinatorID, volunteerID model_user_id.UserID) *Cooperation {
	return &Cooperation{
		coordinatorID: coordinatorID,
		volunteerID:   volunteerID,
	}
}

func NewCooperationFromFields(
	coordinatorID model_user_id.UserID,
	volunteerID model_user_id.UserID,
	status model_cooperation_status.CooperationStatus,
	retryAfter time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) *Cooperation {
	return &Cooperation{
		coordinatorID: coordinatorID,
		volunteerID:   volunteerID,
		status:        status,
		retryAfter:    retryAfter,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
}
