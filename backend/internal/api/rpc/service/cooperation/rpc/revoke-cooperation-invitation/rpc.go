package rpc_service_cooperation_rpc_revoke_cooperation_invitation

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	rpc_error "go.flida.dev/rpc-error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"

	cooperation_v1 "github.com/Roum1212/hahaton/genproto/cooperation/v1"
	service_revoke_cooperation_invitation "github.com/Roum1212/hahaton/internal/app/service/revoke-cooperation-invitation"
	aggregate_cooperation_invitation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation-invitation"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	model_invitation_id "github.com/Roum1212/hahaton/internal/domain/model/invitation-id"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_revoke_cooperation_invitation.Service
}

func (x RPC) RevokeCooperationInvitation(
	ctx context.Context,
	r *connect.Request[cooperation_v1.RevokeCooperationInvitationRequest],
) (*connect.Response[emptypb.Empty], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	invitationID, err := model_invitation_id.NewInvitationID(r.Msg.GetId())
	if err != nil {
		return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "id",
					Description: err.Error(),
				},
			},
		})
	}

	if err = x.service.RevokeCooperationInvitation(ctx, service_revoke_cooperation_invitation.NewRequest(
		session.GetUserID(),
		invitationID,
	)); err != nil {
		switch {
		case errors.Is(err, service_revoke_cooperation_invitation.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_revoke_cooperation_invitation.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		case errors.Is(err, aggregate_cooperation_invitation.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "cooperation.v1.CooperationInvitation",
				ResourceName: "id:" + invitationID.String(),
				Description:  "cooperation invitation not found",
			})
		case errors.Is(err, aggregate_coordinator.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "coordinator.v1.Coordinator",
				ResourceName: "id:" + session.GetUserID().String(),
				Description:  "coordinator not found",
			})
		default:
			return nil, rpc_error.NewInternalError()
		}
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func NewRPC(service service_revoke_cooperation_invitation.Service) RPC {
	return RPC{service: service}
}
