package aggregate_reward_order_command

import (
	"time"
)

type Create struct {
	rewardTitle string
	priceCoins  int64
	requestID   string
	now         time.Time
}

func (x Create) GetNow() time.Time {
	return x.now
}

func (x Create) GetPriceCoins() int64 {
	return x.priceCoins
}

func (x Create) GetRequestID() string {
	return x.requestID
}

func (x Create) GetRewardTitle() string {
	return x.rewardTitle
}

func NewCreate(rewardTitle string, priceCoins int64, requestID string, now time.Time) Create {
	return Create{
		rewardTitle: rewardTitle,
		priceCoins:  priceCoins,
		requestID:   requestID,
		now:         now,
	}
}
