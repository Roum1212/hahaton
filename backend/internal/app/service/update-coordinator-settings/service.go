package service_update_coordinator_settings

import (
	"context"
	"errors"
	"fmt"
	"time"

	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_coordinator_command "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator/command"
)

var ErrInvalidRequest = errors.New("invalid request")

type (
	Service interface {
		UpdateCoordinatorSettings(ctx context.Context, r Request) error
	}
	service struct {
		coordinators aggregate_coordinator.Repository
	}
)

func (x service) UpdateCoordinatorSettings(ctx context.Context, r Request) error {
	if err := r.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	return x.coordinators.Modify(ctx, r.GetUserID(), func(
		ctx context.Context,
		coordinator *aggregate_coordinator.Coordinator,
	) error {
		cooldown := coordinator.GetCooldownSeconds()
		if value, ok := r.GetCooldownSeconds().Get(); ok {
			cooldown = value
		}

		coordinator.Update(aggregate_coordinator_command.NewUpdate(
			coordinator.GetDescription(),
			coordinator.GetAvatarPath(),
			cooldown,
			time.Now().UTC(),
		))

		return nil
	})
}

func NewService(coordinators aggregate_coordinator.Repository) Service {
	return service{
		coordinators: coordinators,
	}
}
