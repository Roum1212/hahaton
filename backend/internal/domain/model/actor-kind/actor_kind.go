package model_actor_kind

import (
	"errors"
)

type ActorKind string

const (
	Coordinator ActorKind = "coordinator"
	Volunteer   ActorKind = "volunteer"
	System      ActorKind = "system"
)

func (x ActorKind) String() string {
	return string(x)
}

func (x ActorKind) Validate() error {
	switch x {
	case Coordinator:
		return nil
	case Volunteer:
		return nil
	case System:
		return nil
	default:
		return errors.New("invalid actor kind")
	}
}

func NewActorKind(value string) (ActorKind, error) {
	x := ActorKind(value)
	if err := x.Validate(); err != nil {
		return ActorKind(""), err
	}

	return x, nil
}
