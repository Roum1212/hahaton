package model_event_title

import (
	"errors"
	"strings"
)

type EventTitle struct {
	value string
}

func (x EventTitle) String() string {
	return x.value
}

func (x EventTitle) Validate() error {
	if len(strings.TrimSpace(x.value)) == 0 {
		return errors.New("event title is required")
	}

	return nil
}

func NewEventTitle(value string) (EventTitle, error) {
	x := EventTitle{value: strings.TrimSpace(value)}
	if err := x.Validate(); err != nil {
		return EventTitle{}, err
	}

	return x, nil
}
