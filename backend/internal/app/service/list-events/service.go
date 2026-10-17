package service_list_events

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	model_city "github.com/Roum1212/hahaton/internal/domain/model/city"
	model_cooperation_status "github.com/Roum1212/hahaton/internal/domain/model/cooperation-status"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_participation_status "github.com/Roum1212/hahaton/internal/domain/model/participation-status"
	model_search_area "github.com/Roum1212/hahaton/internal/domain/model/search-area"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	"github.com/Roum1212/hahaton/internal/pkg/page"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

var (
	ErrFailedPrecondition = errors.New("failed precondition")
	ErrInvalidRequest     = errors.New("invalid request")
	ErrNotFound           = errors.New("not found")
	ErrPermissionDenied   = errors.New("permission denied")
)

type Catalog interface {
	GetCity(ctx context.Context, id string) (model_city.City, bool, error)
}
type (
	Service interface {
		ListEvents(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		events         aggregate_event.Repository
		coordinators   aggregate_coordinator.Repository
		cooperations   aggregate_cooperation.Repository
		participations aggregate_participation.Repository
		catalog        Catalog
	}
)

func (x service) ListEvents(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	allowed, err := x.allowedCoordinators(ctx, r)
	if err != nil {
		return Response{}, err
	}

	city, err := x.searchCity(ctx, r.GetSearchArea())
	if err != nil {
		return Response{}, err
	}

	events, err := x.events.List(ctx)
	if err != nil {
		return Response{}, err
	}

	participations, err := x.participations.List(ctx)
	if err != nil {
		return Response{}, err
	}

	approved, pending := make(map[model_event_id.EventID]int32), make(map[model_event_id.EventID]int32)

	for _, value := range participations {
		if value.GetStatus() == model_participation_status.Approved {
			approved[value.GetEventID()]++
		}

		if value.GetStatus() == model_participation_status.Pending {
			pending[value.GetEventID()]++
		}
	}

	filtered := make([]*aggregate_event.Event, 0)

	for _, event := range events {
		if matches(event, r, allowed, city) {
			filtered = append(filtered, event)
		}
	}

	slices.SortFunc(filtered, func(
		a *aggregate_event.Event,
		b *aggregate_event.Event,
	) int {
		if value := a.GetSchedule().GetStartsAt().Compare(b.GetSchedule().GetStartsAt()); value != 0 {
			return value
		}

		return strings.Compare(a.GetEventID().String(), b.GetEventID().String())
	})
	values, total := page.Slice(filtered, r.GetPage())

	return NewResponse(values, approved, pending, total), nil
}

func (x service) allowedCoordinators(ctx context.Context, r Request) (map[model_user_id.UserID]bool, error) {
	allowed := make(map[model_user_id.UserID]bool)

	if r.GetScope() == scopeOwned {
		if _, err := x.coordinators.Load(ctx, r.GetUserID()); err != nil {
			switch {
			case errors.Is(err, aggregate_coordinator.ErrNotFound):
				return nil, ErrPermissionDenied
			default:
				return nil, err
			}
		}

		allowed[r.GetUserID()] = true

		return allowed, nil
	}

	cooperations, err := x.cooperations.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, value := range cooperations {
		if value.GetVolunteerID() == r.GetUserID() &&
			value.GetStatus() == model_cooperation_status.Active {
			allowed[value.GetCoordinatorID()] = true
		}
	}

	return allowed, nil
}

func matches(
	event *aggregate_event.Event,
	r Request,
	allowed map[model_user_id.UserID]bool,
	city model_city.City,
) bool {
	if !allowed[event.GetCoordinatorID()] {
		return false
	}

	if value, ok := r.GetCoordinatorID().Get(); ok && event.GetCoordinatorID() != value {
		return false
	}

	if value, ok := r.GetStatus().Get(); ok && event.GetStatus() != value {
		return false
	}

	if !within(event.GetSchedule().GetStartsAt(), r.GetStartsAtFrom(), r.GetStartsAtTo()) {
		return false
	}

	if !within(event.GetSchedule().GetEndsAt(), r.GetEndsAtFrom(), r.GetEndsAtTo()) {
		return false
	}

	area, ok := r.GetSearchArea().Get()
	if !ok {
		return true
	}

	if len(area.GetCityID()) != 0 {
		return city.Contains(event.GetPoint())
	}

	return area.Contains(event.GetPoint())
}

func (x service) searchCity(
	ctx context.Context,
	option utils_value.Optional[model_search_area.SearchArea],
) (model_city.City, error) {
	area, ok := option.Get()
	if !ok || len(area.GetCityID()) == 0 {
		return model_city.City{}, nil
	}

	city, found, err := x.catalog.GetCity(ctx, area.GetCityID())
	if err != nil {
		return model_city.City{}, err
	}

	if !found {
		return model_city.City{}, ErrNotFound
	}

	if !city.HasBoundary() {
		return model_city.City{}, ErrFailedPrecondition
	}

	return city, nil
}

func within(value time.Time, from, to utils_value.Optional[time.Time]) bool {
	fromValue, hasFrom := from.Get()
	toValue, hasTo := to.Get()

	return (!hasFrom || !value.Before(fromValue)) && (!hasTo || value.Before(toValue))
}

func NewService(
	events aggregate_event.Repository,
	coordinators aggregate_coordinator.Repository,
	cooperations aggregate_cooperation.Repository,
	participations aggregate_participation.Repository,
	catalog Catalog,
) Service {
	return service{
		events:         events,
		coordinators:   coordinators,
		cooperations:   cooperations,
		participations: participations,
		catalog:        catalog,
	}
}
