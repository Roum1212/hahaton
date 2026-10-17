package service_get_participation

import (
	"fmt"

	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
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

type Response struct {
	participation *aggregate_participation.Participation
}

func (x Response) GetParticipation() *aggregate_participation.Participation {
	return x.participation
}

func NewResponse(participation *aggregate_participation.Participation) Response {
	return Response{
		participation: participation,
	}
}
