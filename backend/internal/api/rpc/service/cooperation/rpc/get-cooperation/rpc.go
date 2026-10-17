package rpc_service_cooperation_rpc_get_cooperation

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	rpc_error "go.flida.dev/rpc-error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	cooperation_v1 "github.com/Roum1212/hahaton/genproto/cooperation/v1"
	rpc_service_cooperation_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/cooperation/serde"
	service_get_cooperation "github.com/Roum1212/hahaton/internal/app/service/get-cooperation"
	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_get_cooperation.Service
}

func (x RPC) GetCooperation(
	ctx context.Context,
	r *connect.Request[cooperation_v1.GetCooperationRequest],
) (*connect.Response[cooperation_v1.GetCooperationResponse], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	counterpartRaw := ""
	counterpartField := "counterpart"
	asCoordinator := false

	switch value := r.Msg.GetCounterpart().(type) {
	case *cooperation_v1.GetCooperationRequest_CoordinatorId:
		counterpartRaw = value.CoordinatorId
		counterpartField = "coordinator_id"
	case *cooperation_v1.GetCooperationRequest_VolunteerId:
		counterpartRaw = value.VolunteerId
		counterpartField = "volunteer_id"
		asCoordinator = true
	}

	counterpartID, err := model_user_id.NewUserID(counterpartRaw)
	if err != nil {
		return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       counterpartField,
					Description: err.Error(),
				},
			},
		})
	}

	coordinatorID, volunteerID := counterpartID, session.GetUserID()
	if asCoordinator {
		coordinatorID, volunteerID = session.GetUserID(), counterpartID
	}

	response, err := x.service.GetCooperation(ctx, service_get_cooperation.NewRequest(
		session.GetUserID(),
		counterpartID,
		asCoordinator,
	))
	if err != nil {
		switch {
		case errors.Is(err, service_get_cooperation.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_get_cooperation.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		case errors.Is(err, aggregate_cooperation.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "cooperation.v1.Cooperation",
				ResourceName: "coordinator_id:" + coordinatorID.String() + ",volunteer_id:" + volunteerID.String(),
				Description:  "cooperation not found",
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

	return connect.NewResponse(&cooperation_v1.GetCooperationResponse{
		Cooperation: rpc_service_cooperation_serde.AggregateToPBCooperation(response.GetCooperation()),
	}), nil
}

func NewRPC(service service_get_cooperation.Service) RPC {
	return RPC{service: service}
}
