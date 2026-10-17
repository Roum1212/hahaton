package aggregate_participation_command

import (
	"time"
)

type Reject struct {
	reason string
	now    time.Time
}

func (x Reject) GetNow() time.Time {
	return x.now
}

func (x Reject) GetReason() string {
	return x.reason
}

func NewReject(reason string, now time.Time) Reject {
	return Reject{
		reason: reason,
		now:    now,
	}
}
