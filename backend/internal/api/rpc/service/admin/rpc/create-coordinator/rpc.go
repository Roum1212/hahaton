package rpc_service_admin_rpc_create_coordinator

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	rpc_error "go.flida.dev/rpc-error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	admin_v1 "github.com/Roum1212/hahaton/genproto/admin/v1"
	service_create_coordinator "github.com/Roum1212/hahaton/internal/app/service/create-coordinator"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type RPC struct {
	service service_create_coordinator.Service
}

func (x RPC) CreateCoordinator(
	ctx context.Context,
	r *connect.Request[admin_v1.CreateCoordinatorRequest],
) (*connect.Response[admin_v1.CreateCoordinatorResponse], error) {
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

	avatarPath, err := model_image_path.NewImagePath(r.Msg.GetAvatarPath())
	if err != nil {
		return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "avatar_path",
					Description: err.Error(),
				},
			},
		})
	}

	if r.Msg.CooperationRequestCooldownSeconds == nil {
		return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "cooperation_request_cooldown_seconds",
					Description: "cooldown is required",
				},
			},
		})
	}

	response, err := x.service.CreateCoordinator(ctx, service_create_coordinator.NewRequest(
		volunteerID,
		r.Msg.GetDescription(),
		avatarPath,
		r.Msg.GetCooperationRequestCooldownSeconds(),
	))
	if err != nil {
		switch {
		case errors.Is(err, service_create_coordinator.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_create_coordinator.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		case errors.Is(err, service_create_coordinator.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "storage.Image",
				ResourceName: "path:" + r.Msg.GetAvatarPath(),
				Description:  "image not found",
			})
		case errors.Is(err, aggregate_volunteer.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "volunteer.v1.Volunteer",
				ResourceName: "id:" + volunteerID.String(),
				Description:  "volunteer not found",
			})
		case errors.Is(err, aggregate_coordinator.ErrAlreadyExists):
			return nil, rpc_error.NewAlreadyExistsError(nil)
		default:
			return nil, rpc_error.NewInternalError()
		}
	}

	return connect.NewResponse(&admin_v1.CreateCoordinatorResponse{Id: response.GetCoordinatorID().String()}), nil
}

func NewRPC(service service_create_coordinator.Service) RPC {
	return RPC{service: service}
}
