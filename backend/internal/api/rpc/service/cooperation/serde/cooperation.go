package rpc_service_cooperation_serde

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/protobuf/types/known/timestamppb"

	cooperation_v1 "github.com/Roum1212/hahaton/genproto/cooperation/v1"
	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_cooperation_invitation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation-invitation"
	model_cooperation_status "github.com/Roum1212/hahaton/internal/domain/model/cooperation-status"
)

func AggregateToPBCooperation(value *aggregate_cooperation.Cooperation) *cooperation_v1.Cooperation {
	result := &cooperation_v1.Cooperation{
		CoordinatorId: value.GetCoordinatorID().String(),
		VolunteerId:   value.GetVolunteerID().String(),
		Status:        StatusToPB(value.GetStatus()),
		CreatedAt:     timestamppb.New(value.GetCreatedAt()),
		UpdatedAt:     timestamppb.New(value.GetUpdatedAt()),
	}
	if value.GetStatus() == model_cooperation_status.Rejected {
		result.RetryAfter = timestamppb.New(value.GetRetryAfter())
	}

	return result
}

func InvitationToPB(
	value *aggregate_cooperation_invitation.CooperationInvitation,
	frontendURL string,
) *cooperation_v1.CooperationInvitation {
	result := &cooperation_v1.CooperationInvitation{
		Id:            value.GetInvitationID().String(),
		CoordinatorId: value.GetCoordinatorID().String(),
		InvitationUrl: strings.TrimRight(frontendURL, "/") + "/invite/" + value.GetToken(),
		CreatedAt:     timestamppb.New(value.GetCreatedAt()),
	}
	if !value.GetExpiresAt().IsZero() {
		result.ExpiresAt = timestamppb.New(value.GetExpiresAt())
	}

	return result
}

func StatusFromPB(value cooperation_v1.CooperationStatus) (model_cooperation_status.CooperationStatus, error) {
	switch value {
	case cooperation_v1.CooperationStatus_COOPERATION_STATUS_PENDING:
		return model_cooperation_status.NewCooperationStatus(model_cooperation_status.Pending.String())
	case cooperation_v1.CooperationStatus_COOPERATION_STATUS_ACTIVE:
		return model_cooperation_status.NewCooperationStatus(model_cooperation_status.Active.String())
	case cooperation_v1.CooperationStatus_COOPERATION_STATUS_REJECTED:
		return model_cooperation_status.NewCooperationStatus(model_cooperation_status.Rejected.String())
	case cooperation_v1.CooperationStatus_COOPERATION_STATUS_TERMINATED:
		return model_cooperation_status.NewCooperationStatus(model_cooperation_status.Terminated.String())
	case cooperation_v1.CooperationStatus_COOPERATION_STATUS_UNSPECIFIED:
		return "", errors.New("unspecified cooperation status")
	default:
		return "", fmt.Errorf("unknown cooperation status: %d", value)
	}
}

func StatusToPB(value model_cooperation_status.CooperationStatus) cooperation_v1.CooperationStatus {
	values := map[model_cooperation_status.CooperationStatus]cooperation_v1.CooperationStatus{
		model_cooperation_status.Pending:    cooperation_v1.CooperationStatus_COOPERATION_STATUS_PENDING,
		model_cooperation_status.Active:     cooperation_v1.CooperationStatus_COOPERATION_STATUS_ACTIVE,
		model_cooperation_status.Rejected:   cooperation_v1.CooperationStatus_COOPERATION_STATUS_REJECTED,
		model_cooperation_status.Terminated: cooperation_v1.CooperationStatus_COOPERATION_STATUS_TERMINATED,
	}

	return values[value]
}
