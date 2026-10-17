package service_create_event

import (
	"fmt"

	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_event_requirements "github.com/Roum1212/hahaton/internal/domain/model/event-requirements"
	model_event_schedule "github.com/Roum1212/hahaton/internal/domain/model/event-schedule"
	model_event_title "github.com/Roum1212/hahaton/internal/domain/model/event-title"
	model_geo_point "github.com/Roum1212/hahaton/internal/domain/model/geo-point"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	model_volunteer_limit "github.com/Roum1212/hahaton/internal/domain/model/volunteer-limit"
)

type Request struct {
	userID         model_user_id.UserID
	title          model_event_title.EventTitle
	description    string
	imagePath      model_image_path.ImagePath
	location       string
	point          model_geo_point.GeoPoint
	schedule       model_event_schedule.EventSchedule
	requirements   model_event_requirements.EventRequirements
	volunteerLimit model_volunteer_limit.VolunteerLimit
}

func (x Request) GetDescription() string {
	return x.description
}

func (x Request) GetImagePath() model_image_path.ImagePath {
	return x.imagePath
}

func (x Request) GetLocation() string {
	return x.location
}

func (x Request) GetPoint() model_geo_point.GeoPoint {
	return x.point
}

func (x Request) GetRequirements() model_event_requirements.EventRequirements {
	return x.requirements
}

func (x Request) GetSchedule() model_event_schedule.EventSchedule {
	return x.schedule
}

func (x Request) GetTitle() model_event_title.EventTitle {
	return x.title
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) GetVolunteerLimit() model_volunteer_limit.VolunteerLimit {
	return x.volunteerLimit
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if err := x.title.Validate(); err != nil {
		return fmt.Errorf("title: %w", err)
	}

	if err := x.point.Validate(); err != nil {
		return fmt.Errorf("point: %w", err)
	}

	if err := x.schedule.Validate(); err != nil {
		return fmt.Errorf("schedule: %w", err)
	}

	if err := x.requirements.Validate(); err != nil {
		return fmt.Errorf("requirements: %w", err)
	}

	if err := x.volunteerLimit.Validate(); err != nil {
		return fmt.Errorf("volunteer limit: %w", err)
	}

	if err := x.imagePath.Validate(); err != nil {
		return fmt.Errorf("image path: %w", err)
	}

	return nil
}

func NewRequest(
	userID model_user_id.UserID,
	title model_event_title.EventTitle,
	description string,
	imagePath model_image_path.ImagePath,
	location string,
	point model_geo_point.GeoPoint,
	schedule model_event_schedule.EventSchedule,
	requirements model_event_requirements.EventRequirements,
	volunteerLimit model_volunteer_limit.VolunteerLimit,
) Request {
	return Request{
		userID:         userID,
		title:          title,
		description:    description,
		imagePath:      imagePath,
		location:       location,
		point:          point,
		schedule:       schedule,
		requirements:   requirements,
		volunteerLimit: volunteerLimit,
	}
}

type Response struct {
	eventID model_event_id.EventID
}

func (x Response) GetEventID() model_event_id.EventID {
	return x.eventID
}

func NewResponse(eventID model_event_id.EventID) Response {
	return Response{
		eventID: eventID,
	}
}
