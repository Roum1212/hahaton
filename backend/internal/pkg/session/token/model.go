package session_token

import (
	"time"

	model_session_id "github.com/Roum1212/hahaton/internal/domain/model/session-id"
)

type Token struct {
	sessionID model_session_id.SessionID
	issuedAt  time.Time
	expiresAt time.Time
}

func (x Token) GetExpiresAt() time.Time {
	return x.expiresAt
}

func (x Token) GetIssuedAt() time.Time {
	return x.issuedAt
}

func (x Token) GetSessionID() model_session_id.SessionID {
	return x.sessionID
}

func NewToken(
	sessionID model_session_id.SessionID,
	issuedAt time.Time,
	expiresAt time.Time,
) Token {
	return Token{
		sessionID: sessionID,
		issuedAt:  issuedAt.Truncate(time.Second),
		expiresAt: expiresAt.Truncate(time.Second),
	}
}
