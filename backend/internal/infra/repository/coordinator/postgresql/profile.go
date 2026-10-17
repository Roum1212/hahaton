package repository_coordinator_postgresql

import (
	"context"
	"errors"

	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

func (x Coordinators) HasProfile(ctx context.Context, userID model_user_id.UserID) (bool, error) {
	_, err := x.Load(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, aggregate_coordinator.ErrNotFound):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
