package aggregate_participation_command

import (
	"time"
)

type Restore struct {
	now time.Time
}

func (x Restore) GetNow() time.Time {
	return x.now
}

func NewRestore(now time.Time) Restore {
	return Restore{
		now: now,
	}
}
