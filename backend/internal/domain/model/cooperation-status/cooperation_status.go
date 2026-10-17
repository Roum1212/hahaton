package model_cooperation_status

import (
	"errors"
)

type CooperationStatus string

const (
	Pending    CooperationStatus = "pending"
	Active     CooperationStatus = "active"
	Rejected   CooperationStatus = "rejected"
	Terminated CooperationStatus = "terminated"
)

func (x CooperationStatus) String() string {
	return string(x)
}

func (x CooperationStatus) Validate() error {
	switch x {
	case Pending:
		return nil
	case Active:
		return nil
	case Rejected:
		return nil
	case Terminated:
		return nil
	default:
		return errors.New("invalid cooperation status")
	}
}

func NewCooperationStatus(value string) (CooperationStatus, error) {
	x := CooperationStatus(value)
	if err := x.Validate(); err != nil {
		return CooperationStatus(""), err
	}

	return x, nil
}
