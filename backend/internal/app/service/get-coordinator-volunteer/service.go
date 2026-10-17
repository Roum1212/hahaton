package service_get_coordinator_volunteer

import (
	"context"
	"errors"
	"fmt"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	aggregate_reward_order "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order"
	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
	model_cooperation_status "github.com/Roum1212/hahaton/internal/domain/model/cooperation-status"
)

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrPermissionDenied = errors.New("permission denied")
)

type (
	Service interface {
		GetCoordinatorVolunteer(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		coordinators   aggregate_coordinator.Repository
		volunteers     aggregate_volunteer.Repository
		cooperations   aggregate_cooperation.Repository
		participations aggregate_participation.Repository
		orders         aggregate_reward_order.Repository
	}
)

func (x service) GetCoordinatorVolunteer(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	if _, err := x.coordinators.Load(ctx, r.GetUserID()); err != nil {
		switch {
		case errors.Is(err, aggregate_coordinator.ErrNotFound):
			return Response{}, ErrPermissionDenied
		default:
			return Response{}, err
		}
	}

	allowed, contacts := false, false

	cooperation, err := x.cooperations.Load(ctx, r.GetUserID(), r.GetVolunteerID())
	if err != nil && !errors.Is(err, aggregate_cooperation.ErrNotFound) {
		return Response{}, err
	}

	if err == nil {
		allowed = true
		contacts = cooperation.GetStatus() == model_cooperation_status.Active
	}

	participations, err := x.participations.List(ctx)
	if err != nil {
		return Response{}, err
	}

	for _, value := range participations {
		if value.GetCoordinatorID() == r.GetUserID() && value.GetVolunteerID() == r.GetVolunteerID() {
			allowed = true
		}
	}

	orders, err := x.orders.List(ctx)
	if err != nil {
		return Response{}, err
	}

	for _, value := range orders {
		if value.GetCoordinatorID() == r.GetUserID() && value.GetVolunteerID() == r.GetVolunteerID() {
			allowed = true
		}
	}

	if !allowed {
		return Response{}, ErrPermissionDenied
	}

	volunteer, err := x.volunteers.Load(ctx, r.GetVolunteerID())
	if err != nil {
		return Response{}, err
	}

	return NewResponse(volunteer, contacts), nil
}

func NewService(
	coordinators aggregate_coordinator.Repository,
	volunteers aggregate_volunteer.Repository,
	cooperations aggregate_cooperation.Repository,
	participations aggregate_participation.Repository,
	orders aggregate_reward_order.Repository,
) Service {
	return service{
		coordinators:   coordinators,
		volunteers:     volunteers,
		cooperations:   cooperations,
		participations: participations,
		orders:         orders,
	}
}
