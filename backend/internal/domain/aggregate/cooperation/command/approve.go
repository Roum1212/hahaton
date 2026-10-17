package aggregate_cooperation_command

import (
	"time"
)

type Approve struct {
	now time.Time
}

func (x Approve) GetNow() time.Time {
	return x.now
}

func NewApprove(now time.Time) Approve {
	return Approve{
		now: now,
	}
}
