package rpc_service_coordinator_rpc_update_coordinator_settings

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	rpc_error "go.flida.dev/rpc-error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"

	coordinator_v1 "github.com/Roum1212/hahaton/genproto/coordinator/v1"
	service_update_coordinator_settings "github.com/Roum1212/hahaton/internal/app/service/update-coordinator-settings"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_update_coordinator_settings.Service
}

func (x RPC) UpdateCoordinatorSettings(
	ctx context.Context,
	r *connect.Request[coordinator_v1.UpdateCoordinatorSettingsRequest],
) (*connect.Response[emptypb.Empty], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	var opts []service_update_coordinator_settings.RequestOption

	if value := r.Msg.CooperationRequestCooldownSeconds; value != nil {
		opts = append(opts, service_update_coordinator_settings.WithCooldownSeconds(*value))
	}

	if err := x.service.UpdateCoordinatorSettings(
		ctx,
		service_update_coordinator_settings.NewRequest(session.GetUserID(), opts...),
	); err != nil {
		switch {
		case errors.Is(err, service_update_coordinator_settings.ErrInvalidRequest):
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

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func NewRPC(service service_update_coordinator_settings.Service) RPC {
	return RPC{service: service}
}
