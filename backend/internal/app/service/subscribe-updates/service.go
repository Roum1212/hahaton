package service_subscribe_updates

import (
	"context"
	"errors"
	"fmt"
	"time"

	aggregate_session "github.com/Roum1212/hahaton/internal/domain/aggregate/session"
)

var ErrInvalidRequest = errors.New("invalid request")

type (
	Service interface {
		SubscribeUpdates(ctx context.Context, r Request) error
	}
	service struct{ sessions aggregate_session.Repository }
)

func (x service) SubscribeUpdates(ctx context.Context, r Request) error {
	if err := r.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	session, err := x.sessions.Load(ctx, r.GetSessionID())
	if err != nil {
		switch {
		case errors.Is(err, aggregate_session.ErrNotFound):
			return aggregate_session.ErrInvalidCredentials
		default:
			return err
		}
	}

	if session.GetUserID() != r.GetUserID() || !session.IsActive(time.Now().UTC()) {
		return aggregate_session.ErrInvalidCredentials
	}

	return nil
}

func NewService(sessions aggregate_session.Repository) Service { return service{sessions: sessions} }
