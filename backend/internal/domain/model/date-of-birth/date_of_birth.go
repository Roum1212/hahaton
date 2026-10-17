package model_date_of_birth

import (
	"errors"
	"fmt"
	"time"
)

const minYear = 1900

type DateOfBirth struct {
	year  int
	month time.Month
	day   int
}

func (x DateOfBirth) GetDay() int {
	return x.day
}

func (x DateOfBirth) GetMonth() time.Month {
	return x.month
}

func (x DateOfBirth) GetYear() int {
	return x.year
}

func (x DateOfBirth) IsZero() bool {
	return x.isZero()
}

func (x DateOfBirth) String() string {
	if x.isZero() {
		return ""
	}

	return x.Time().Format(time.DateOnly)
}

func (x DateOfBirth) Time() time.Time {
	if x.isZero() {
		return time.Time{}
	}

	return x.time()
}

func (x DateOfBirth) Validate() error {
	if x.IsZero() {
		return nil
	}

	birthdate := x.time()

	if birthdate.Year() != x.year || birthdate.Month() != x.month || birthdate.Day() != x.day {
		return errors.New("invalid date")
	}

	if birthdate.Year() < minYear {
		return errors.New("birthdate too far in the past")
	}

	if birthdate.After(time.Now()) {
		return errors.New("birthdate can't be in the future")
	}

	return nil
}

func (x DateOfBirth) isZero() bool {
	return x.year == 0 && x.month == 0 && x.day == 0
}

func (x DateOfBirth) time() time.Time {
	return time.Date(x.year, x.month, x.day, 0, 0, 0, 0, time.Local)
}

func NewDateOfBirth(year int, month time.Month, day int) (DateOfBirth, error) {
	x := DateOfBirth{
		year:  year,
		month: month,
		day:   day,
	}
	if err := x.Validate(); err != nil {
		return DateOfBirth{}, fmt.Errorf("invalid date of birth: %w", err)
	}

	return x, nil
}

func NewDateOfBirthFromString(s string) (DateOfBirth, error) {
	if s == "" {
		return NewDateOfBirth(0, 0, 0)
	}

	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return DateOfBirth{}, fmt.Errorf("failed to parse date of birth: %w", err)
	}

	return NewDateOfBirth(t.Year(), t.Month(), t.Day())
}
