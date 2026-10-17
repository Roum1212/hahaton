package service_create_coordinator

import (
	"context"
	"errors"
	"fmt"
	"time"

	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_coordinator_command "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator/command"
	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
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
		CreateCoordinator(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		volunteers   aggregate_volunteer.Repository
		coordinators aggregate_coordinator.Repository
		images       ImageStorage
	}
)

func (x service) CreateCoordinator(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	if _, err := x.volunteers.Load(ctx, r.GetVolunteerID()); err != nil {
		return Response{}, err
	}

	path := r.GetAvatarPath()
	if err := x.checkImage(ctx, r.GetVolunteerID(), path); err != nil {
		return Response{}, err
	}

	coordinator := aggregate_coordinator.NewCoordinator(r.GetVolunteerID())
	if err := coordinator.Create(aggregate_coordinator_command.NewCreate(
		r.GetDescription(),
		path,
		r.GetCooldownSeconds(),
		time.Now().UTC(),
	)); err != nil {
		return Response{}, err
	}

	if err := x.coordinators.Create(ctx, coordinator); err != nil {
		return Response{}, err
	}

	return NewResponse(coordinator.GetUserID()), nil
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

func NewService(
	volunteers aggregate_volunteer.Repository,
	coordinators aggregate_coordinator.Repository,
	images ImageStorage,
) Service {
	return service{
		volunteers:   volunteers,
		coordinators: coordinators,
		images:       images,
	}
}
