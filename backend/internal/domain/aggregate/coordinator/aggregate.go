package aggregate_coordinator

import (
	"errors"
	"time"

	aggregate_coordinator_command "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator/command"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

var ErrAlreadyExists = errors.New("coordinator already exists")

type Coordinator struct {
	userID          model_user_id.UserID
	description     string
	avatarPath      model_image_path.ImagePath
	cooldownSeconds int64
	createdAt       time.Time
	updatedAt       time.Time
}

func (x *Coordinator) GetAvatarPath() model_image_path.ImagePath {
	return x.avatarPath
}

func (x *Coordinator) GetCooldownSeconds() int64 {
	return x.cooldownSeconds
}

func (x *Coordinator) GetCreatedAt() time.Time {
	return x.createdAt
}

func (x *Coordinator) GetDescription() string {
	return x.description
}

func (x *Coordinator) GetUpdatedAt() time.Time {
	return x.updatedAt
}

func (x *Coordinator) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x *Coordinator) Create(command aggregate_coordinator_command.Create) error {
	if !x.createdAt.IsZero() {
		return ErrAlreadyExists
	}

	x.description = command.GetDescription()
	x.avatarPath = command.GetAvatarPath()
	x.cooldownSeconds = command.GetCooldownSeconds()
	x.createdAt = command.GetNow()
	x.updatedAt = command.GetNow()

	return nil
}

func (x *Coordinator) Update(command aggregate_coordinator_command.Update) {
	x.description = command.GetDescription()
	x.avatarPath = command.GetAvatarPath()
	x.cooldownSeconds = command.GetCooldownSeconds()
	x.updatedAt = command.GetNow()
}

func NewCoordinator(userID model_user_id.UserID) *Coordinator {
	return &Coordinator{
		userID: userID,
	}
}

func NewCoordinatorFromFields(
	userID model_user_id.UserID,
	description string,
	avatarPath model_image_path.ImagePath,
	cooldownSeconds int64,
	createdAt time.Time,
	updatedAt time.Time,
) *Coordinator {
	return &Coordinator{
		userID:          userID,
		description:     description,
		avatarPath:      avatarPath,
		cooldownSeconds: cooldownSeconds,
		createdAt:       createdAt,
		updatedAt:       updatedAt,
	}
}
