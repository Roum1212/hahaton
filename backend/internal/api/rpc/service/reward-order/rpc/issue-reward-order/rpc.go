package rpc_service_reward_order_rpc_issue_reward_order

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"

	reward_order_v1 "github.com/Roum1212/hahaton/genproto/reward-order/v1"
	service_issue_reward_order "github.com/Roum1212/hahaton/internal/app/service/issue-reward-order"
	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_reward "github.com/Roum1212/hahaton/internal/domain/aggregate/reward"
	aggregate_reward_order "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order"
	model_reward_order_id "github.com/Roum1212/hahaton/internal/domain/model/reward-order-id"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_issue_reward_order.Service
}

func (x RPC) IssueRewardOrder(
	ctx context.Context,
	r *connect.Request[reward_order_v1.IssueRewardOrderRequest],
) (*connect.Response[emptypb.Empty], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	rewardOrderID, err := model_reward_order_id.NewRewardOrderID(r.Msg.GetId())
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

	if err := x.service.IssueRewardOrder(
		ctx,
		service_issue_reward_order.NewRequest(session.GetUserID(), rewardOrderID),
	); err != nil {
		switch {
		case errors.Is(err, service_issue_reward_order.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_issue_reward_order.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		case errors.Is(err, service_issue_reward_order.ErrFailedPrecondition):
			return nil, rpc_error.NewFailedPreconditionError(nil)
		case errors.Is(err, aggregate_reward_order.ErrInvalidState):
			return nil, rpc_error.NewFailedPreconditionError(nil)
		case errors.Is(err, aggregate_reward_order.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "reward_order.v1.RewardOrder",
				ResourceName: "id:" + rewardOrderID.String(),
				Description:  "reward order not found",
			})
		case errors.Is(err, aggregate_coordinator.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "coordinator.v1.Coordinator",
				ResourceName: "id:" + session.GetUserID().String(),
				Description:  "coordinator not found",
			})
		case errors.Is(err, aggregate_reward.ErrInvalidState):
			return nil, rpc_error.NewFailedPreconditionError(nil)
		case errors.Is(err, aggregate_reward.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "reward.v1.Reward",
				ResourceName: "",
				Description:  "reward not found",
			})
		case errors.Is(err, aggregate_cooperation.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "cooperation.v1.Cooperation",
				ResourceName: "",
				Description:  "cooperation not found",
			})
		default:
			return nil, rpc_error.NewInternalError()
		}
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func NewRPC(service service_issue_reward_order.Service) RPC {
	return RPC{service: service}
}
