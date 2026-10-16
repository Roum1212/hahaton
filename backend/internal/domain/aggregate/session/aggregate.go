package aggregate_session

import (
	"errors"
	"time"

	aggregate_session_command "github.com/Roum1212/hahaton/internal/domain/aggregate/session/command"
	model_session_id "github.com/Roum1212/hahaton/internal/domain/model/session-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

var (
	ErrAlreadyAuthenticated = errors.New("session is already authenticated")
	ErrInvalidCredentials   = errors.New("invalid session credentials")
	ErrInvalidState         = errors.New("invalid session state")
)

type Session struct {
	sessionID model_session_id.SessionID
	userID    model_user_id.UserID
	createdAt time.Time
	updatedAt time.Time
	expiresAt time.Time
}

func (x *Session) GetCreatedAt() time.Time {
	return x.createdAt
}

func (x *Session) GetExpiresAt() time.Time {
	return x.expiresAt
}

func (x *Session) GetSessionID() model_session_id.SessionID {
	return x.sessionID
}

func (x *Session) GetUpdatedAt() time.Time {
	return x.updatedAt
}

func (x *Session) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x *Session) IsActive(now time.Time) bool {
	return !x.createdAt.IsZero() && now.Before(x.expiresAt)
}

func (x *Session) IsAnonymous() bool {
	return x.userID.IsZero()
}

func (x *Session) Authenticate(command aggregate_session_command.Authenticate) error {
	if !x.IsActive(command.GetNow()) {
		return ErrInvalidCredentials
	}

	if !x.userID.IsZero() {
		return ErrAlreadyAuthenticated
	}

	x.userID = command.GetUserID()
	x.updatedAt = command.GetNow()
	x.expiresAt = command.GetExpiresAt()

	return nil
}

func (x *Session) Open(command aggregate_session_command.Open) error {
	if !x.createdAt.IsZero() {
		return ErrInvalidState
	}

	x.createdAt = command.GetNow()
	x.updatedAt = command.GetNow()
	x.expiresAt = command.GetExpiresAt()

	return nil
}

func NewSession(sessionID model_session_id.SessionID) *Session {
	return &Session{sessionID: sessionID}
}

func NewSessionFromFields(
	sessionID model_session_id.SessionID,
	userID model_user_id.UserID,
	createdAt time.Time,
	updatedAt time.Time,
	expiresAt time.Time,
) *Session {
	return &Session{
		sessionID: sessionID,
		userID:    userID,
		createdAt: createdAt,
		updatedAt: updatedAt,
		expiresAt: expiresAt,
	}
}
