package rpc_service_cooperation_rpc_accept_cooperation_invitation

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	rpc_error "go.flida.dev/rpc-error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	cooperation_v1 "github.com/Roum1212/hahaton/genproto/cooperation/v1"
	rpc_service_cooperation_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/cooperation/serde"
	service_accept_cooperation_invitation "github.com/Roum1212/hahaton/internal/app/service/accept-cooperation-invitation"
	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_cooperation_invitation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation-invitation"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_accept_cooperation_invitation.Service
}

func (x RPC) AcceptCooperationInvitation(
	ctx context.Context,
	r *connect.Request[cooperation_v1.AcceptCooperationInvitationRequest],
) (*connect.Response[cooperation_v1.AcceptCooperationInvitationResponse], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	token := r.Msg.GetToken()

	response, err := x.service.AcceptCooperationInvitation(ctx, service_accept_cooperation_invitation.NewRequest(
		session.GetUserID(),
		token,
	))
	if err != nil {
		switch {
		case errors.Is(err, service_accept_cooperation_invitation.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_accept_cooperation_invitation.ErrFailedPrecondition):
			return nil, rpc_error.NewFailedPreconditionError(nil)
		case errors.Is(err, aggregate_cooperation_invitation.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "cooperation.v1.CooperationInvitation",
				ResourceName: "",
				Description:  "cooperation invitation not found",
			})
		case errors.Is(err, aggregate_cooperation.ErrAlreadyExists):
			return nil, rpc_error.NewAlreadyExistsError(nil)
		case errors.Is(err, aggregate_cooperation.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "cooperation.v1.Cooperation",
				ResourceName: "",
				Description:  "cooperation not found",
			})
		case errors.Is(err, aggregate_cooperation.ErrSelfCooperation):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		default:
			return nil, rpc_error.NewInternalError()
		}
	}

	return connect.NewResponse(&cooperation_v1.AcceptCooperationInvitationResponse{
		Cooperation: rpc_service_cooperation_serde.AggregateToPBCooperation(response.GetCooperation()),
	}), nil
}

func NewRPC(service service_accept_cooperation_invitation.Service) RPC {
	return RPC{
		service: service,
	}
}
