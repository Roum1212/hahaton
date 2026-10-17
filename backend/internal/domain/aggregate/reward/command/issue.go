package aggregate_reward_command

import (
	"time"
)

type Issue struct {
	now time.Time
}

func (x Issue) GetNow() time.Time {
	return x.now
}

func NewIssue(now time.Time) Issue {
	return Issue{
		now: now,
	}
}
