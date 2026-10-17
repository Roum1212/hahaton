package aggregate_cooperation_command

import (
	"time"
)

type AcceptInvitation struct {
	now time.Time
}

func (x AcceptInvitation) GetNow() time.Time {
	return x.now
}

func NewAcceptInvitation(now time.Time) AcceptInvitation {
	return AcceptInvitation{
		now: now,
	}
}
