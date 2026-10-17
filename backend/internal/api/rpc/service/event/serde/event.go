package rpc_service_event_serde

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	event_v1 "github.com/Roum1212/hahaton/genproto/event/v1"
	rpc_service_geo_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/geo/serde"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	model_event_requirement "github.com/Roum1212/hahaton/internal/domain/model/event-requirement"
	model_event_requirements "github.com/Roum1212/hahaton/internal/domain/model/event-requirements"
	model_event_result "github.com/Roum1212/hahaton/internal/domain/model/event-result"
	model_event_status "github.com/Roum1212/hahaton/internal/domain/model/event-status"
	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
)

func AggregateToPBEvent(event *aggregate_event.Event, approvedCount, pendingCount int32) *event_v1.Event {
	requirements := event.GetRequirements().GetValues()

	values := make([]event_v1.EventCandidateRequirement, len(requirements))
	for index, value := range requirements {
		values[index] = RequirementToPB(value)
	}

	return &event_v1.Event{
		Id:                       event.GetEventID().String(),
		CoordinatorId:            event.GetCoordinatorID().String(),
		Description:              event.GetDescription(),
		ImageUrl:                 event.GetImagePath().GetURL(),
		Title:                    event.GetTitle().String(),
		Location:                 event.GetLocation(),
		Point:                    rpc_service_geo_serde.PointToPB(event.GetPoint()),
		StartsAt:                 timestamppb.New(event.GetSchedule().GetStartsAt()),
		EndsAt:                   timestamppb.New(event.GetSchedule().GetEndsAt()),
		CandidateRequirements:    values,
		VolunteerLimit:           event.GetVolunteerLimit().GetValue(),
		Status:                   StatusToPB(event.GetStatus()),
		ApprovedParticipantCount: approvedCount,
		PendingApplicationCount:  pendingCount,
	}
}

func RequirementFromPB(value event_v1.EventCandidateRequirement) (model_event_requirement.Requirement, error) {
	switch value {
	case event_v1.EventCandidateRequirement_EVENT_CANDIDATE_REQUIREMENT_AGE_14_PLUS:
		return model_event_requirement.NewRequirement(model_event_requirement.Age14Plus.Int32())
	case event_v1.EventCandidateRequirement_EVENT_CANDIDATE_REQUIREMENT_AGE_16_PLUS:
		return model_event_requirement.NewRequirement(model_event_requirement.Age16Plus.Int32())
	case event_v1.EventCandidateRequirement_EVENT_CANDIDATE_REQUIREMENT_AGE_18_PLUS:
		return model_event_requirement.NewRequirement(model_event_requirement.Age18Plus.Int32())
	case event_v1.EventCandidateRequirement_EVENT_CANDIDATE_REQUIREMENT_PRELIMINARY_TRAINING:
		return model_event_requirement.NewRequirement(model_event_requirement.PreliminaryTraining.Int32())
	case event_v1.EventCandidateRequirement_EVENT_CANDIDATE_REQUIREMENT_RUSSIAN_CITIZENSHIP:
		return model_event_requirement.NewRequirement(model_event_requirement.RussianCitizenship.Int32())
	case event_v1.EventCandidateRequirement_EVENT_CANDIDATE_REQUIREMENT_VOLUNTEERING_EXPERIENCE:
		return model_event_requirement.NewRequirement(model_event_requirement.VolunteeringExperience.Int32())
	case event_v1.EventCandidateRequirement_EVENT_CANDIDATE_REQUIREMENT_UNSPECIFIED:
		return 0, errors.New("unspecified event requirement")
	default:
		return 0, fmt.Errorf("unknown event requirement: %d", value)
	}
}

func RequirementToPB(value model_event_requirement.Requirement) event_v1.EventCandidateRequirement {
	switch value {
	case model_event_requirement.Age14Plus:
		return event_v1.EventCandidateRequirement_EVENT_CANDIDATE_REQUIREMENT_AGE_14_PLUS
	case model_event_requirement.Age16Plus:
		return event_v1.EventCandidateRequirement_EVENT_CANDIDATE_REQUIREMENT_AGE_16_PLUS
	case model_event_requirement.Age18Plus:
		return event_v1.EventCandidateRequirement_EVENT_CANDIDATE_REQUIREMENT_AGE_18_PLUS
	case model_event_requirement.PreliminaryTraining:
		return event_v1.EventCandidateRequirement_EVENT_CANDIDATE_REQUIREMENT_PRELIMINARY_TRAINING
	case model_event_requirement.RussianCitizenship:
		return event_v1.EventCandidateRequirement_EVENT_CANDIDATE_REQUIREMENT_RUSSIAN_CITIZENSHIP
	case model_event_requirement.VolunteeringExperience:
		return event_v1.EventCandidateRequirement_EVENT_CANDIDATE_REQUIREMENT_VOLUNTEERING_EXPERIENCE
	default:
		return event_v1.EventCandidateRequirement_EVENT_CANDIDATE_REQUIREMENT_UNSPECIFIED
	}
}

func RequirementsFromPB(
	values []event_v1.EventCandidateRequirement,
) (model_event_requirements.EventRequirements, error) {
	requirements := make([]model_event_requirement.Requirement, len(values))
	for index, value := range values {
		requirement, err := RequirementFromPB(value)
		if err != nil {
			return model_event_requirements.EventRequirements{}, err
		}

		requirements[index] = requirement
	}

	return model_event_requirements.NewEventRequirements(requirements)
}

func ResultFromPB(value *event_v1.EventParticipationResult) (model_event_result.EventResult, error) {
	if value == nil {
		return model_event_result.EventResult{}, errors.New("participation result is required")
	}

	id, err := model_participation_id.NewParticipationID(value.GetParticipationId())
	if err != nil {
		return model_event_result.EventResult{}, fmt.Errorf("participation id: %w", err)
	}

	switch outcome := value.GetOutcome().(type) {
	case *event_v1.EventParticipationResult_Absent:
		if !outcome.Absent {
			return model_event_result.EventResult{}, errors.New("absent must be true")
		}

		return model_event_result.NewEventResult(id, 0, true)
	case *event_v1.EventParticipationResult_ConfirmedMinutes:
		return model_event_result.NewEventResult(id, outcome.ConfirmedMinutes, false)
	default:
		return model_event_result.EventResult{}, errors.New("one participation outcome is required")
	}
}

func StatusFromPB(value event_v1.EventStatus) (model_event_status.EventStatus, error) {
	switch value {
	case event_v1.EventStatus_EVENT_STATUS_PUBLISHED:
		return model_event_status.NewEventStatus(model_event_status.Published.String())
	case event_v1.EventStatus_EVENT_STATUS_CANCELLED:
		return model_event_status.NewEventStatus(model_event_status.Canceled.String())
	case event_v1.EventStatus_EVENT_STATUS_COMPLETED:
		return model_event_status.NewEventStatus(model_event_status.Completed.String())
	case event_v1.EventStatus_EVENT_STATUS_UNSPECIFIED:
		return "", errors.New("unspecified event status")
	default:
		return "", fmt.Errorf("unknown event status: %d", value)
	}
}

func StatusToPB(value model_event_status.EventStatus) event_v1.EventStatus {
	statuses := map[model_event_status.EventStatus]event_v1.EventStatus{
		model_event_status.Published: event_v1.EventStatus_EVENT_STATUS_PUBLISHED,
		model_event_status.Canceled:  event_v1.EventStatus_EVENT_STATUS_CANCELLED,
		model_event_status.Completed: event_v1.EventStatus_EVENT_STATUS_COMPLETED,
	}

	return statuses[value]
}
