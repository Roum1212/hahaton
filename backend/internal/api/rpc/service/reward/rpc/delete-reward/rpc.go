package rpc_service_reward_rpc_delete_reward

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"

	reward_v1 "github.com/Roum1212/hahaton/genproto/reward/v1"
	service_delete_reward "github.com/Roum1212/hahaton/internal/app/service/delete-reward"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_reward "github.com/Roum1212/hahaton/internal/domain/aggregate/reward"
	aggregate_reward_order "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order"
	aggregate_wallet "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet"
	model_reward_id "github.com/Roum1212/hahaton/internal/domain/model/reward-id"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_delete_reward.Service
}

func (x RPC) DeleteReward(
	ctx context.Context,
	r *connect.Request[reward_v1.DeleteRewardRequest],
) (*connect.Response[emptypb.Empty], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	rewardID, err := model_reward_id.NewRewardID(r.Msg.GetId())
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

	if err := x.service.DeleteReward(ctx, service_delete_reward.NewRequest(session.GetUserID(), rewardID)); err != nil {
		switch {
		case errors.Is(err, service_delete_reward.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_delete_reward.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		case errors.Is(err, aggregate_reward.ErrInvalidState):
			return nil, rpc_error.NewFailedPreconditionError(nil)
		case errors.Is(err, aggregate_reward.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "reward.v1.Reward",
				ResourceName: "id:" + rewardID.String(),
				Description:  "reward not found",
			})
		case errors.Is(err, aggregate_reward.ErrReservedQuantity):
			return nil, rpc_error.NewFailedPreconditionError(nil)
		case errors.Is(err, aggregate_coordinator.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "coordinator.v1.Coordinator",
				ResourceName: "id:" + session.GetUserID().String(),
				Description:  "coordinator not found",
			})
		case errors.Is(err, aggregate_reward_order.ErrInvalidState):
			return nil, rpc_error.NewFailedPreconditionError(nil)
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

func NewRPC(service service_delete_reward.Service) RPC {
	return RPC{service: service}
}
