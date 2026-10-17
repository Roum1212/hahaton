package aggregate_cooperation_command

import (
	"time"
)

type Terminate struct {
	now time.Time
}

func (x Terminate) GetNow() time.Time {
	return x.now
}

func NewTerminate(now time.Time) Terminate {
	return Terminate{
		now: now,
	}
}
