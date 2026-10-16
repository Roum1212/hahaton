package aggregate_session_command

import (
	"time"
)

type Open struct {
	now       time.Time
	expiresAt time.Time
}

func (x Open) GetExpiresAt() time.Time {
	return x.expiresAt
}

func (x Open) GetNow() time.Time {
	return x.now
}

func NewOpen(
	now time.Time,
	expiresAt time.Time,
) Open {
	return Open{
		now:       now,
		expiresAt: expiresAt,
	}
}
