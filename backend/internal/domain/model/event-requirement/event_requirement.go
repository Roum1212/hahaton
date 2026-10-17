package model_event_requirement

import (
	"errors"
)

type Requirement int32

const (
	Age14Plus Requirement = iota + 1
	Age16Plus
	Age18Plus
	PreliminaryTraining
	RussianCitizenship
	VolunteeringExperience
)

func (x Requirement) Int32() int32 {
	return int32(x)
}

func (x Requirement) IsAgeRequirement() bool {
	return x == Age14Plus || x == Age16Plus || x == Age18Plus
}

func (x Requirement) Validate() error {
	switch x {
	case Age14Plus:
		return nil
	case Age16Plus:
		return nil
	case Age18Plus:
		return nil
	case PreliminaryTraining:
		return nil
	case RussianCitizenship:
		return nil
	case VolunteeringExperience:
		return nil
	default:
		return errors.New("unknown event requirement")
	}
}

func NewRequirement(value int32) (Requirement, error) {
	x := Requirement(value)
	if err := x.Validate(); err != nil {
		return 0, err
	}

	return x, nil
}
