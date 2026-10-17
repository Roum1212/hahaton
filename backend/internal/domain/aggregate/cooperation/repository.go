package aggregate_cooperation

import (
	"context"
	"errors"

	model_cooperation_change "github.com/Roum1212/hahaton/internal/domain/model/cooperation-change"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

var ErrNotFound = errors.New("cooperation not found")

type Repository interface {
	AppendChange(ctx context.Context, value model_cooperation_change.CooperationChange) error
	Create(ctx context.Context, value *Cooperation, afterCreate func(ctx context.Context) error) error
	List(ctx context.Context) ([]*Cooperation, error)
	Load(
		ctx context.Context,
		coordinatorID model_user_id.UserID,
		volunteerID model_user_id.UserID,
	) (*Cooperation, error)
	Modify(
		ctx context.Context,
		coordinatorID model_user_id.UserID,
		volunteerID model_user_id.UserID,
		change func(
			ctx context.Context,
			value *Cooperation,
		) error,
	) error
	Save(ctx context.Context, value *Cooperation) error
}
