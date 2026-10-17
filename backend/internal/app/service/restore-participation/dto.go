package service_restore_participation

import (
	"fmt"

	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Request struct {
	userID          model_user_id.UserID
	participationID model_participation_id.ParticipationID
}

func (x Request) GetParticipationID() model_participation_id.ParticipationID {
	return x.participationID
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if err := x.participationID.Validate(); err != nil {
		return fmt.Errorf("participation id: %w", err)
	}

	return nil
}

func NewRequest(userID model_user_id.UserID, participationID model_participation_id.ParticipationID) Request {
	return Request{
		userID:          userID,
		participationID: participationID,
	}
}
