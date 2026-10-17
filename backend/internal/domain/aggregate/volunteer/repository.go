package aggregate_volunteer

import (
	"context"
	"errors"

	model_identity "github.com/Roum1212/hahaton/internal/domain/model/identity"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

var ErrNotFound = errors.New("volunteer not found")

type Repository interface {
	Create(ctx context.Context, volunteer *Volunteer, afterCreate func(ctx context.Context) error) error
	FindByIdentity(ctx context.Context, identity model_identity.Identity) (*Volunteer, error)
	List(ctx context.Context) ([]*Volunteer, error)
	Load(ctx context.Context, userID model_user_id.UserID) (*Volunteer, error)
	Modify(
		ctx context.Context,
		userID model_user_id.UserID,
		change func(
			ctx context.Context,
			volunteer *Volunteer,
		) error,
	) error
	Save(ctx context.Context, volunteer *Volunteer) error
}
