package aggregate_reward_order_command

import (
	"time"
)

type Cancel struct {
	reason string
	now    time.Time
}

func (x Cancel) GetNow() time.Time {
	return x.now
}

func (x Cancel) GetReason() string {
	return x.reason
}

func NewCancel(reason string, now time.Time) Cancel {
	return Cancel{
		reason: reason,
		now:    now,
	}
}
