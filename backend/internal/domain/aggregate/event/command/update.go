package aggregate_event_command

import (
	"time"

	model_event_requirements "github.com/Roum1212/hahaton/internal/domain/model/event-requirements"
	model_event_schedule "github.com/Roum1212/hahaton/internal/domain/model/event-schedule"
	model_event_title "github.com/Roum1212/hahaton/internal/domain/model/event-title"
	model_geo_point "github.com/Roum1212/hahaton/internal/domain/model/geo-point"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	model_volunteer_limit "github.com/Roum1212/hahaton/internal/domain/model/volunteer-limit"
)

type Update struct {
	title                   model_event_title.EventTitle
	description             string
	imagePath               model_image_path.ImagePath
	location                string
	point                   model_geo_point.GeoPoint
	schedule                model_event_schedule.EventSchedule
	requirements            model_event_requirements.EventRequirements
	volunteerLimit          model_volunteer_limit.VolunteerLimit
	now                     time.Time
	hasActiveParticipations bool
}

func (x Update) GetDescription() string {
	return x.description
}

func (x Update) GetHasActiveParticipations() bool {
	return x.hasActiveParticipations
}

func (x Update) GetImagePath() model_image_path.ImagePath {
	return x.imagePath
}

func (x Update) GetLocation() string {
	return x.location
}

func (x Update) GetNow() time.Time {
	return x.now
}

func (x Update) GetPoint() model_geo_point.GeoPoint {
	return x.point
}

func (x Update) GetRequirements() model_event_requirements.EventRequirements {
	return x.requirements
}

func (x Update) GetSchedule() model_event_schedule.EventSchedule {
	return x.schedule
}

func (x Update) GetTitle() model_event_title.EventTitle {
	return x.title
}

func (x Update) GetVolunteerLimit() model_volunteer_limit.VolunteerLimit {
	return x.volunteerLimit
}

func NewUpdate(
	title model_event_title.EventTitle,
	description string,
	imagePath model_image_path.ImagePath,
	location string,
	point model_geo_point.GeoPoint,
	schedule model_event_schedule.EventSchedule,
	requirements model_event_requirements.EventRequirements,
	volunteerLimit model_volunteer_limit.VolunteerLimit,
	now time.Time,
	hasActiveParticipations bool,
) Update {
	return Update{
		title:                   title,
		description:             description,
		imagePath:               imagePath,
		location:                location,
		point:                   point,
		schedule:                schedule,
		requirements:            requirements,
		volunteerLimit:          volunteerLimit,
		now:                     now,
		hasActiveParticipations: hasActiveParticipations,
	}
}
