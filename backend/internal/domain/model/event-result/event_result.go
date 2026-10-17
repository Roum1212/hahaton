package model_event_result

import (
	"errors"
	"math"

	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
)

const (
	MinutesPerStep = 30
	CoinsPerStep   = 5
)

type EventResult struct {
	participationID model_participation_id.ParticipationID
	minutes         int64
	absent          bool
}

func (x EventResult) GetMinutes() int64 {
	return x.minutes
}

func (x EventResult) GetParticipationID() model_participation_id.ParticipationID {
	return x.participationID
}

func (x EventResult) IsAbsent() bool {
	return x.absent
}

func (x EventResult) Validate() error {
	if err := x.participationID.Validate(); err != nil {
		return err
	}

	if x.absent {
		if x.minutes != 0 {
			return errors.New("absent participant cannot have confirmed minutes")
		}

		return nil
	}

	if x.minutes <= 0 || x.minutes%MinutesPerStep != 0 || x.minutes/MinutesPerStep > math.MaxInt64/CoinsPerStep {
		return errors.New("confirmed minutes must be a positive multiple of 30 without coin overflow")
	}

	return nil
}

func NewEventResult(
	participationID model_participation_id.ParticipationID,
	minutes int64,
	absent bool,
) (EventResult, error) {
	x := EventResult{
		participationID: participationID,
		minutes:         minutes,
		absent:          absent,
	}
	if err := x.Validate(); err != nil {
		return EventResult{}, err
	}

	return x, nil
}
