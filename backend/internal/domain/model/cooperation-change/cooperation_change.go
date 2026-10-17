package model_cooperation_change

import (
	"errors"
	"time"

	model_cooperation_status "github.com/Roum1212/hahaton/internal/domain/model/cooperation-status"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type CooperationChange struct {
	coordinatorID model_user_id.UserID
	volunteerID   model_user_id.UserID
	status        model_cooperation_status.CooperationStatus
	actorID       model_user_id.UserID
	retryAfter    time.Time
	createdAt     time.Time
}

func (x CooperationChange) GetActorID() model_user_id.UserID {
	return x.actorID
}

func (x CooperationChange) GetCoordinatorID() model_user_id.UserID {
	return x.coordinatorID
}

func (x CooperationChange) GetCreatedAt() time.Time {
	return x.createdAt
}

func (x CooperationChange) GetRetryAfter() time.Time {
	return x.retryAfter
}

func (x CooperationChange) GetStatus() model_cooperation_status.CooperationStatus {
	return x.status
}

func (x CooperationChange) GetVolunteerID() model_user_id.UserID {
	return x.volunteerID
}

func (x CooperationChange) Validate() error {
	if err := x.coordinatorID.Validate(); err != nil {
		return err
	}

	if err := x.volunteerID.Validate(); err != nil {
		return err
	}

	if err := x.actorID.Validate(); err != nil {
		return err
	}

	if err := x.status.Validate(); err != nil {
		return err
	}

	if x.createdAt.IsZero() {
		return errors.New("cooperation change time is required")
	}

	return nil
}

func NewCooperationChange(
	coordinatorID model_user_id.UserID,
	volunteerID model_user_id.UserID,
	status model_cooperation_status.CooperationStatus,
	actorID model_user_id.UserID,
	retryAfter time.Time,
	createdAt time.Time,
) (CooperationChange, error) {
	x := CooperationChange{
		coordinatorID: coordinatorID,
		volunteerID:   volunteerID,
		status:        status,
		actorID:       actorID,
		retryAfter:    retryAfter,
		createdAt:     createdAt,
	}
	if err := x.Validate(); err != nil {
		return CooperationChange{}, err
	}

	return x, nil
}
