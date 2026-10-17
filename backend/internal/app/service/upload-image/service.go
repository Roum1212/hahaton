package service_upload_image

import (
	"context"
	"errors"
	"fmt"

	model_image_data "github.com/Roum1212/hahaton/internal/domain/model/image-data"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	types_uuid_v7 "github.com/Roum1212/hahaton/internal/pkg/types/uuid-v7"
)

var ErrInvalidRequest = errors.New("invalid request")

type Storage interface {
	Save(ctx context.Context, path model_image_path.ImagePath, data model_image_data.ImageData) error
}

type Service interface {
	UploadImage(ctx context.Context, r Request) (Response, error)
}

type service struct {
	storage Storage
}

func (x service) UploadImage(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	id, err := types_uuid_v7.GenerateUUIDv7()
	if err != nil {
		return Response{}, err
	}

	path, err := model_image_path.NewImagePathFromFields(r.GetUserID(), id, r.GetData().GetExtension())
	if err != nil {
		return Response{}, err
	}

	if err := x.storage.Save(ctx, path, r.GetData()); err != nil {
		return Response{}, err
	}

	return NewResponse(path), nil
}

func NewService(storage Storage) Service {
	return service{storage: storage}
}
