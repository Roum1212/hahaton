package service_get_participation

import (
	"context"
	"errors"
	"fmt"

	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
)

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrPermissionDenied = errors.New("permission denied")
)

type (
	Service interface {
		GetParticipation(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		participations aggregate_participation.Repository
		coordinators   aggregate_coordinator.Repository
	}
)

func (x service) GetParticipation(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	value, err := x.participations.Load(ctx, r.GetParticipationID())
	if err != nil {
		return Response{}, err
	}

	if value.GetVolunteerID() != r.GetUserID() {
		if value.GetCoordinatorID() != r.GetUserID() {
			return Response{}, ErrPermissionDenied
		}

		if _, err := x.coordinators.Load(ctx, r.GetUserID()); err != nil {
			switch {
			case errors.Is(err, aggregate_coordinator.ErrNotFound):
				return Response{}, ErrPermissionDenied
			default:
				return Response{}, err
			}
		}
	}

	return NewResponse(value), nil
}

func NewService(
	participations aggregate_participation.Repository,
	coordinators aggregate_coordinator.Repository,
) Service {
	return service{
		participations: participations,
		coordinators:   coordinators,
	}
}
