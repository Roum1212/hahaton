package aggregate_cooperation_command

import (
	"time"
)

type Reject struct {
	retryAfter time.Time
	now        time.Time
}

func (x Reject) GetNow() time.Time {
	return x.now
}

func (x Reject) GetRetryAfter() time.Time {
	return x.retryAfter
}

func NewReject(retryAfter, now time.Time) Reject {
	return Reject{
		retryAfter: retryAfter,
		now:        now,
	}
}
