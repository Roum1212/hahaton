package model_volunteer_limit

import (
	"errors"
)

type VolunteerLimit struct {
	value int32
}

func (x VolunteerLimit) GetValue() int32 {
	return x.value
}

func (x VolunteerLimit) Validate() error {
	if x.value <= 0 {
		return errors.New("volunteer limit must be positive")
	}

	return nil
}

func NewVolunteerLimit(value int32) (VolunteerLimit, error) {
	x := VolunteerLimit{value: value}
	if err := x.Validate(); err != nil {
		return VolunteerLimit{}, err
	}

	return x, nil
}
