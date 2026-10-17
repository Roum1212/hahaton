package service_create_event

import (
	"context"
	"errors"
	"fmt"
	"time"

	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_event_chat "github.com/Roum1212/hahaton/internal/domain/aggregate/event-chat"
	aggregate_event_command "github.com/Roum1212/hahaton/internal/domain/aggregate/event/command"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
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
		CreateEvent(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		events       aggregate_event.Repository
		coordinators aggregate_coordinator.Repository
		chats        aggregate_event_chat.Repository
		images       ImageStorage
	}
)

func (x service) CreateEvent(ctx context.Context, r Request) (Response, error) {
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

	path := r.GetImagePath()
	if err := x.checkImage(ctx, r.GetUserID(), path); err != nil {
		return Response{}, err
	}

	id, err := model_event_id.GenerateEventID()
	if err != nil {
		return Response{}, err
	}

	event := aggregate_event.NewEvent(id, r.GetUserID())

	command := aggregate_event_command.NewCreate(
		r.GetTitle(),
		r.GetDescription(),
		path,
		r.GetLocation(),
		r.GetPoint(),
		r.GetSchedule(),
		r.GetRequirements(),
		r.GetVolunteerLimit(),
		time.Now().UTC(),
	)
	if err := event.Create(command); err != nil {
		return Response{}, err
	}

	if err := x.events.Create(ctx, event, func(
		ctx context.Context,
	) error {
		return x.chats.Create(ctx, aggregate_event_chat.NewEventChat(id))
	}); err != nil {
		return Response{}, err
	}

	return NewResponse(id), nil
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

func NewService(
	events aggregate_event.Repository,
	coordinators aggregate_coordinator.Repository,
	chats aggregate_event_chat.Repository,
	images ImageStorage,
) Service {
	return service{
		events:       events,
		coordinators: coordinators,
		chats:        chats,
		images:       images,
	}
}
