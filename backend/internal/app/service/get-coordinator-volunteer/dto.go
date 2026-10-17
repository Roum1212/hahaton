package service_get_coordinator_volunteer

import (
	"fmt"

	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Request struct {
	userID      model_user_id.UserID
	volunteerID model_user_id.UserID
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) GetVolunteerID() model_user_id.UserID {
	return x.volunteerID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if err := x.volunteerID.Validate(); err != nil {
		return fmt.Errorf("volunteer id: %w", err)
	}

	return nil
}

func NewRequest(userID, volunteerID model_user_id.UserID) Request {
	return Request{
		userID:      userID,
		volunteerID: volunteerID,
	}
}

type Response struct {
	volunteer    *aggregate_volunteer.Volunteer
	showContacts bool
}

func (x Response) GetShowContacts() bool {
	return x.showContacts
}

func (x Response) GetVolunteer() *aggregate_volunteer.Volunteer {
	return x.volunteer
}

func NewResponse(volunteer *aggregate_volunteer.Volunteer, showContacts bool) Response {
	return Response{
		volunteer:    volunteer,
		showContacts: showContacts,
	}
}
