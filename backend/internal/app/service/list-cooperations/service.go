package service_list_cooperations

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	"github.com/Roum1212/hahaton/internal/pkg/page"
)

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrPermissionDenied = errors.New("permission denied")
)

type (
	Service interface {
		ListCooperations(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		cooperations aggregate_cooperation.Repository
		coordinators aggregate_coordinator.Repository
	}
)

func (x service) ListCooperations(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	if r.GetScope() == scopeAsCoordinator {
		if _, err := x.coordinators.Load(ctx, r.GetUserID()); err != nil {
			switch {
			case errors.Is(err, aggregate_coordinator.ErrNotFound):
				return Response{}, ErrPermissionDenied
			default:
				return Response{}, err
			}
		}
	}

	values, err := x.cooperations.List(ctx)
	if err != nil {
		return Response{}, err
	}

	filtered := make([]*aggregate_cooperation.Cooperation, 0)

	for _, value := range values {
		if r.GetScope() == scopeAsVolunteer && value.GetVolunteerID() != r.GetUserID() {
			continue
		}

		if r.GetScope() == scopeAsCoordinator && value.GetCoordinatorID() != r.GetUserID() {
			continue
		}

		if filter, ok := r.GetStatus().Get(); ok && value.GetStatus() != filter {
			continue
		}

		filtered = append(filtered, value)
	}

	slices.SortFunc(filtered, func(
		a *aggregate_cooperation.Cooperation,
		b *aggregate_cooperation.Cooperation,
	) int {
		if value := strings.Compare(a.GetCoordinatorID().String(), b.GetCoordinatorID().String()); value != 0 {
			return value
		}

		return strings.Compare(a.GetVolunteerID().String(), b.GetVolunteerID().String())
	})
	values, total := page.Slice(filtered, r.GetPage())

	return NewResponse(values, total), nil
}

func NewService(cooperations aggregate_cooperation.Repository, coordinators aggregate_coordinator.Repository) Service {
	return service{
		cooperations: cooperations,
		coordinators: coordinators,
	}
}
