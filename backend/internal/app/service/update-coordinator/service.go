package service_update_coordinator

import (
	"context"
	"errors"
	"fmt"
	"time"

	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_coordinator_command "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator/command"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrNotFound         = errors.New("not found")
	ErrPermissionDenied = errors.New("permission denied")
)

type (
	ImageStorage interface {
		Exists(ctx context.Context, path model_image_path.ImagePath) (bool, error)
	}
	Service interface {
		UpdateCoordinator(ctx context.Context, r Request) error
	}
	service struct {
		coordinators aggregate_coordinator.Repository
		images       ImageStorage
	}
)

func (x service) UpdateCoordinator(ctx context.Context, r Request) error {
	if err := r.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	if path, ok := r.GetAvatarPath().Get(); ok {
		if err := x.checkImage(ctx, r.GetUserID(), path); err != nil {
			return err
		}
	}

	return x.coordinators.Modify(ctx, r.GetUserID(), func(
		ctx context.Context,
		coordinator *aggregate_coordinator.Coordinator,
	) error {
		description := coordinator.GetDescription()
		if value, ok := r.GetDescription().Get(); ok {
			description = value
		}

		path := coordinator.GetAvatarPath()
		if value, ok := r.GetAvatarPath().Get(); ok {
			path = value
		}

		coordinator.Update(aggregate_coordinator_command.NewUpdate(
			description,
			path,
			coordinator.GetCooldownSeconds(),
			time.Now().UTC(),
		))

		return nil
	})
}

func (x service) checkImage(ctx context.Context, owner model_user_id.UserID, path model_image_path.ImagePath) error {
	if path.IsZero() {
		return nil
	}

	if path.GetUserID() != owner {
		return ErrPermissionDenied
	}

	exists, err := x.images.Exists(ctx, path)
	if err != nil {
		return err
	}

	if !exists {
		return ErrNotFound
	}

	return nil
}

func NewService(coordinators aggregate_coordinator.Repository, images ImageStorage) Service {
	return service{
		coordinators: coordinators,
		images:       images,
	}
}
