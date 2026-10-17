package service_find_volunteers

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
		FindVolunteers(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		volunteers   aggregate_volunteer.Repository
		coordinators aggregate_coordinator.Repository
	}
)

func (x service) FindVolunteers(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	volunteers, err := x.volunteers.List(ctx)
	if err != nil {
		return Response{}, err
	}

	coordinators, err := x.coordinators.List(ctx)
	if err != nil {
		return Response{}, err
	}

	profiles := make(map[model_user_id.UserID]bool, len(coordinators))
	for _, value := range coordinators {
		profiles[value.GetUserID()] = true
	}

	filtered := make([]*aggregate_volunteer.Volunteer, 0)

	for _, value := range volunteers {
		if filter, ok := r.GetVolunteerID().Get(); ok && value.GetUserID() != filter {
			continue
		}

		if filter, ok := r.GetEmailAddress().Get(); ok && !value.GetEmailAddress().Equals(filter) {
			continue
		}

		if filter, ok := r.GetPhoneNumber().Get(); ok && value.GetPhoneNumber() != filter {
			continue
		}

		filtered = append(filtered, value)
	}

	slices.SortFunc(filtered, func(
		a *aggregate_volunteer.Volunteer,
		b *aggregate_volunteer.Volunteer,
	) int {
		return strings.Compare(a.GetUserID().String(), b.GetUserID().String())
	})
	values, total := page.Slice(filtered, r.GetPage())

	return NewResponse(values, profiles, total), nil
}

func NewService(volunteers aggregate_volunteer.Repository, coordinators aggregate_coordinator.Repository) Service {
	return service{
		volunteers:   volunteers,
		coordinators: coordinators,
	}
}
