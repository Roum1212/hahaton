package service_get_cooperation

import (
	"context"
	"errors"
	"fmt"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
)

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrPermissionDenied = errors.New("permission denied")
)

type (
	Service interface {
		GetCooperation(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		cooperations aggregate_cooperation.Repository
		coordinators aggregate_coordinator.Repository
	}
)

func (x service) GetCooperation(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	coordinatorID, volunteerID := r.GetCounterpartID(), r.GetUserID()
	if r.GetAsCoordinator() {
		if _, err := x.coordinators.Load(ctx, r.GetUserID()); err != nil {
			switch {
			case errors.Is(err, aggregate_coordinator.ErrNotFound):
				return Response{}, ErrPermissionDenied
			default:
				return Response{}, err
			}
		}

		coordinatorID, volunteerID = r.GetUserID(), r.GetCounterpartID()
	}

	value, err := x.cooperations.Load(ctx, coordinatorID, volunteerID)
	if err != nil {
		return Response{}, err
	}

	return NewResponse(value), nil
}

func NewService(cooperations aggregate_cooperation.Repository, coordinators aggregate_coordinator.Repository) Service {
	return service{
		cooperations: cooperations,
		coordinators: coordinators,
	}
}
