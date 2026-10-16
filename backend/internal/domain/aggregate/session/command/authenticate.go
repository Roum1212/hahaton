package aggregate_session_command

import (
	"time"

	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Authenticate struct {
	userID    model_user_id.UserID
	now       time.Time
	expiresAt time.Time
}

func (x Authenticate) GetExpiresAt() time.Time {
	return x.expiresAt
}

func (x Authenticate) GetNow() time.Time {
	return x.now
}

func (x Authenticate) GetUserID() model_user_id.UserID {
	return x.userID
}

func NewAuthenticate(
	userID model_user_id.UserID,
	now time.Time,
	expiresAt time.Time,
) Authenticate {
	return Authenticate{
		userID:    userID,
		now:       now,
		expiresAt: expiresAt,
	}
}
