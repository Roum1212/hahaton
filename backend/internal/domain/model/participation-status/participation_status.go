package model_participation_status

import (
	"errors"
)

type ParticipationStatus string

const (
	Pending     ParticipationStatus = "pending"
	Approved    ParticipationStatus = "approved"
	Rejected    ParticipationStatus = "rejected"
	Canceled    ParticipationStatus = "canceled"
	Confirmed   ParticipationStatus = "confirmed"
	NotAttended ParticipationStatus = "notattended"
	Closed      ParticipationStatus = "closed"
)

func (x ParticipationStatus) String() string {
	return string(x)
}

func (x ParticipationStatus) Validate() error {
	switch x {
	case Pending:
		return nil
	case Approved:
		return nil
	case Rejected:
		return nil
	case Canceled:
		return nil
	case Confirmed:
		return nil
	case NotAttended:
		return nil
	case Closed:
		return nil
	default:
		return errors.New("invalid participation status")
	}
}

func NewParticipationStatus(value string) (ParticipationStatus, error) {
	x := ParticipationStatus(value)
	if err := x.Validate(); err != nil {
		return ParticipationStatus(""), err
	}

	return x, nil
}
