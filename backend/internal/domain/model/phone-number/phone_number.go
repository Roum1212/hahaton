package model_phone_number

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var pattern = regexp.MustCompile(`^\+[1-9][0-9]{1,14}$`)

type PhoneNumber struct {
	value string
}

func (x PhoneNumber) String() string {
	return x.value
}

func (x PhoneNumber) Validate() error {
	if !pattern.MatchString(x.value) {
		return errors.New("phone number must use e.164 format")
	}

	return nil
}

func NewPhoneNumber(value string) (PhoneNumber, error) {
	x := PhoneNumber{value: strings.TrimSpace(value)}
	if err := x.Validate(); err != nil {
		return PhoneNumber{}, fmt.Errorf("invalid phone number: %w", err)
	}

	return x, nil
}
