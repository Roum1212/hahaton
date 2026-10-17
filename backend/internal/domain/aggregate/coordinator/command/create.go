package aggregate_coordinator_command

import (
	"time"

	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
)

type Create struct {
	description     string
	avatarPath      model_image_path.ImagePath
	cooldownSeconds int64
	now             time.Time
}

func (x Create) GetAvatarPath() model_image_path.ImagePath {
	return x.avatarPath
}

func (x Create) GetCooldownSeconds() int64 {
	return x.cooldownSeconds
}

func (x Create) GetDescription() string {
	return x.description
}

func (x Create) GetNow() time.Time {
	return x.now
}

func NewCreate(description string, avatarPath model_image_path.ImagePath, cooldownSeconds int64, now time.Time) Create {
	return Create{
		description:     description,
		avatarPath:      avatarPath,
		cooldownSeconds: cooldownSeconds,
		now:             now,
	}
}
