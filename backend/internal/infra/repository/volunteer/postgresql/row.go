package repository_volunteer_postgresql

import (
	"fmt"
	"time"

	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
	model_date_of_birth "github.com/Roum1212/hahaton/internal/domain/model/date-of-birth"
	model_email_address "github.com/Roum1212/hahaton/internal/domain/model/email-address"
	model_full_name "github.com/Roum1212/hahaton/internal/domain/model/full-name"
	model_identity "github.com/Roum1212/hahaton/internal/domain/model/identity"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	model_phone_number "github.com/Roum1212/hahaton/internal/domain/model/phone-number"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type volunteerRow struct {
	id           string
	firstName    string
	lastName     string
	middleName   string
	dateOfBirth  time.Time
	emailAddress string
	phoneNumber  string
	avatarPath   string
	createdAt    time.Time
	updatedAt    time.Time
}

func (x volunteerRow) toVolunteer(identity model_identity.Identity) (*aggregate_volunteer.Volunteer, error) {
	userID, err := model_user_id.NewUserID(x.id)
	if err != nil {
		return nil, fmt.Errorf("restore user id: %w", err)
	}

	fullName, err := model_full_name.NewFullName(x.firstName, x.lastName, x.middleName)
	if err != nil {
		return nil, fmt.Errorf("restore full name: %w", err)
	}

	dateOfBirth, err := model_date_of_birth.NewDateOfBirth(
		x.dateOfBirth.Year(),
		x.dateOfBirth.Month(),
		x.dateOfBirth.Day(),
	)
	if err != nil {
		return nil, fmt.Errorf("restore date of birth: %w", err)
	}

	emailAddress, err := model_email_address.NewEmailAddress(x.emailAddress)
	if err != nil {
		return nil, fmt.Errorf("restore email address: %w", err)
	}

	phoneNumber, err := model_phone_number.NewPhoneNumber(x.phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("restore phone number: %w", err)
	}

	var avatarPath model_image_path.ImagePath

	if len(x.avatarPath) != 0 {
		avatarPath, err = model_image_path.NewImagePath(x.avatarPath)
		if err != nil {
			return nil, fmt.Errorf("restore avatar path: %w", err)
		}
	}

	return aggregate_volunteer.NewVolunteerFromFields(
		userID, identity, fullName, dateOfBirth, emailAddress, phoneNumber,
		avatarPath, x.createdAt, x.updatedAt,
	), nil
}
