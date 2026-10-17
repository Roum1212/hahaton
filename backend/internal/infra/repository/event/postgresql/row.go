package repository_event_postgresql

import (
	"fmt"
	"time"

	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_event_requirement "github.com/Roum1212/hahaton/internal/domain/model/event-requirement"
	model_event_requirements "github.com/Roum1212/hahaton/internal/domain/model/event-requirements"
	model_event_schedule "github.com/Roum1212/hahaton/internal/domain/model/event-schedule"
	model_event_status "github.com/Roum1212/hahaton/internal/domain/model/event-status"
	model_event_title "github.com/Roum1212/hahaton/internal/domain/model/event-title"
	model_geo_point "github.com/Roum1212/hahaton/internal/domain/model/geo-point"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	model_volunteer_limit "github.com/Roum1212/hahaton/internal/domain/model/volunteer-limit"
)

type row struct {
	eventID             string
	coordinatorID       string
	title               string
	description         string
	imagePath           string
	location            string
	latitude            float64
	longitude           float64
	startsAt            time.Time
	endsAt              time.Time
	requirements        []int32
	volunteerLimit      int32
	status              string
	completionRequestID string
	resultFingerprint   string
	createdAt           time.Time
	updatedAt           time.Time
}

func (x *row) scan(source interface{ Scan(dest ...any) error }) error {
	return source.Scan(
		&x.eventID,
		&x.coordinatorID,
		&x.title,
		&x.description,
		&x.imagePath,
		&x.location,
		&x.latitude,
		&x.longitude,
		&x.startsAt,
		&x.endsAt,
		&x.requirements,
		&x.volunteerLimit,
		&x.status,
		&x.completionRequestID,
		&x.resultFingerprint,
		&x.createdAt,
		&x.updatedAt,
	)
}

func (x *row) toAggregate() (*aggregate_event.Event, error) {
	eventID, err := model_event_id.NewEventID(x.eventID)
	if err != nil {
		return nil, fmt.Errorf("restore event id: %w", err)
	}

	coordinatorID, err := model_user_id.NewUserID(x.coordinatorID)
	if err != nil {
		return nil, fmt.Errorf("restore coordinator id: %w", err)
	}

	title, err := model_event_title.NewEventTitle(x.title)
	if err != nil {
		return nil, fmt.Errorf("restore title: %w", err)
	}

	var imagePath model_image_path.ImagePath
	if len(x.imagePath) != 0 {
		value, err := model_image_path.NewImagePath(x.imagePath)
		if err != nil {
			return nil, fmt.Errorf("restore value: %w", err)
		}

		imagePath = value
	}

	point, err := model_geo_point.NewGeoPoint(x.latitude, x.longitude)
	if err != nil {
		return nil, fmt.Errorf("restore point: %w", err)
	}

	schedule, err := model_event_schedule.NewEventSchedule(x.startsAt, x.endsAt)
	if err != nil {
		return nil, fmt.Errorf("restore schedule: %w", err)
	}

	values := make([]model_event_requirement.Requirement, len(x.requirements))
	for index, value := range x.requirements {
		requirement, err := model_event_requirement.NewRequirement(value)
		if err != nil {
			return nil, fmt.Errorf("restore requirement: %w", err)
		}

		values[index] = requirement
	}

	requirements, err := model_event_requirements.NewEventRequirements(values)
	if err != nil {
		return nil, fmt.Errorf("restore requirements: %w", err)
	}

	volunteerLimit, err := model_volunteer_limit.NewVolunteerLimit(x.volunteerLimit)
	if err != nil {
		return nil, fmt.Errorf("restore volunteer limit: %w", err)
	}

	status, err := model_event_status.NewEventStatus(x.status)
	if err != nil {
		return nil, fmt.Errorf("restore status: %w", err)
	}

	return aggregate_event.NewEventFromFields(
		eventID,
		coordinatorID,
		title,
		x.description,
		imagePath,
		x.location,
		point,
		schedule,
		requirements,
		volunteerLimit,
		status,
		x.completionRequestID,
		x.resultFingerprint,
		x.createdAt,
		x.updatedAt,
	), nil
}
