package rpc_service_coordinator_rpc_update_coordinator

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"

	coordinator_v1 "github.com/Roum1212/hahaton/genproto/coordinator/v1"
	service_update_coordinator "github.com/Roum1212/hahaton/internal/app/service/update-coordinator"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_update_coordinator.Service
}

func (x RPC) UpdateCoordinator(
	ctx context.Context,
	r *connect.Request[coordinator_v1.UpdateCoordinatorRequest],
) (*connect.Response[emptypb.Empty], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	var opts []service_update_coordinator.RequestOption

	if value := r.Msg.AvatarPath; value != nil {
		path, err := model_image_path.NewImagePath(*value)
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

		opts = append(opts, service_update_coordinator.WithAvatarPath(path))
	}

	if value := r.Msg.Description; value != nil {
		opts = append(opts, service_update_coordinator.WithDescription(*value))
	}

	if err := x.service.UpdateCoordinator(
		ctx,
		service_update_coordinator.NewRequest(session.GetUserID(), opts...),
	); err != nil {
		switch {
		case errors.Is(err, service_update_coordinator.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_update_coordinator.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		case errors.Is(err, service_update_coordinator.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "storage.Image",
				ResourceName: "path:" + r.Msg.GetAvatarPath(),
				Description:  "image not found",
			})
		case errors.Is(err, aggregate_coordinator.ErrNotFound):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		default:
			return nil, rpc_error.NewInternalError()
		}
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func NewRPC(service service_update_coordinator.Service) RPC {
	return RPC{
		service: service,
	}
}
