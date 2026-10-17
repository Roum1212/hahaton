package aggregate_volunteer

import (
	"errors"
	"time"

	aggregate_volunteer_command "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer/command"
	model_date_of_birth "github.com/Roum1212/hahaton/internal/domain/model/date-of-birth"
	model_email_address "github.com/Roum1212/hahaton/internal/domain/model/email-address"
	model_full_name "github.com/Roum1212/hahaton/internal/domain/model/full-name"
	model_identity "github.com/Roum1212/hahaton/internal/domain/model/identity"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	model_phone_number "github.com/Roum1212/hahaton/internal/domain/model/phone-number"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

var ErrAlreadyRegistered = errors.New("volunteer is already registered")

type Volunteer struct {
	userID       model_user_id.UserID
	identity     model_identity.Identity
	fullName     model_full_name.FullName
	dateOfBirth  model_date_of_birth.DateOfBirth
	emailAddress model_email_address.EmailAddress
	phoneNumber  model_phone_number.PhoneNumber
	avatarPath   model_image_path.ImagePath
	createdAt    time.Time
	updatedAt    time.Time
}

func (x *Volunteer) GetAvatarPath() model_image_path.ImagePath {
	return x.avatarPath
}

func (x *Volunteer) GetCreatedAt() time.Time {
	return x.createdAt
}

func (x *Volunteer) GetDateOfBirth() model_date_of_birth.DateOfBirth {
	return x.dateOfBirth
}

func (x *Volunteer) GetEmailAddress() model_email_address.EmailAddress {
	return x.emailAddress
}

func (x *Volunteer) GetFullName() model_full_name.FullName {
	return x.fullName
}

func (x *Volunteer) GetIdentity() model_identity.Identity {
	return x.identity
}

func (x *Volunteer) GetPhoneNumber() model_phone_number.PhoneNumber {
	return x.phoneNumber
}

func (x *Volunteer) GetUpdatedAt() time.Time {
	return x.updatedAt
}

func (x *Volunteer) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x *Volunteer) Register(command aggregate_volunteer_command.Register) error {
	if !x.createdAt.IsZero() {
		return ErrAlreadyRegistered
	}

	x.identity = command.GetIdentity()
	x.fullName = command.GetFullName()
	x.dateOfBirth = command.GetDateOfBirth()
	x.emailAddress = command.GetEmailAddress()
	x.phoneNumber = command.GetPhoneNumber()
	x.createdAt = command.GetNow()
	x.updatedAt = command.GetNow()

	return nil
}

func (x *Volunteer) Update(command aggregate_volunteer_command.Update) {
	x.fullName = command.GetFullName()
	x.dateOfBirth = command.GetDateOfBirth()
	x.avatarPath = command.GetAvatarPath()
	x.updatedAt = command.GetNow()
}

func NewVolunteer(userID model_user_id.UserID) *Volunteer {
	return &Volunteer{userID: userID}
}

func NewVolunteerFromFields(
	userID model_user_id.UserID,
	identity model_identity.Identity,
	fullName model_full_name.FullName,
	dateOfBirth model_date_of_birth.DateOfBirth,
	emailAddress model_email_address.EmailAddress,
	phoneNumber model_phone_number.PhoneNumber,
	avatarPath model_image_path.ImagePath,
	createdAt time.Time,
	updatedAt time.Time,
) *Volunteer {
	return &Volunteer{
		userID:       userID,
		identity:     identity,
		fullName:     fullName,
		dateOfBirth:  dateOfBirth,
		emailAddress: emailAddress,
		phoneNumber:  phoneNumber,
		avatarPath:   avatarPath,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}
