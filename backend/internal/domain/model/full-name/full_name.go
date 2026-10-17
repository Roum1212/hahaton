package model_full_name

import (
	"errors"
	"fmt"
	"strings"
)

type FullName struct {
	firstName  string
	lastName   string
	middleName string
}

func (x FullName) GetFirstName() string {
	return x.firstName
}

func (x FullName) GetLastName() string {
	return x.lastName
}

func (x FullName) GetMiddleName() string {
	return x.middleName
}

func (x FullName) Validate() error {
	if len(x.firstName) == 0 {
		return errors.New("first name is empty")
	}

	if len(x.lastName) == 0 {
		return errors.New("last name is empty")
	}

	return nil
}

func NewFullName(firstName, lastName, middleName string) (FullName, error) {
	x := FullName{
		firstName:  strings.TrimSpace(firstName),
		lastName:   strings.TrimSpace(lastName),
		middleName: strings.TrimSpace(middleName),
	}
	if err := x.Validate(); err != nil {
		return FullName{}, fmt.Errorf("invalid full name: %w", err)
	}

	return x, nil
}
