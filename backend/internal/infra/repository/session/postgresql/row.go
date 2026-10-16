package repository_session_postgresql

import (
	"fmt"
	"time"

	aggregate_session "github.com/Roum1212/hahaton/internal/domain/aggregate/session"
	model_session_id "github.com/Roum1212/hahaton/internal/domain/model/session-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type sessionRow struct {
	userID    string
	createdAt time.Time
	updatedAt time.Time
	expiresAt time.Time
}

func (x sessionRow) toSession(sessionID model_session_id.SessionID) (*aggregate_session.Session, error) {
	var userID model_user_id.UserID

	if len(x.userID) != 0 {
		value, err := model_user_id.NewUserID(x.userID)
		if err != nil {
			return nil, fmt.Errorf("restore session user id: %w", err)
		}

		userID = value
	}

	return aggregate_session.NewSessionFromFields(
		sessionID,
		userID,
		x.createdAt,
		x.updatedAt,
		x.expiresAt,
	), nil
}
