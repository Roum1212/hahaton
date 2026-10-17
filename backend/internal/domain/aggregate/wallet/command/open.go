package aggregate_wallet_command

import (
	"time"
)

type Open struct {
	now time.Time
}

func (x Open) GetNow() time.Time {
	return x.now
}

func NewOpen(now time.Time) Open {
	return Open{now: now}
}
