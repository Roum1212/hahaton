package service_list_participations

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	"github.com/Roum1212/hahaton/internal/pkg/page"
)

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrPermissionDenied = errors.New("permission denied")
)

type (
	Service interface {
		ListParticipations(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		participations aggregate_participation.Repository
		coordinators   aggregate_coordinator.Repository
	}
)

func (x service) ListParticipations(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	if r.GetScope() == scopeMyEvents {
		if _, err := x.coordinators.Load(ctx, r.GetUserID()); err != nil {
			switch {
			case errors.Is(err, aggregate_coordinator.ErrNotFound):
				return Response{}, ErrPermissionDenied
			default:
				return Response{}, err
			}
		}
	}

	values, err := x.participations.List(ctx)
	if err != nil {
		return Response{}, err
	}

	filtered := make([]*aggregate_participation.Participation, 0)

	for _, value := range values {
		if r.GetScope() == scopeMyApplications && value.GetVolunteerID() != r.GetUserID() {
			continue
		}

		if r.GetScope() == scopeMyEvents && value.GetCoordinatorID() != r.GetUserID() {
			continue
		}

		if filter, ok := r.GetEventID().Get(); ok && value.GetEventID() != filter {
			continue
		}

		if filter, ok := r.GetVolunteerID().Get(); ok && value.GetVolunteerID() != filter {
			continue
		}

		if filter, ok := r.GetStatus().Get(); ok && value.GetStatus() != filter {
			continue
		}

		filtered = append(filtered, value)
	}

	slices.SortFunc(filtered, func(
		a *aggregate_participation.Participation,
		b *aggregate_participation.Participation,
	) int {
		if value := a.GetCreatedAt().Compare(b.GetCreatedAt()); value != 0 {
			return value
		}

		return strings.Compare(a.GetParticipationID().String(), b.GetParticipationID().String())
	})
	values, total := page.Slice(filtered, r.GetPage())

	return NewResponse(values, total), nil
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
