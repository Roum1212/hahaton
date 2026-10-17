package rpc_service_coordinator_rpc_get_coordinator_settings

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	rpc_error "go.flida.dev/rpc-error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	coordinator_v1 "github.com/Roum1212/hahaton/genproto/coordinator/v1"
	service_get_coordinator_settings "github.com/Roum1212/hahaton/internal/app/service/get-coordinator-settings"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_get_coordinator_settings.Service
}

func (x RPC) GetCoordinatorSettings(
	ctx context.Context,
	r *connect.Request[coordinator_v1.GetCoordinatorSettingsRequest],
) (*connect.Response[coordinator_v1.GetCoordinatorSettingsResponse], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	response, err := x.service.GetCoordinatorSettings(
		ctx,
		service_get_coordinator_settings.NewRequest(session.GetUserID()),
	)
	if err != nil {
		switch {
		case errors.Is(err, service_get_coordinator_settings.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, aggregate_coordinator.ErrNotFound):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		default:
			return nil, rpc_error.NewInternalError()
		}
	}

	return connect.NewResponse(&coordinator_v1.GetCoordinatorSettingsResponse{
		Settings: &coordinator_v1.CoordinatorSettings{
			CooperationRequestCooldownSeconds: response.GetCooldownSeconds(),
		},
	}), nil
}

func NewRPC(service service_get_coordinator_settings.Service) RPC {
	return RPC{service: service}
}
