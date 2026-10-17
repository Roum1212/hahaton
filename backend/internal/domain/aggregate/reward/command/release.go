package aggregate_reward_command

import (
	"time"
)

type Release struct {
	now time.Time
}

func (x Release) GetNow() time.Time {
	return x.now
}

func NewRelease(now time.Time) Release {
	return Release{
		now: now,
	}
}
