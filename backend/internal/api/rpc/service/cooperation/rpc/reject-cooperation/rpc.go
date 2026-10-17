package rpc_service_cooperation_rpc_reject_cooperation

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"

	cooperation_v1 "github.com/Roum1212/hahaton/genproto/cooperation/v1"
	service_reject_cooperation "github.com/Roum1212/hahaton/internal/app/service/reject-cooperation"
	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_reject_cooperation.Service
}

func (x RPC) RejectCooperation(
	ctx context.Context,
	r *connect.Request[cooperation_v1.RejectCooperationRequest],
) (*connect.Response[emptypb.Empty], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	volunteerID, err := model_user_id.NewUserID(r.Msg.GetVolunteerId())
	if err != nil {
		return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "volunteer_id",
					Description: err.Error(),
				},
			},
		})
	}

	if err = x.service.RejectCooperation(
		ctx,
		service_reject_cooperation.NewRequest(session.GetUserID(), volunteerID),
	); err != nil {
		switch {
		case errors.Is(err, service_reject_cooperation.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_reject_cooperation.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		case errors.Is(err, aggregate_cooperation.ErrInvalidState):
			return nil, rpc_error.NewFailedPreconditionError(nil)
		case errors.Is(err, aggregate_cooperation.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "cooperation.v1.Cooperation",
				ResourceName: "coordinator_id:" + session.GetUserID().String() + ",volunteer_id:" + volunteerID.String(),
				Description:  "cooperation not found",
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

func NewRPC(service service_reject_cooperation.Service) RPC {
	return RPC{service: service}
}
