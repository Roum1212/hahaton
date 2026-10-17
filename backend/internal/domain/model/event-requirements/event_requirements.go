package model_event_requirements

import (
	"errors"
	"slices"

	model_event_requirement "github.com/Roum1212/hahaton/internal/domain/model/event-requirement"
)

type EventRequirements struct {
	values []model_event_requirement.Requirement
}

func (x EventRequirements) GetValues() []model_event_requirement.Requirement {
	return slices.Clone(x.values)
}

func (x EventRequirements) Validate() error {
	seen := make(map[model_event_requirement.Requirement]bool, len(x.values))
	hasAge := false

	for _, value := range x.values {
		if seen[value] {
			return errors.New("event requirements must not repeat")
		}

		if err := value.Validate(); err != nil {
			return err
		}

		seen[value] = true
		if value.IsAgeRequirement() {
			if hasAge {
				return errors.New("only one age requirement is allowed")
			}

			hasAge = true
		}
	}

	return nil
}

func NewEventRequirements(values []model_event_requirement.Requirement) (EventRequirements, error) {
	x := EventRequirements{values: slices.Clone(values)}
	if err := x.Validate(); err != nil {
		return EventRequirements{}, err
	}

	return x, nil
}
