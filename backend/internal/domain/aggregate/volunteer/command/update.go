package aggregate_volunteer_command

import (
	"time"

	model_date_of_birth "github.com/Roum1212/hahaton/internal/domain/model/date-of-birth"
	model_full_name "github.com/Roum1212/hahaton/internal/domain/model/full-name"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
)

type Update struct {
	fullName    model_full_name.FullName
	dateOfBirth model_date_of_birth.DateOfBirth
	avatarPath  model_image_path.ImagePath
	now         time.Time
}

func (x Update) GetAvatarPath() model_image_path.ImagePath {
	return x.avatarPath
}

func (x Update) GetDateOfBirth() model_date_of_birth.DateOfBirth {
	return x.dateOfBirth
}

func (x Update) GetFullName() model_full_name.FullName {
	return x.fullName
}

func (x Update) GetNow() time.Time {
	return x.now
}

func NewUpdate(
	fullName model_full_name.FullName,
	dateOfBirth model_date_of_birth.DateOfBirth,
	avatarPath model_image_path.ImagePath,
	now time.Time,
) Update {
	return Update{
		fullName:    fullName,
		dateOfBirth: dateOfBirth,
		avatarPath:  avatarPath,
		now:         now,
	}
}
