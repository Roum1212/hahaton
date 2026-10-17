package service_get_coordinator_settings

import (
	"context"
	"errors"
	"fmt"

	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
)

var ErrInvalidRequest = errors.New("invalid request")

type (
	Service interface {
		GetCoordinatorSettings(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		coordinators aggregate_coordinator.Repository
	}
)

func (x service) GetCoordinatorSettings(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	coordinator, err := x.coordinators.Load(ctx, r.GetUserID())
	if err != nil {
		return Response{}, err
	}

	return NewResponse(coordinator.GetCooldownSeconds()), nil
}

func NewService(coordinators aggregate_coordinator.Repository) Service {
	return service{
		coordinators: coordinators,
	}
}
