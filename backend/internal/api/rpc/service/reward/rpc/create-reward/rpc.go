package rpc_service_reward_rpc_create_reward

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	reward_v1 "github.com/Roum1212/hahaton/genproto/reward/v1"
	service_create_reward "github.com/Roum1212/hahaton/internal/app/service/create-reward"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_reward "github.com/Roum1212/hahaton/internal/domain/aggregate/reward"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_create_reward.Service
}

func (x RPC) CreateReward(
	ctx context.Context,
	r *connect.Request[reward_v1.CreateRewardRequest],
) (*connect.Response[reward_v1.CreateRewardResponse], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	title := r.Msg.GetTitle()
	description := r.Msg.GetDescription()

	imagePath, err := model_image_path.NewImagePath(r.Msg.GetImagePath())
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

	priceCoins := r.Msg.GetPriceCoins()
	quantityLimit := r.Msg.GetQuantityLimit()

	response, err := x.service.CreateReward(ctx, service_create_reward.NewRequest(
		session.GetUserID(),
		title,
		description,
		imagePath,
		priceCoins,
		quantityLimit,
	))
	if err != nil {
		switch {
		case errors.Is(err, service_create_reward.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_create_reward.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		case errors.Is(err, service_create_reward.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "storage.Image",
				ResourceName: "path:" + r.Msg.GetImagePath(),
				Description:  "image not found",
			})
		case errors.Is(err, aggregate_reward.ErrAlreadyExists):
			return nil, rpc_error.NewAlreadyExistsError(nil)
		case errors.Is(err, aggregate_reward.ErrInvalidState):
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

	return connect.NewResponse(&reward_v1.CreateRewardResponse{Id: response.GetRewardID().String()}), nil
}

func NewRPC(service service_create_reward.Service) RPC {
	return RPC{service: service}
}
