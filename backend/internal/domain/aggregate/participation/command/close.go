package aggregate_participation_command

import (
	"time"
)

type Close struct {
	now time.Time
}

func (x Close) GetNow() time.Time {
	return x.now
}

func NewClose(now time.Time) Close {
	return Close{
		now: now,
	}
}
