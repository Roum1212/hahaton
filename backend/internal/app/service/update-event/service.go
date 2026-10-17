package service_update_event

import (
	"context"
	"errors"
	"fmt"
	"time"

	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_event_command "github.com/Roum1212/hahaton/internal/domain/aggregate/event/command"
	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	model_event_schedule "github.com/Roum1212/hahaton/internal/domain/model/event-schedule"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	model_participation_status "github.com/Roum1212/hahaton/internal/domain/model/participation-status"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrNotFound         = errors.New("not found")
	ErrPermissionDenied = errors.New("permission denied")
)

type (
	ImageStorage interface {
		Exists(ctx context.Context, path model_image_path.ImagePath) (bool, error)
	}
	Service interface {
		UpdateEvent(ctx context.Context, r Request) error
	}
	service struct {
		events         aggregate_event.Repository
		coordinators   aggregate_coordinator.Repository
		participations aggregate_participation.Repository
		images         ImageStorage
	}
)

func (x service) UpdateEvent(ctx context.Context, r Request) error {
	if err := r.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	if _, err := x.coordinators.Load(ctx, r.GetUserID()); err != nil {
		switch {
		case errors.Is(err, aggregate_coordinator.ErrNotFound):
			return ErrPermissionDenied
		default:
			return err
		}
	}

	if path, ok := r.GetImagePath().Get(); ok {
		if err := x.checkImage(ctx, r.GetUserID(), path); err != nil {
			return err
		}
	}

	return x.events.Modify(ctx, r.GetEventID(), func(
		ctx context.Context,
		event *aggregate_event.Event,
	) error {
		return x.update(ctx, r, event)
	})
}

func (x service) checkImage(ctx context.Context, owner model_user_id.UserID, path model_image_path.ImagePath) error {
	if path.IsZero() {
		return nil
	}

	if path.GetUserID() != owner {
		return ErrPermissionDenied
	}

	exists, err := x.images.Exists(ctx, path)
	if err != nil {
		return err
	}

	if !exists {
		return ErrNotFound
	}

	return nil
}

func (x service) update(ctx context.Context, r Request, event *aggregate_event.Event) error {
	if event.GetCoordinatorID() != r.GetUserID() {
		return ErrPermissionDenied
	}

	participations, err := x.participations.List(ctx)
	if err != nil {
		return err
	}

	hasActive := false

	for _, value := range participations {
		if value.GetEventID() == event.GetEventID() && (value.GetStatus() == model_participation_status.Pending ||
			value.GetStatus() == model_participation_status.Approved) {
			hasActive = true
		}
	}

	title := event.GetTitle()
	if value, ok := r.GetTitle().Get(); ok {
		title = value
	}

	description := event.GetDescription()
	if value, ok := r.GetDescription().Get(); ok {
		description = value
	}

	path := event.GetImagePath()
	if value, ok := r.GetImagePath().Get(); ok {
		path = value
	}

	location := event.GetLocation()
	if value, ok := r.GetLocation().Get(); ok {
		location = value
	}

	point := event.GetPoint()
	if value, ok := r.GetPoint().Get(); ok {
		point = value
	}

	startsAt, endsAt := event.GetSchedule().GetStartsAt(), event.GetSchedule().GetEndsAt()
	if value, ok := r.GetStartsAt().Get(); ok {
		startsAt = value
	}

	if value, ok := r.GetEndsAt().Get(); ok {
		endsAt = value
	}

	schedule, err := model_event_schedule.NewEventSchedule(startsAt, endsAt)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	requirements := event.GetRequirements()
	if value, ok := r.GetRequirements().Get(); ok {
		requirements = value
	}

	limit := event.GetVolunteerLimit()
	if value, ok := r.GetVolunteerLimit().Get(); ok {
		limit = value
	}

	return event.Update(aggregate_event_command.NewUpdate(
		title,
		description,
		path,
		location,
		point,
		schedule,
		requirements,
		limit,
		time.Now().UTC(),
		hasActive,
	))
}

func NewService(
	events aggregate_event.Repository,
	coordinators aggregate_coordinator.Repository,
	participations aggregate_participation.Repository,
	images ImageStorage,
) Service {
	return service{
		events:         events,
		coordinators:   coordinators,
		participations: participations,
		images:         images,
	}
}
