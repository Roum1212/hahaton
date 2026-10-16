package service_authenticate_session

import (
	"errors"
	"fmt"
	"strings"

	model_session_id "github.com/Roum1212/hahaton/internal/domain/model/session-id"
)

type Request struct {
	sessionID model_session_id.SessionID
	idToken   string
}

func (x Request) GetIDToken() string {
	return x.idToken
}

func (x Request) GetSessionID() model_session_id.SessionID {
	return x.sessionID
}

func (x Request) Validate() error {
	if err := x.sessionID.Validate(); err != nil {
		return fmt.Errorf("invalid session id: %w", err)
	}

	if len(strings.TrimSpace(x.idToken)) == 0 {
		return errors.New("id token is required")
	}

	return nil
}

func NewRequest(sessionID model_session_id.SessionID, idToken string) Request {
	return Request{
		sessionID: sessionID,
		idToken:   idToken,
	}
}
