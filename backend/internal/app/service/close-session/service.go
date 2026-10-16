package service_close_session

import (
	"context"
	"errors"
	"fmt"

	aggregate_session "github.com/Roum1212/hahaton/internal/domain/aggregate/session"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

var ErrInvalidRequest = errors.New("invalid request")

type Service interface {
	CloseSession(ctx context.Context, r Request) error
}

type service struct {
	repository aggregate_session.Repository
}

func (x service) CloseSession(ctx context.Context, r Request) error {
	if err := r.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil
	}

	if err := x.repository.Delete(ctx, session.GetSessionID()); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}

func NewService(repository aggregate_session.Repository) Service {
	return service{repository: repository}
}
