package aggregate_participation_command

import (
	"time"
)

type RevokeChat struct {
	now time.Time
}

func (x RevokeChat) GetNow() time.Time {
	return x.now
}

func NewRevokeChat(now time.Time) RevokeChat {
	return RevokeChat{
		now: now,
	}
}
