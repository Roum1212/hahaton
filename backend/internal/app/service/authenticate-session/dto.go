package service_authenticate_session

import (
	"fmt"

	model_session_id "github.com/Roum1212/hahaton/internal/domain/model/session-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Request struct {
	sessionID model_session_id.SessionID
	userID    model_user_id.UserID
}

func (x Request) GetSessionID() model_session_id.SessionID {
	return x.sessionID
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.sessionID.Validate(); err != nil {
		return fmt.Errorf("invalid session id: %w", err)
	}

	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	return nil
}

func NewRequest(sessionID model_session_id.SessionID, userID model_user_id.UserID) Request {
	return Request{
		sessionID: sessionID,
		userID:    userID,
	}
}
