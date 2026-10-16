package service_get_session

import (
	"context"
	"errors"
	"fmt"
	"time"

	aggregate_session "github.com/Roum1212/hahaton/internal/domain/aggregate/session"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

var ErrInvalidRequest = errors.New("invalid request")

type CoordinatorReader interface {
	HasProfile(ctx context.Context, userID model_user_id.UserID) (bool, error)
}

type Service interface {
	GetSession(ctx context.Context, r Request) (Response, error)
}

type service struct {
	coordinators CoordinatorReader
}

func (x service) GetSession(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	session, ok := session_context.GetSession(ctx)
	if !ok || !session.IsActive(time.Now().UTC()) {
		return Response{}, aggregate_session.ErrInvalidCredentials
	}

	if session.IsAnonymous() {
		return NewResponse(session, false), nil
	}

	hasCoordinatorProfile, err := x.coordinators.HasProfile(ctx, session.GetUserID())
	if err != nil {
		return Response{}, err
	}

	return NewResponse(session, hasCoordinatorProfile), nil
}

func NewService(coordinators CoordinatorReader) Service {
	return service{
		coordinators: coordinators,
	}
}
