package aggregate_coordinator

import (
	"context"
	"errors"

	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

var ErrNotFound = errors.New("coordinator not found")

type Repository interface {
	Create(ctx context.Context, value *Coordinator) error
	List(ctx context.Context) ([]*Coordinator, error)
	Load(ctx context.Context, userID model_user_id.UserID) (*Coordinator, error)
	Modify(
		ctx context.Context,
		userID model_user_id.UserID,
		change func(
			ctx context.Context,
			value *Coordinator,
		) error,
	) error
	Save(ctx context.Context, value *Coordinator) error
}
