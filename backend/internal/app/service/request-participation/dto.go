package service_request_participation

import (
	"fmt"

	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Request struct {
	userID  model_user_id.UserID
	eventID model_event_id.EventID
}

func (x Request) GetEventID() model_event_id.EventID {
	return x.eventID
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if err := x.eventID.Validate(); err != nil {
		return fmt.Errorf("event id: %w", err)
	}

	return nil
}

func NewRequest(userID model_user_id.UserID, eventID model_event_id.EventID) Request {
	return Request{
		userID:  userID,
		eventID: eventID,
	}
}

type Response struct {
	participationID model_participation_id.ParticipationID
}

func (x Response) GetParticipationID() model_participation_id.ParticipationID {
	return x.participationID
}

func NewResponse(participationID model_participation_id.ParticipationID) Response {
	return Response{
		participationID: participationID,
	}
}
