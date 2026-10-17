package model_event_status

import (
	"errors"
)

type EventStatus string

const (
	Published EventStatus = "published"
	Canceled  EventStatus = "canceled"
	Completed EventStatus = "completed"
)

func (x EventStatus) String() string {
	return string(x)
}

func (x EventStatus) Validate() error {
	switch x {
	case Published:
		return nil
	case Canceled:
		return nil
	case Completed:
		return nil
	default:
		return errors.New("invalid event status")
	}
}

func NewEventStatus(value string) (EventStatus, error) {
	x := EventStatus(value)
	if err := x.Validate(); err != nil {
		return EventStatus(""), err
	}

	return x, nil
}
