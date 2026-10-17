package model_participation_change

import (
	"errors"
	"time"

	model_actor_kind "github.com/Roum1212/hahaton/internal/domain/model/actor-kind"
	model_operation_id "github.com/Roum1212/hahaton/internal/domain/model/operation-id"
	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
	model_participation_status "github.com/Roum1212/hahaton/internal/domain/model/participation-status"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

type ParticipationChange struct {
	changeID           model_operation_id.OperationID
	participationID    model_participation_id.ParticipationID
	previousStatus     utils_value.Optional[model_participation_status.ParticipationStatus]
	status             model_participation_status.ParticipationStatus
	actorID            model_user_id.UserID
	actorKind          model_actor_kind.ActorKind
	cancellationReason string
	rejectionReason    string
	sequence           uint64
	createdAt          time.Time
}

func (x ParticipationChange) GetActorID() model_user_id.UserID {
	return x.actorID
}

func (x ParticipationChange) GetActorKind() model_actor_kind.ActorKind {
	return x.actorKind
}

func (x ParticipationChange) GetCancellationReason() string {
	return x.cancellationReason
}

func (x ParticipationChange) GetChangeID() model_operation_id.OperationID {
	return x.changeID
}

func (x ParticipationChange) GetCreatedAt() time.Time {
	return x.createdAt
}

func (x ParticipationChange) GetParticipationID() model_participation_id.ParticipationID {
	return x.participationID
}

func (x ParticipationChange) GetPreviousStatus() utils_value.Optional[model_participation_status.ParticipationStatus] {
	return x.previousStatus
}

func (x ParticipationChange) GetRejectionReason() string {
	return x.rejectionReason
}

func (x ParticipationChange) GetSequence() uint64 {
	return x.sequence
}

func (x ParticipationChange) GetStatus() model_participation_status.ParticipationStatus {
	return x.status
}

func (x ParticipationChange) Validate() error {
	if err := x.changeID.Validate(); err != nil {
		return err
	}

	if err := x.participationID.Validate(); err != nil {
		return err
	}

	if previous, ok := x.previousStatus.Get(); ok {
		if err := previous.Validate(); err != nil {
			return err
		}
	}

	if err := x.status.Validate(); err != nil {
		return err
	}

	if err := x.actorKind.Validate(); err != nil {
		return err
	}

	if x.actorKind == model_actor_kind.System {
		if !x.actorID.IsZero() {
			return errors.New("system change must not have an actor id")
		}
	} else if err := x.actorID.Validate(); err != nil {
		return err
	}

	if x.createdAt.IsZero() {
		return errors.New("participation change time is required")
	}

	return nil
}

func NewParticipationChange(
	changeID model_operation_id.OperationID,
	participationID model_participation_id.ParticipationID,
	previousStatus utils_value.Optional[model_participation_status.ParticipationStatus],
	status model_participation_status.ParticipationStatus,
	actorID model_user_id.UserID,
	actorKind model_actor_kind.ActorKind,
	cancellationReason string,
	rejectionReason string,
	sequence uint64,
	createdAt time.Time,
) (ParticipationChange, error) {
	x := ParticipationChange{
		changeID:           changeID,
		participationID:    participationID,
		previousStatus:     previousStatus,
		status:             status,
		actorID:            actorID,
		actorKind:          actorKind,
		cancellationReason: cancellationReason,
		rejectionReason:    rejectionReason,
		sequence:           sequence,
		createdAt:          createdAt,
	}
	if err := x.Validate(); err != nil {
		return ParticipationChange{}, err
	}

	return x, nil
}
