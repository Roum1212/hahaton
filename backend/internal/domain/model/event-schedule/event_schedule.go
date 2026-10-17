package model_event_schedule

import (
	"errors"
	"time"
)

const DurationStep = 30 * time.Minute

type EventSchedule struct {
	startsAt time.Time
	endsAt   time.Time
}

func (x EventSchedule) GetEndsAt() time.Time {
	return x.endsAt
}

func (x EventSchedule) GetStartsAt() time.Time {
	return x.startsAt
}

func (x EventSchedule) Validate() error {
	switch {
	case x.startsAt.IsZero(), x.endsAt.IsZero():
		return errors.New("event dates are required")
	case !x.endsAt.After(x.startsAt):
		return errors.New("event must end after its start")
	case x.endsAt.Sub(x.startsAt)%DurationStep != 0:
		return errors.New("event duration must be a multiple of 30 minutes")
	}

	return nil
}

func NewEventSchedule(startsAt, endsAt time.Time) (EventSchedule, error) {
	x := EventSchedule{
		startsAt: startsAt.UTC(),
		endsAt:   endsAt.UTC(),
	}
	if err := x.Validate(); err != nil {
		return EventSchedule{}, err
	}

	return x, nil
}
