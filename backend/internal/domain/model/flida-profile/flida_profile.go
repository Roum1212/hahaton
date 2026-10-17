package model_flida_profile

import (
	"errors"
	"fmt"

	model_date_of_birth "github.com/Roum1212/hahaton/internal/domain/model/date-of-birth"
	model_email_address "github.com/Roum1212/hahaton/internal/domain/model/email-address"
	model_full_name "github.com/Roum1212/hahaton/internal/domain/model/full-name"
	model_identity "github.com/Roum1212/hahaton/internal/domain/model/identity"
	model_phone_number "github.com/Roum1212/hahaton/internal/domain/model/phone-number"
)

type Profile struct {
	identity     model_identity.Identity
	fullName     model_full_name.FullName
	dateOfBirth  model_date_of_birth.DateOfBirth
	emailAddress model_email_address.EmailAddress
	phoneNumber  model_phone_number.PhoneNumber
}

func (x Profile) GetDateOfBirth() model_date_of_birth.DateOfBirth {
	return x.dateOfBirth
}

func (x Profile) GetEmailAddress() model_email_address.EmailAddress {
	return x.emailAddress
}

func (x Profile) GetFullName() model_full_name.FullName {
	return x.fullName
}

func (x Profile) GetIdentity() model_identity.Identity {
	return x.identity
}

func (x Profile) GetPhoneNumber() model_phone_number.PhoneNumber {
	return x.phoneNumber
}

func (x Profile) Validate() error {
	if err := x.identity.Validate(); err != nil {
		return err
	}

	if err := x.fullName.Validate(); err != nil {
		return err
	}

	if x.dateOfBirth.IsZero() {
		return errors.New("date of birth is required")
	}

	if err := x.dateOfBirth.Validate(); err != nil {
		return err
	}

	if err := x.emailAddress.Validate(); err != nil {
		return err
	}

	return x.phoneNumber.Validate()
}

func NewProfile(
	identity model_identity.Identity,
	fullName model_full_name.FullName,
	dateOfBirth model_date_of_birth.DateOfBirth,
	emailAddress model_email_address.EmailAddress,
	phoneNumber model_phone_number.PhoneNumber,
) (Profile, error) {
	x := Profile{
		identity:     identity,
		fullName:     fullName,
		dateOfBirth:  dateOfBirth,
		emailAddress: emailAddress,
		phoneNumber:  phoneNumber,
	}
	if err := x.Validate(); err != nil {
		return Profile{}, fmt.Errorf("invalid flida profile: %w", err)
	}

	return x, nil
}
