package service_get_coordinator

import (
	"context"
	"errors"
	"fmt"

	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
)

var ErrInvalidRequest = errors.New("invalid request")

type (
	Service interface {
		GetCoordinator(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		coordinators aggregate_coordinator.Repository
		volunteers   aggregate_volunteer.Repository
	}
)

func (x service) GetCoordinator(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	id := r.GetUserID()
	if value, ok := r.GetID().Get(); ok {
		id = value
	}

	coordinator, err := x.coordinators.Load(ctx, id)
	if err != nil {
		return Response{}, err
	}

	volunteer, err := x.volunteers.Load(ctx, id)
	if err != nil {
		return Response{}, err
	}

	return NewResponse(coordinator, volunteer), nil
}

func NewService(coordinators aggregate_coordinator.Repository, volunteers aggregate_volunteer.Repository) Service {
	return service{
		coordinators: coordinators,
		volunteers:   volunteers,
	}
}
