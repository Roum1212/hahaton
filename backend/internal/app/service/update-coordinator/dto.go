package service_update_coordinator

import (
	"fmt"

	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

type Request struct {
	userID model_user_id.UserID

	avatarPath  utils_value.Optional[model_image_path.ImagePath]
	description utils_value.Optional[string]
}

func (x Request) GetAvatarPath() utils_value.Optional[model_image_path.ImagePath] {
	return x.avatarPath
}

func (x Request) GetDescription() utils_value.Optional[string] {
	return x.description
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if value, ok := x.avatarPath.Get(); ok {
		if err := value.Validate(); err != nil {
			return fmt.Errorf("avatar path: %w", err)
		}
	}

	return nil
}

type RequestOption func(x *Request)

func WithAvatarPath(value model_image_path.ImagePath) RequestOption {
	return func(x *Request) {
		x.avatarPath = utils_value.NewOptional(value)
	}
}

func WithDescription(value string) RequestOption {
	return func(x *Request) {
		x.description = utils_value.NewOptional(value)
	}
}

func NewRequest(userID model_user_id.UserID, opts ...RequestOption) Request {
	x := Request{
		userID: userID,
	}
	for _, option := range opts {
		option(&x)
	}

	return x
}
