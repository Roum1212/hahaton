package aggregate_event

import (
	"context"
	"errors"

	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
)

var (
	ErrNotFound      = errors.New("event not found")
	ErrAlreadyExists = errors.New("event already exists")
)

type Repository interface {
	Create(ctx context.Context, value *Event, afterCreate func(ctx context.Context) error) error
	List(ctx context.Context) ([]*Event, error)
	Load(ctx context.Context, eventID model_event_id.EventID) (*Event, error)
	Modify(
		ctx context.Context,
		eventID model_event_id.EventID,
		change func(
			ctx context.Context,
			value *Event,
		) error,
	) error
	Save(ctx context.Context, value *Event) error
}
