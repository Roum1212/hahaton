package service_list_coordinators

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	"github.com/Roum1212/hahaton/internal/pkg/page"
)

var ErrInvalidRequest = errors.New("invalid request")

type (
	Service interface {
		ListCoordinators(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		coordinators aggregate_coordinator.Repository
		volunteers   aggregate_volunteer.Repository
	}
)

func (x service) ListCoordinators(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	values, err := x.coordinators.List(ctx)
	if err != nil {
		return Response{}, err
	}

	profiles, err := x.volunteers.List(ctx)
	if err != nil {
		return Response{}, err
	}

	volunteers := make(map[model_user_id.UserID]*aggregate_volunteer.Volunteer, len(profiles))
	for _, value := range profiles {
		volunteers[value.GetUserID()] = value
	}

	slices.SortFunc(values, func(
		a *aggregate_coordinator.Coordinator,
		b *aggregate_coordinator.Coordinator,
	) int {
		if result := strings.Compare(
			displayName(volunteers[a.GetUserID()]),
			displayName(volunteers[b.GetUserID()]),
		); result != 0 {
			return result
		}

		return strings.Compare(a.GetUserID().String(), b.GetUserID().String())
	})
	values, total := page.Slice(values, r.GetPage())

	return NewResponse(values, volunteers, total), nil
}

func displayName(volunteer *aggregate_volunteer.Volunteer) string {
	name := volunteer.GetFullName()

	return strings.TrimSpace(name.GetLastName() + " " + name.GetFirstName() + " " + name.GetMiddleName())
}

func NewService(coordinators aggregate_coordinator.Repository, volunteers aggregate_volunteer.Repository) Service {
	return service{
		coordinators: coordinators,
		volunteers:   volunteers,
	}
}
