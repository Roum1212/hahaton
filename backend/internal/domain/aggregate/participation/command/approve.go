package aggregate_participation_command

import (
	"time"
)

type Approve struct {
	visibleAfterSequence uint64
	now                  time.Time
}

func (x Approve) GetNow() time.Time {
	return x.now
}

func (x Approve) GetVisibleAfterSequence() uint64 {
	return x.visibleAfterSequence
}

func NewApprove(visibleAfterSequence uint64, now time.Time) Approve {
	return Approve{
		visibleAfterSequence: visibleAfterSequence,
		now:                  now,
	}
}
