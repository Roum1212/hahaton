package aggregate_event_command

import (
	"time"
)

type Complete struct {
	requestID         string
	resultFingerprint string
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

func NewComplete(requestID, resultFingerprint string, now time.Time) Complete {
	return Complete{
		requestID:         requestID,
		resultFingerprint: resultFingerprint,
		now:               now,
	}
}
