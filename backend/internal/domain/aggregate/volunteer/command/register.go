package aggregate_volunteer_command

import (
	"time"

	model_date_of_birth "github.com/Roum1212/hahaton/internal/domain/model/date-of-birth"
	model_email_address "github.com/Roum1212/hahaton/internal/domain/model/email-address"
	model_full_name "github.com/Roum1212/hahaton/internal/domain/model/full-name"
	model_identity "github.com/Roum1212/hahaton/internal/domain/model/identity"
	model_phone_number "github.com/Roum1212/hahaton/internal/domain/model/phone-number"
)

type Register struct {
	identity     model_identity.Identity
	fullName     model_full_name.FullName
	dateOfBirth  model_date_of_birth.DateOfBirth
	emailAddress model_email_address.EmailAddress
	phoneNumber  model_phone_number.PhoneNumber
	now          time.Time
}

func (x Register) GetDateOfBirth() model_date_of_birth.DateOfBirth {
	return x.dateOfBirth
}

func (x Register) GetEmailAddress() model_email_address.EmailAddress {
	return x.emailAddress
}

func (x Register) GetFullName() model_full_name.FullName {
	return x.fullName
}

func (x Register) GetIdentity() model_identity.Identity {
	return x.identity
}

func (x Register) GetNow() time.Time {
	return x.now
}

func (x Register) GetPhoneNumber() model_phone_number.PhoneNumber {
	return x.phoneNumber
}

func NewRegister(
	identity model_identity.Identity,
	fullName model_full_name.FullName,
	dateOfBirth model_date_of_birth.DateOfBirth,
	emailAddress model_email_address.EmailAddress,
	phoneNumber model_phone_number.PhoneNumber,
	now time.Time,
) Register {
	return Register{
		identity:     identity,
		fullName:     fullName,
		dateOfBirth:  dateOfBirth,
		emailAddress: emailAddress,
		phoneNumber:  phoneNumber,
		now:          now,
	}
}
