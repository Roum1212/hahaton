package aggregate_reward_command

import (
	"time"
)

type Delete struct {
	now time.Time
}

func (x Delete) GetNow() time.Time {
	return x.now
}

func NewDelete(now time.Time) Delete {
	return Delete{
		now: now,
	}
}
