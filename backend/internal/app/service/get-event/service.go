package service_get_event

import (
	"context"
	"errors"
	"fmt"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	model_cooperation_status "github.com/Roum1212/hahaton/internal/domain/model/cooperation-status"
	model_participation_status "github.com/Roum1212/hahaton/internal/domain/model/participation-status"
)

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrPermissionDenied = errors.New("permission denied")
)

type (
	Service interface {
		GetEvent(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		events         aggregate_event.Repository
		participations aggregate_participation.Repository
		cooperations   aggregate_cooperation.Repository
	}
)

func (x service) GetEvent(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	event, err := x.events.Load(ctx, r.GetEventID())
	if err != nil {
		return Response{}, err
	}

	participations, err := x.participations.List(ctx)
	if err != nil {
		return Response{}, err
	}

	allowed := event.GetCoordinatorID() == r.GetUserID()

	cooperation, err := x.cooperations.Load(ctx, event.GetCoordinatorID(), r.GetUserID())
	if err != nil && !errors.Is(err, aggregate_cooperation.ErrNotFound) {
		return Response{}, err
	}

	if err == nil && cooperation.GetStatus() == model_cooperation_status.Active {
		allowed = true
	}

	var approved, pending int32

	for _, participation := range participations {
		if participation.GetEventID() != event.GetEventID() {
			continue
		}

		if participation.GetVolunteerID() == r.GetUserID() {
			allowed = true
		}

		if participation.GetStatus() == model_participation_status.Approved {
			approved++
		}

		if participation.GetStatus() == model_participation_status.Pending {
			pending++
		}
	}

	if !allowed {
		return Response{}, ErrPermissionDenied
	}

	return NewResponse(event, approved, pending), nil
}

func NewService(
	events aggregate_event.Repository,
	participations aggregate_participation.Repository,
	cooperations aggregate_cooperation.Repository,
) Service {
	return service{
		events:         events,
		participations: participations,
		cooperations:   cooperations,
	}
}
