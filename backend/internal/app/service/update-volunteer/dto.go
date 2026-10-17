package service_update_volunteer

import (
	"errors"
	"fmt"

	model_date_of_birth "github.com/Roum1212/hahaton/internal/domain/model/date-of-birth"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

type Request struct {
	userID model_user_id.UserID

	avatarPath  utils_value.Optional[model_image_path.ImagePath]
	dateOfBirth utils_value.Optional[model_date_of_birth.DateOfBirth]
	firstName   utils_value.Optional[string]
	lastName    utils_value.Optional[string]
	middleName  utils_value.Optional[string]
}

func (x Request) GetAvatarPath() utils_value.Optional[model_image_path.ImagePath] {
	return x.avatarPath
}

func (x Request) GetDateOfBirth() utils_value.Optional[model_date_of_birth.DateOfBirth] {
	return x.dateOfBirth
}

func (x Request) GetFirstName() utils_value.Optional[string] {
	return x.firstName
}

func (x Request) GetLastName() utils_value.Optional[string] {
	return x.lastName
}

func (x Request) GetMiddleName() utils_value.Optional[string] {
	return x.middleName
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if value, ok := x.dateOfBirth.Get(); ok {
		if value.IsZero() {
			return errors.New("date of birth is required")
		}

		if err := value.Validate(); err != nil {
			return err
		}
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

func WithDateOfBirth(value model_date_of_birth.DateOfBirth) RequestOption {
	return func(x *Request) {
		x.dateOfBirth = utils_value.NewOptional(value)
	}
}

func WithFirstName(value string) RequestOption {
	return func(x *Request) {
		x.firstName = utils_value.NewOptional(value)
	}
}

func WithLastName(value string) RequestOption {
	return func(x *Request) {
		x.lastName = utils_value.NewOptional(value)
	}
}

func WithMiddleName(value string) RequestOption {
	return func(x *Request) {
		x.middleName = utils_value.NewOptional(value)
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
