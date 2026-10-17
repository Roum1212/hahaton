package repository_coordinator_postgresql

import (
	"fmt"
	"time"

	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type row struct {
	userID          string
	description     string
	avatarPath      string
	cooldownSeconds int64
	createdAt       time.Time
	updatedAt       time.Time
}

func (x *row) scan(source interface{ Scan(dest ...any) error }) error {
	return source.Scan(&x.userID, &x.description, &x.avatarPath, &x.cooldownSeconds, &x.createdAt, &x.updatedAt)
}

func (x *row) toAggregate() (*aggregate_coordinator.Coordinator, error) {
	userID, err := model_user_id.NewUserID(x.userID)
	if err != nil {
		return nil, fmt.Errorf("restore user id: %w", err)
	}

	var avatarPath model_image_path.ImagePath
	if len(x.avatarPath) != 0 {
		value, err := model_image_path.NewImagePath(x.avatarPath)
		if err != nil {
			return nil, fmt.Errorf("restore value: %w", err)
		}

		avatarPath = value
	}

	return aggregate_coordinator.NewCoordinatorFromFields(
		userID,
		x.description,
		avatarPath,
		x.cooldownSeconds,
		x.createdAt,
		x.updatedAt,
	), nil
}
