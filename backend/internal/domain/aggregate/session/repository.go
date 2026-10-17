package aggregate_session

import (
	"context"
	"errors"

	model_session_id "github.com/Roum1212/hahaton/internal/domain/model/session-id"
)

var ErrNotFound = errors.New("session not found")

type Repository interface {
	Create(ctx context.Context, session *Session) error
	Delete(ctx context.Context, sessionID model_session_id.SessionID) error
	Load(ctx context.Context, sessionID model_session_id.SessionID) (*Session, error)
	Replace(
		ctx context.Context,
		sessionID model_session_id.SessionID,
		replace func(
			ctx context.Context,
			current *Session,
		) (*Session, error),
	) error
	Save(ctx context.Context, session *Session) error
}
