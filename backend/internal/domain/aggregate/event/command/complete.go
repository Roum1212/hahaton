package aggregate_event_command

import (
	"time"
)

type Complete struct {
	requestID         string
	resultFingerprint string
	force             bool
	now               time.Time
}

func (x Complete) GetNow() time.Time {
	return x.now
}

func (x Complete) GetRequestID() string {
	return x.requestID
}

func (x Complete) GetResultFingerprint() string {
	return x.resultFingerprint
}

func (x Complete) IsForced() bool {
	return x.force
}

func NewComplete(requestID, resultFingerprint string, now time.Time) Complete {
	return Complete{
		requestID:         requestID,
		resultFingerprint: resultFingerprint,
		now:               now,
	}
}

func NewForceComplete(requestID, resultFingerprint string, now time.Time) Complete {
	return Complete{
		requestID:         requestID,
		resultFingerprint: resultFingerprint,
		force:             true,
		now:               now,
	}
}
