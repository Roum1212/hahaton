package session_context

import (
	"context"

	"connectrpc.com/authn"

	aggregate_session "github.com/Roum1212/hahaton/internal/domain/aggregate/session"
)

func GetSession(ctx context.Context) (*aggregate_session.Session, bool) {
	session, ok := authn.GetInfo(ctx).(*aggregate_session.Session)

	return session, ok
}

func WithSession(ctx context.Context, session *aggregate_session.Session) context.Context {
	return authn.SetInfo(ctx, session)
}
