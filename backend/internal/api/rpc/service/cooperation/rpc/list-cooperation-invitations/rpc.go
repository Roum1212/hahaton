package rpc_service_cooperation_rpc_list_cooperation_invitations

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	cooperation_v1 "github.com/Roum1212/hahaton/genproto/cooperation/v1"
	rpc_service_cooperation_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/cooperation/serde"
	service_list_cooperation_invitations "github.com/Roum1212/hahaton/internal/app/service/list-cooperation-invitations"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service     service_list_cooperation_invitations.Service
	frontendURL string
}

func (x RPC) ListCooperationInvitations(
	ctx context.Context,
	r *connect.Request[cooperation_v1.ListCooperationInvitationsRequest],
) (*connect.Response[cooperation_v1.ListCooperationInvitationsResponse], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	includeExpired := r.Msg.GetIncludeExpired()

	page, err := model_page.NewPage(r.Msg.GetLimit(), r.Msg.GetOffset())
	if err != nil {
		return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "pagination",
					Description: err.Error(),
				},
			},
		})
	}

	response, err := x.service.ListCooperationInvitations(ctx, service_list_cooperation_invitations.NewRequest(
		session.GetUserID(),
		includeExpired,
		page,
	))
	if err != nil {
		switch {
		case errors.Is(err, service_list_cooperation_invitations.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_list_cooperation_invitations.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
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

	values := make([]*cooperation_v1.CooperationInvitation, 0, len(response.GetInvitations()))
	for _, value := range response.GetInvitations() {
		values = append(values, rpc_service_cooperation_serde.InvitationToPB(value, x.frontendURL))
	}

	return connect.NewResponse(&cooperation_v1.ListCooperationInvitationsResponse{
		Invitations: values,
		TotalCount:  response.GetTotalCount(),
	}), nil
}

func NewRPC(service service_list_cooperation_invitations.Service, frontendURL string) RPC {
	return RPC{
		service:     service,
		frontendURL: frontendURL,
	}
}
