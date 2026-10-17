package aggregate_participation

import (
	"context"
	"errors"

	model_participation_change "github.com/Roum1212/hahaton/internal/domain/model/participation-change"
	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
)

var ErrNotFound = errors.New("participation not found")

type Repository interface {
	AppendChange(ctx context.Context, value model_participation_change.ParticipationChange) error
	Create(ctx context.Context, value *Participation) error
	List(ctx context.Context) ([]*Participation, error)
	ListChanges(ctx context.Context) ([]model_participation_change.ParticipationChange, error)
	Load(ctx context.Context, participationID model_participation_id.ParticipationID) (*Participation, error)
	Save(ctx context.Context, value *Participation) error
}
