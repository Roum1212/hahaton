package aggregate_cooperation_command

import (
	"time"
)

type Request struct {
	now time.Time
}

func (x Request) GetNow() time.Time {
	return x.now
}

func NewRequest(now time.Time) Request {
	return Request{
		now: now,
	}
}
