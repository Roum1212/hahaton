package rpc_service_image_rpc_upload_image

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	rpc_error "go.flida.dev/rpc-error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	image_v1 "github.com/Roum1212/hahaton/genproto/image/v1"
	service_upload_image "github.com/Roum1212/hahaton/internal/app/service/upload-image"
	model_image_data "github.com/Roum1212/hahaton/internal/domain/model/image-data"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_upload_image.Service
}

func (x RPC) UploadImage(
	ctx context.Context,
	r *connect.Request[image_v1.UploadImageRequest],
) (*connect.Response[image_v1.UploadImageResponse], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok || session.GetUserID().IsZero() {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	data, err := model_image_data.NewImageData(r.Msg.GetData())
	if err != nil {
		return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "data",
					Description: err.Error(),
				},
			},
		})
	}

	response, err := x.service.UploadImage(ctx, service_upload_image.NewRequest(session.GetUserID(), data))
	if err != nil {
		switch {
		case errors.Is(err, service_upload_image.ErrInvalidRequest):
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

	return connect.NewResponse(&image_v1.UploadImageResponse{
		ImagePath: response.GetPath().String(),
		ImageUrl:  "/storage/" + response.GetPath().String(),
	}), nil
}

func NewRPC(service service_upload_image.Service) RPC {
	return RPC{
		service: service,
	}
}
