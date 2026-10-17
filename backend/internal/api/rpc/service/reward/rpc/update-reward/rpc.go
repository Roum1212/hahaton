package rpc_service_reward_rpc_update_reward

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	rpc_error "go.flida.dev/rpc-error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"

	reward_v1 "github.com/Roum1212/hahaton/genproto/reward/v1"
	service_update_reward "github.com/Roum1212/hahaton/internal/app/service/update-reward"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_reward "github.com/Roum1212/hahaton/internal/domain/aggregate/reward"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	model_reward_id "github.com/Roum1212/hahaton/internal/domain/model/reward-id"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_update_reward.Service
}

func (x RPC) UpdateReward(
	ctx context.Context,
	r *connect.Request[reward_v1.UpdateRewardRequest],
) (*connect.Response[emptypb.Empty], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	var opts []service_update_reward.RequestOption

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

	if value := r.Msg.Description; value != nil {
		opts = append(opts, service_update_reward.WithDescription(*value))
	}

	if value := r.Msg.ImagePath; value != nil {
		path, err := model_image_path.NewImagePath(*value)
		if err != nil {
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "image_path",
						Description: err.Error(),
					},
				},
			})
		}

		opts = append(opts, service_update_reward.WithImagePath(path))
	}

	if value := r.Msg.Title; value != nil {
		opts = append(opts, service_update_reward.WithTitle(*value))
	}

	if value := r.Msg.PriceCoins; value != nil {
		opts = append(opts, service_update_reward.WithPriceCoins(*value))
	}

	if value := r.Msg.QuantityLimit; value != nil {
		opts = append(opts, service_update_reward.WithQuantityLimit(*value))
	}

	if err := x.service.UpdateReward(
		ctx,
		service_update_reward.NewRequest(session.GetUserID(), rewardID, opts...),
	); err != nil {
		switch {
		case errors.Is(err, service_update_reward.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_update_reward.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		case errors.Is(err, service_update_reward.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "storage.Image",
				ResourceName: "path:" + r.Msg.GetImagePath(),
				Description:  "image not found",
			})
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
		default:
			return nil, rpc_error.NewInternalError()
		}
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func NewRPC(service service_update_reward.Service) RPC {
	return RPC{service: service}
}
