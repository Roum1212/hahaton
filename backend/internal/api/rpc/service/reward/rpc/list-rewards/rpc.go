package rpc_service_reward_rpc_list_rewards

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	reward_v1 "github.com/Roum1212/hahaton/genproto/reward/v1"
	rpc_service_reward_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/reward/serde"
	service_list_rewards "github.com/Roum1212/hahaton/internal/app/service/list-rewards"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_list_rewards.Service
}

func (x RPC) ListRewards(
	ctx context.Context,
	r *connect.Request[reward_v1.ListRewardsRequest],
) (*connect.Response[reward_v1.ListRewardsResponse], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	var opts []service_list_rewards.RequestOption

	scope := int32(r.Msg.GetScope())

	if r.Msg.CoordinatorId != nil {
		value, err := model_user_id.NewUserID(*r.Msg.CoordinatorId)
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

		opts = append(opts, service_list_rewards.WithCoordinatorID(value))
	}

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

	response, err := x.service.ListRewards(
		ctx,
		service_list_rewards.NewRequest(session.GetUserID(), scope, page, opts...),
	)
	if err != nil {
		switch {
		case errors.Is(err, service_list_rewards.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_list_rewards.ErrPermissionDenied):
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

	values := make([]*reward_v1.Reward, 0, len(response.GetRewards()))
	for _, value := range response.GetRewards() {
		values = append(values, rpc_service_reward_serde.AggregateToPBReward(value))
	}

	return connect.NewResponse(&reward_v1.ListRewardsResponse{
		Rewards:    values,
		TotalCount: response.GetTotalCount(),
	}), nil
}

func NewRPC(service service_list_rewards.Service) RPC {
	return RPC{service: service}
}
