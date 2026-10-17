package rpc_service_cooperation_rpc_request_cooperation

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	cooperation_v1 "github.com/Roum1212/hahaton/genproto/cooperation/v1"
	rpc_service_cooperation_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/cooperation/serde"
	service_request_cooperation "github.com/Roum1212/hahaton/internal/app/service/request-cooperation"
	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_request_cooperation.Service
}

func (x RPC) RequestCooperation(
	ctx context.Context,
	r *connect.Request[cooperation_v1.RequestCooperationRequest],
) (*connect.Response[cooperation_v1.RequestCooperationResponse], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	coordinatorID, err := model_user_id.NewUserID(r.Msg.GetCoordinatorId())
	if err != nil {
		return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "coordinator_id",
					Description: err.Error(),
				},
			},
		})
	}

	response, err := x.service.RequestCooperation(ctx, service_request_cooperation.NewRequest(
		session.GetUserID(),
		coordinatorID,
	))
	if err != nil {
		switch {
		case errors.Is(err, service_request_cooperation.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, aggregate_cooperation.ErrAlreadyExists):
			return nil, rpc_error.NewAlreadyExistsError(nil)
		case errors.Is(err, aggregate_cooperation.ErrCooldown):
			return nil, rpc_error.NewFailedPreconditionError(nil)
		case errors.Is(err, aggregate_cooperation.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "cooperation.v1.Cooperation",
				ResourceName: "coordinator_id:" + coordinatorID.String() + ",volunteer_id:" + session.GetUserID().String(),
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
		case errors.Is(err, aggregate_coordinator.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "coordinator.v1.Coordinator",
				ResourceName: "id:" + coordinatorID.String(),
				Description:  "coordinator not found",
			})
		default:
			return nil, rpc_error.NewInternalError()
		}
	}

	return connect.NewResponse(&cooperation_v1.RequestCooperationResponse{
		Cooperation: rpc_service_cooperation_serde.AggregateToPBCooperation(response.GetCooperation()),
	}), nil
}

func NewRPC(service service_request_cooperation.Service) RPC {
	return RPC{service: service}
}
