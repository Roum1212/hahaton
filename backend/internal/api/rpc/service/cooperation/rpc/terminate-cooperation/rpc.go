package rpc_service_cooperation_rpc_terminate_cooperation

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"

	cooperation_v1 "github.com/Roum1212/hahaton/genproto/cooperation/v1"
	service_terminate_cooperation "github.com/Roum1212/hahaton/internal/app/service/terminate-cooperation"
	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	aggregate_reward "github.com/Roum1212/hahaton/internal/domain/aggregate/reward"
	aggregate_reward_order "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order"
	aggregate_wallet "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_terminate_cooperation.Service
}

func (x RPC) TerminateCooperation(
	ctx context.Context,
	r *connect.Request[cooperation_v1.TerminateCooperationRequest],
) (*connect.Response[emptypb.Empty], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	counterpartRaw := ""
	counterpartField := "counterpart"
	asCoordinator := false

	switch value := r.Msg.GetCounterpart().(type) {
	case *cooperation_v1.TerminateCooperationRequest_CoordinatorId:
		counterpartRaw = value.CoordinatorId
		counterpartField = "coordinator_id"
	case *cooperation_v1.TerminateCooperationRequest_VolunteerId:
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

	if err := x.service.TerminateCooperation(ctx, service_terminate_cooperation.NewRequest(
		session.GetUserID(),
		counterpartID,
		asCoordinator,
	)); err != nil {
		switch {
		case errors.Is(err, service_terminate_cooperation.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_terminate_cooperation.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		case errors.Is(err, aggregate_cooperation.ErrInvalidState):
			return nil, rpc_error.NewFailedPreconditionError(nil)
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
		case errors.Is(err, aggregate_participation.ErrInvalidState):
			return nil, rpc_error.NewFailedPreconditionError(nil)
		case errors.Is(err, aggregate_reward_order.ErrInvalidState):
			return nil, rpc_error.NewFailedPreconditionError(nil)
		case errors.Is(err, aggregate_reward.ErrInvalidState):
			return nil, rpc_error.NewFailedPreconditionError(nil)
		case errors.Is(err, aggregate_reward.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "reward.v1.Reward",
				ResourceName: "",
				Description:  "reward not found",
			})
		case errors.Is(err, aggregate_wallet.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "wallet.Wallet",
				ResourceName: "",
				Description:  "wallet not found",
			})
		case errors.Is(err, aggregate_wallet.ErrOverflow):
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

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func NewRPC(service service_terminate_cooperation.Service) RPC {
	return RPC{service: service}
}
