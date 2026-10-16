package service_upload_image

import (
	"fmt"

	model_image_data "github.com/Roum1212/hahaton/internal/domain/model/image-data"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Request struct {
	userID model_user_id.UserID
	data   model_image_data.ImageData
}

func (x Request) GetData() model_image_data.ImageData {
	return x.data
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	return x.data.Validate()
}

func NewRequest(userID model_user_id.UserID, data model_image_data.ImageData) Request {
	return Request{
		userID: userID,
		data:   data,
	}
}

type Response struct {
	path model_image_path.ImagePath
}

func (x Response) GetPath() model_image_path.ImagePath {
	return x.path
}

func NewResponse(path model_image_path.ImagePath) Response {
	return Response{path: path}
}
