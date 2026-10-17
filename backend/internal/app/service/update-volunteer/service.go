package service_update_volunteer

import (
	"context"
	"errors"
	"fmt"
	"time"

	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
	aggregate_volunteer_command "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer/command"
	model_full_name "github.com/Roum1212/hahaton/internal/domain/model/full-name"
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
		UpdateVolunteer(ctx context.Context, r Request) error
	}
	service struct {
		volunteers aggregate_volunteer.Repository
		images     ImageStorage
	}
)

func (x service) UpdateVolunteer(ctx context.Context, r Request) error {
	if err := r.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	if path, ok := r.GetAvatarPath().Get(); ok {
		if err := x.checkImage(ctx, r.GetUserID(), path); err != nil {
			return err
		}
	}

	return x.volunteers.Modify(ctx, r.GetUserID(), func(
		ctx context.Context,
		volunteer *aggregate_volunteer.Volunteer,
	) error {
		fullName := volunteer.GetFullName()
		first := fullName.GetFirstName()
		last := fullName.GetLastName()
		middle := fullName.GetMiddleName()

		if value, ok := r.GetFirstName().Get(); ok {
			first = value
		}

		if value, ok := r.GetLastName().Get(); ok {
			last = value
		}

		if value, ok := r.GetMiddleName().Get(); ok {
			middle = value
		}

		name, err := model_full_name.NewFullName(first, last, middle)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrInvalidRequest, err)
		}

		dateOfBirth := volunteer.GetDateOfBirth()
		if value, ok := r.GetDateOfBirth().Get(); ok {
			dateOfBirth = value
		}

		path := volunteer.GetAvatarPath()
		if value, ok := r.GetAvatarPath().Get(); ok {
			path = value
		}

		volunteer.Update(aggregate_volunteer_command.NewUpdate(name, dateOfBirth, path, time.Now().UTC()))

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

func NewService(volunteers aggregate_volunteer.Repository, images ImageStorage) Service {
	return service{
		volunteers: volunteers,
		images:     images,
	}
}
