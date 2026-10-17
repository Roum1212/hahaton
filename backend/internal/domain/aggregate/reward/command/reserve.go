package aggregate_reward_command

import (
	"time"
)

type Reserve struct {
	expectedPrice int64
	now           time.Time
}

func (x Reserve) GetExpectedPrice() int64 {
	return x.expectedPrice
}

func (x Reserve) GetNow() time.Time {
	return x.now
}

func NewReserve(expectedPrice int64, now time.Time) Reserve {
	return Reserve{
		expectedPrice: expectedPrice,
		now:           now,
	}
}
