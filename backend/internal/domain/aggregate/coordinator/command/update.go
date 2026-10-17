package aggregate_coordinator_command

import (
	"time"

	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
)

type Update struct {
	description     string
	avatarPath      model_image_path.ImagePath
	cooldownSeconds int64
	now             time.Time
}

func (x Update) GetAvatarPath() model_image_path.ImagePath {
	return x.avatarPath
}

func (x Update) GetCooldownSeconds() int64 {
	return x.cooldownSeconds
}

func (x Update) GetDescription() string {
	return x.description
}

func (x Update) GetNow() time.Time {
	return x.now
}

func NewUpdate(description string, avatarPath model_image_path.ImagePath, cooldownSeconds int64, now time.Time) Update {
	return Update{
		description:     description,
		avatarPath:      avatarPath,
		cooldownSeconds: cooldownSeconds,
		now:             now,
	}
}
