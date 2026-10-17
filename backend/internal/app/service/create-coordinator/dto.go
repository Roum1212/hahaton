package service_create_coordinator

import (
	"errors"
	"fmt"

	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Request struct {
	volunteerID     model_user_id.UserID
	description     string
	avatarPath      model_image_path.ImagePath
	cooldownSeconds int64
}

func (x Request) GetAvatarPath() model_image_path.ImagePath {
	return x.avatarPath
}

func (x Request) GetCooldownSeconds() int64 {
	return x.cooldownSeconds
}

func (x Request) GetDescription() string {
	return x.description
}

func (x Request) GetVolunteerID() model_user_id.UserID {
	return x.volunteerID
}

func (x Request) Validate() error {
	if err := x.volunteerID.Validate(); err != nil {
		return fmt.Errorf("volunteer id: %w", err)
	}

	if x.cooldownSeconds < 0 {
		return errors.New("cooldown must not be negative")
	}

	if err := x.avatarPath.Validate(); err != nil {
		return fmt.Errorf("avatar path: %w", err)
	}

	return nil
}

func NewRequest(
	volunteerID model_user_id.UserID,
	description string,
	avatarPath model_image_path.ImagePath,
	cooldownSeconds int64,
) Request {
	return Request{
		volunteerID:     volunteerID,
		description:     description,
		avatarPath:      avatarPath,
		cooldownSeconds: cooldownSeconds,
	}
}

type Response struct {
	coordinatorID model_user_id.UserID
}

func (x Response) GetCoordinatorID() model_user_id.UserID {
	return x.coordinatorID
}

func NewResponse(coordinatorID model_user_id.UserID) Response {
	return Response{
		coordinatorID: coordinatorID,
	}
}
