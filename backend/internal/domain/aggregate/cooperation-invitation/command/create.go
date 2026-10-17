package aggregate_cooperation_invitation_command

import (
	"time"
)

type Create struct {
	token     string
	expiresAt time.Time
	now       time.Time
}

func (x Create) GetExpiresAt() time.Time {
	return x.expiresAt
}

func (x Create) GetNow() time.Time {
	return x.now
}

func (x Create) GetToken() string {
	return x.token
}

func NewCreate(token string, expiresAt, now time.Time) Create {
	return Create{
		token:     token,
		expiresAt: expiresAt,
		now:       now,
	}
}
