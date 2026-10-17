package rpc_service_reward_rpc_get_reward

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	rpc_error "go.flida.dev/rpc-error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	reward_v1 "github.com/Roum1212/hahaton/genproto/reward/v1"
	rpc_service_reward_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/reward/serde"
	service_get_reward "github.com/Roum1212/hahaton/internal/app/service/get-reward"
	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_reward "github.com/Roum1212/hahaton/internal/domain/aggregate/reward"
	model_reward_id "github.com/Roum1212/hahaton/internal/domain/model/reward-id"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_get_reward.Service
}

func (x RPC) GetReward(
	ctx context.Context,
	r *connect.Request[reward_v1.GetRewardRequest],
) (*connect.Response[reward_v1.GetRewardResponse], error) {
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

	response, err := x.service.GetReward(ctx, service_get_reward.NewRequest(session.GetUserID(), rewardID))
	if err != nil {
		switch {
		case errors.Is(err, service_get_reward.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_get_reward.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		case errors.Is(err, aggregate_reward.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "reward.v1.Reward",
				ResourceName: "id:" + rewardID.String(),
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

	return connect.NewResponse(&reward_v1.GetRewardResponse{
		Reward: rpc_service_reward_serde.AggregateToPBReward(response.GetReward()),
	}), nil
}

func NewRPC(service service_get_reward.Service) RPC {
	return RPC{service: service}
}
