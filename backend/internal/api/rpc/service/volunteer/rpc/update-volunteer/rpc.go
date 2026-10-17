package rpc_service_volunteer_rpc_update_volunteer

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	rpc_error "go.flida.dev/rpc-error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"

	volunteer_v1 "github.com/Roum1212/hahaton/genproto/volunteer/v1"
	rpc_service_volunteer_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/volunteer/serde"
	service_update_volunteer "github.com/Roum1212/hahaton/internal/app/service/update-volunteer"
	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_update_volunteer.Service
}

func (x RPC) UpdateVolunteer(
	ctx context.Context,
	r *connect.Request[volunteer_v1.UpdateVolunteerRequest],
) (*connect.Response[emptypb.Empty], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	var opts []service_update_volunteer.RequestOption

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

		opts = append(opts, service_update_volunteer.WithAvatarPath(path))
	}

	if r.Msg.DateOfBirth != nil {
		value, err := rpc_service_volunteer_serde.DateFromPB(r.Msg.GetDateOfBirth())
		if err != nil {
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "date_of_birth",
						Description: err.Error(),
					},
				},
			})
		}

		opts = append(opts, service_update_volunteer.WithDateOfBirth(value))
	}

	if value := r.Msg.FirstName; value != nil {
		opts = append(opts, service_update_volunteer.WithFirstName(*value))
	}

	if value := r.Msg.LastName; value != nil {
		opts = append(opts, service_update_volunteer.WithLastName(*value))
	}

	if value := r.Msg.MiddleName; value != nil {
		opts = append(opts, service_update_volunteer.WithMiddleName(*value))
	}

	if err := x.service.UpdateVolunteer(
		ctx,
		service_update_volunteer.NewRequest(session.GetUserID(), opts...),
	); err != nil {
		switch {
		case errors.Is(err, service_update_volunteer.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_update_volunteer.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		case errors.Is(err, service_update_volunteer.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "storage.Image",
				ResourceName: "path:" + r.Msg.GetAvatarPath(),
				Description:  "image not found",
			})
		case errors.Is(err, aggregate_volunteer.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "volunteer.v1.Volunteer",
				ResourceName: "id:" + session.GetUserID().String(),
				Description:  "volunteer not found",
			})
		default:
			return nil, rpc_error.NewInternalError()
		}
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func NewRPC(service service_update_volunteer.Service) RPC {
	return RPC{service: service}
}
