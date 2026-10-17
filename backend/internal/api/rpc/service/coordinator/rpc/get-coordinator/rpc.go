package rpc_service_coordinator_rpc_get_coordinator

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	rpc_error "go.flida.dev/rpc-error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	coordinator_v1 "github.com/Roum1212/hahaton/genproto/coordinator/v1"
	rpc_service_coordinator_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/coordinator/serde"
	service_get_coordinator "github.com/Roum1212/hahaton/internal/app/service/get-coordinator"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_get_coordinator.Service
}

func (x RPC) GetCoordinator(
	ctx context.Context,
	r *connect.Request[coordinator_v1.GetCoordinatorRequest],
) (*connect.Response[coordinator_v1.GetCoordinatorResponse], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	var opts []service_get_coordinator.RequestOption

	coordinatorID := session.GetUserID()

	if r.Msg.Id != nil {
		value, err := model_user_id.NewUserID(*r.Msg.Id)
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

		opts = append(opts, service_get_coordinator.WithID(value))
		coordinatorID = value
	}

	response, err := x.service.GetCoordinator(ctx, service_get_coordinator.NewRequest(session.GetUserID(), opts...))
	if err != nil {
		switch {
		case errors.Is(err, service_get_coordinator.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, aggregate_coordinator.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "coordinator.v1.Coordinator",
				ResourceName: "id:" + coordinatorID.String(),
				Description:  "coordinator not found",
			})
		case errors.Is(err, aggregate_volunteer.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "volunteer.v1.Volunteer",
				ResourceName: "id:" + coordinatorID.String(),
				Description:  "volunteer not found",
			})
		default:
			return nil, rpc_error.NewInternalError()
		}
	}

	return connect.NewResponse(&coordinator_v1.GetCoordinatorResponse{
		Coordinator: rpc_service_coordinator_serde.AggregateToPBCoordinator(
			response.GetCoordinator(),
			response.GetVolunteer(),
		),
	}), nil
}

func NewRPC(service service_get_coordinator.Service) RPC {
	return RPC{service: service}
}
