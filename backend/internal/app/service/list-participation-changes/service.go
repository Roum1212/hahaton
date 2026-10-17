package service_list_participation_changes

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"

	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	model_participation_change "github.com/Roum1212/hahaton/internal/domain/model/participation-change"
	"github.com/Roum1212/hahaton/internal/pkg/page"
)

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrPermissionDenied = errors.New("permission denied")
)

type (
	Service interface {
		ListParticipationChanges(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		participations aggregate_participation.Repository
		coordinators   aggregate_coordinator.Repository
	}
)

func (x service) ListParticipationChanges(ctx context.Context, r Request) (Response, error) {
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

	values, err := x.participations.ListChanges(ctx)
	if err != nil {
		return Response{}, err
	}

	filtered := make([]model_participation_change.ParticipationChange, 0)

	for _, change := range values {
		if change.GetParticipationID() == r.GetParticipationID() {
			filtered = append(filtered, change)
		}
	}

	slices.SortFunc(filtered, func(
		a model_participation_change.ParticipationChange,
		b model_participation_change.ParticipationChange,
	) int {
		return cmp.Compare(a.GetSequence(), b.GetSequence())
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
