package aggregate_event_command

import (
	"time"
)

type Cancel struct {
	now time.Time
}

func (x Cancel) GetNow() time.Time {
	return x.now
}

func NewCancel(now time.Time) Cancel {
	return Cancel{
		now: now,
	}
}
