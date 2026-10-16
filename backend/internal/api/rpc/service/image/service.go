package rpc_service_image

import (
	"context"

	"connectrpc.com/connect"

	image_v1 "github.com/Roum1212/hahaton/genproto/image/v1"
	"github.com/Roum1212/hahaton/genproto/image/v1/image_v1connect"
	rpc_service_image_rpc_upload_image "github.com/Roum1212/hahaton/internal/api/rpc/service/image/rpc/upload-image"
)

var _ image_v1connect.ImageServiceHandler = Service{}

type Service struct {
	uploadImageRPC rpc_service_image_rpc_upload_image.RPC
}

func (x Service) UploadImage(
	ctx context.Context,
	r *connect.Request[image_v1.UploadImageRequest],
) (*connect.Response[image_v1.UploadImageResponse], error) {
	return x.uploadImageRPC.UploadImage(ctx, r)
}

func NewService(uploadImageRPC rpc_service_image_rpc_upload_image.RPC) Service {
	return Service{
		uploadImageRPC: uploadImageRPC,
	}
}
