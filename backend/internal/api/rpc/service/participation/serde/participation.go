package rpc_service_participation_serde

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	participation_v1 "github.com/Roum1212/hahaton/genproto/participation/v1"
	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	model_actor_kind "github.com/Roum1212/hahaton/internal/domain/model/actor-kind"
	model_participation_change "github.com/Roum1212/hahaton/internal/domain/model/participation-change"
	model_participation_status "github.com/Roum1212/hahaton/internal/domain/model/participation-status"
)

func AggregateToPBParticipation(value *aggregate_participation.Participation) *participation_v1.Participation {
	result := &participation_v1.Participation{
		Id:                 value.GetParticipationID().String(),
		EventId:            value.GetEventID().String(),
		VolunteerId:        value.GetVolunteerID().String(),
		Status:             StatusToPB(value.GetStatus()),
		CancellationReason: CancellationReasonToPB(value.GetCancellationReason()),
		CreatedAt:          timestamppb.New(value.GetCreatedAt()),
		UpdatedAt:          timestamppb.New(value.GetUpdatedAt()),
	}
	if len(value.GetRejectionReason()) != 0 {
		reason := value.GetRejectionReason()
		result.RejectionReason = &reason
	}

	if value.GetStatus() == model_participation_status.Confirmed {
		coins, minutes, id := value.GetAwardedCoins(), value.GetAwardedMinutes(), value.GetAccrualOperationID().String()
		result.AwardedCoins = &coins
		result.AwardedMinutes = &minutes
		result.AccrualOperationId = &id
	}

	return result
}

func CancellationReasonToPB(value string) participation_v1.ParticipationCancellationReason {
	switch value {
	case "volunteer_canceled":
		return participation_v1.ParticipationCancellationReason_PARTICIPATION_CANCELLATION_REASON_VOLUNTEER_CANCELLED
	case "event_canceled":
		return participation_v1.ParticipationCancellationReason_PARTICIPATION_CANCELLATION_REASON_EVENT_CANCELLED
	case "cooperation_terminated":
		return participation_v1.ParticipationCancellationReason_PARTICIPATION_CANCELLATION_REASON_COOPERATION_TERMINATED
	default:
		return participation_v1.ParticipationCancellationReason_PARTICIPATION_CANCELLATION_REASON_UNSPECIFIED
	}
}

func ChangeToPB(value model_participation_change.ParticipationChange) *participation_v1.ParticipationChange {
	status := value.GetStatus()

	result := &participation_v1.ParticipationChange{
		Sequence:           value.GetSequence(),
		Status:             StatusToPB(status),
		CancellationReason: CancellationReasonToPB(value.GetCancellationReason()),
		CreatedAt:          timestamppb.New(value.GetCreatedAt()),
	}
	if len(value.GetRejectionReason()) != 0 {
		reason := value.GetRejectionReason()
		result.RejectionReason = &reason
	}

	switch value.GetActorKind() {
	case model_actor_kind.Coordinator:
		result.Actor = &participation_v1.ParticipationChange_CoordinatorId{CoordinatorId: value.GetActorID().String()}
	case model_actor_kind.Volunteer:
		result.Actor = &participation_v1.ParticipationChange_VolunteerId{VolunteerId: value.GetActorID().String()}
	case model_actor_kind.System:
	}

	return result
}

func StatusFromPB(value participation_v1.ParticipationStatus) (model_participation_status.ParticipationStatus, error) {
	switch value {
	case participation_v1.ParticipationStatus_PARTICIPATION_STATUS_PENDING:
		return model_participation_status.NewParticipationStatus(model_participation_status.Pending.String())
	case participation_v1.ParticipationStatus_PARTICIPATION_STATUS_APPROVED:
		return model_participation_status.NewParticipationStatus(model_participation_status.Approved.String())
	case participation_v1.ParticipationStatus_PARTICIPATION_STATUS_REJECTED:
		return model_participation_status.NewParticipationStatus(model_participation_status.Rejected.String())
	case participation_v1.ParticipationStatus_PARTICIPATION_STATUS_CANCELLED:
		return model_participation_status.NewParticipationStatus(model_participation_status.Canceled.String())
	case participation_v1.ParticipationStatus_PARTICIPATION_STATUS_CONFIRMED:
		return model_participation_status.NewParticipationStatus(model_participation_status.Confirmed.String())
	case participation_v1.ParticipationStatus_PARTICIPATION_STATUS_NOT_ATTENDED:
		return model_participation_status.NewParticipationStatus(model_participation_status.NotAttended.String())
	case participation_v1.ParticipationStatus_PARTICIPATION_STATUS_CLOSED:
		return model_participation_status.NewParticipationStatus(model_participation_status.Closed.String())
	case participation_v1.ParticipationStatus_PARTICIPATION_STATUS_UNSPECIFIED:
		return "", errors.New("unspecified participation status")
	default:
		return "", fmt.Errorf("unknown participation status: %d", value)
	}
}

func StatusToPB(value model_participation_status.ParticipationStatus) participation_v1.ParticipationStatus {
	values := map[model_participation_status.ParticipationStatus]participation_v1.ParticipationStatus{
		model_participation_status.Pending:     participation_v1.ParticipationStatus_PARTICIPATION_STATUS_PENDING,
		model_participation_status.Approved:    participation_v1.ParticipationStatus_PARTICIPATION_STATUS_APPROVED,
		model_participation_status.Rejected:    participation_v1.ParticipationStatus_PARTICIPATION_STATUS_REJECTED,
		model_participation_status.Canceled:    participation_v1.ParticipationStatus_PARTICIPATION_STATUS_CANCELLED,
		model_participation_status.Confirmed:   participation_v1.ParticipationStatus_PARTICIPATION_STATUS_CONFIRMED,
		model_participation_status.NotAttended: participation_v1.ParticipationStatus_PARTICIPATION_STATUS_NOT_ATTENDED,
		model_participation_status.Closed:      participation_v1.ParticipationStatus_PARTICIPATION_STATUS_CLOSED,
	}

	return values[value]
}
