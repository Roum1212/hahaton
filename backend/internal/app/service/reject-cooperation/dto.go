package service_reject_cooperation

import (
	"fmt"

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
