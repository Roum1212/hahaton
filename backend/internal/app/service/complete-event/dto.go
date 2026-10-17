package service_complete_event

import (
	"errors"
	"fmt"
	"strings"

	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_event_result "github.com/Roum1212/hahaton/internal/domain/model/event-result"
	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Request struct {
	userID    model_user_id.UserID
	eventID   model_event_id.EventID
	results   []model_event_result.EventResult
	requestID string
	force     bool
}

func (x Request) GetEventID() model_event_id.EventID {
	return x.eventID
}

func (x Request) GetRequestID() string {
	return x.requestID
}

func (x Request) GetResults() []model_event_result.EventResult {
	return x.results
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) IsForced() bool {
	return x.force
}

func (x Request) Validate() error {
	if !x.force {
		if err := x.userID.Validate(); err != nil {
			return fmt.Errorf("user id: %w", err)
		}
	}

	if err := x.eventID.Validate(); err != nil {
		return fmt.Errorf("event id: %w", err)
	}

	if len(strings.TrimSpace(x.requestID)) == 0 {
		return errors.New("request id is required")
	}

	seen := make(map[model_participation_id.ParticipationID]bool, len(x.results))
	for _, result := range x.results {
		if err := result.Validate(); err != nil {
			return err
		}

		if seen[result.GetParticipationID()] {
			return errors.New("duplicate participation result")
		}

		seen[result.GetParticipationID()] = true
	}

	return nil
}

func NewForceRequest(
	eventID model_event_id.EventID,
	results []model_event_result.EventResult,
	requestID string,
) Request {
	return Request{
		eventID:   eventID,
		results:   results,
		requestID: requestID,
		force:     true,
	}
}

func NewRequest(
	userID model_user_id.UserID,
	eventID model_event_id.EventID,
	results []model_event_result.EventResult,
	requestID string,
) Request {
	return Request{
		userID:    userID,
		eventID:   eventID,
		results:   results,
		requestID: requestID,
	}
}

type Response struct {
	event *aggregate_event.Event
}

func (x Response) GetEvent() *aggregate_event.Event {
	return x.event
}

func NewResponse(event *aggregate_event.Event) Response {
	return Response{
		event: event,
	}
}
