package rpc_service_cooperation_rpc_create_cooperation_invitation

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	rpc_error "go.flida.dev/rpc-error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	cooperation_v1 "github.com/Roum1212/hahaton/genproto/cooperation/v1"
	rpc_serde "github.com/Roum1212/hahaton/internal/api/rpc/serde"
	rpc_service_cooperation_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/cooperation/serde"
	service_create_cooperation_invitation "github.com/Roum1212/hahaton/internal/app/service/create-cooperation-invitation"
	aggregate_cooperation_invitation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation-invitation"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service     service_create_cooperation_invitation.Service
	frontendURL string
}

func (x RPC) CreateCooperationInvitation(
	ctx context.Context,
	r *connect.Request[cooperation_v1.CreateCooperationInvitationRequest],
) (*connect.Response[cooperation_v1.CreateCooperationInvitationResponse], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	var opts []service_create_cooperation_invitation.RequestOption

	if r.Msg.ExpiresAt != nil {
		value, err := rpc_serde.TimestampFromPB(r.Msg.GetExpiresAt())
		if err != nil {
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "expires_at",
						Description: err.Error(),
					},
				},
			})
		}

		opts = append(opts, service_create_cooperation_invitation.WithExpiresAt(value))
	}

	response, err := x.service.CreateCooperationInvitation(
		ctx,
		service_create_cooperation_invitation.NewRequest(session.GetUserID(), opts...),
	)
	if err != nil {
		switch {
		case errors.Is(err, service_create_cooperation_invitation.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_create_cooperation_invitation.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		case errors.Is(err, aggregate_cooperation_invitation.ErrAlreadyExists):
			return nil, rpc_error.NewAlreadyExistsError(nil)
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

	return connect.NewResponse(&cooperation_v1.CreateCooperationInvitationResponse{
		Invitation: rpc_service_cooperation_serde.InvitationToPB(response.GetInvitation(), x.frontendURL),
	}), nil
}

func NewRPC(service service_create_cooperation_invitation.Service, frontendURL string) RPC {
	return RPC{
		service:     service,
		frontendURL: frontendURL,
	}
}
