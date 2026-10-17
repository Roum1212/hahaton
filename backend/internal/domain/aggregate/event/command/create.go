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

type Create struct {
	title          model_event_title.EventTitle
	description    string
	imagePath      model_image_path.ImagePath
	location       string
	point          model_geo_point.GeoPoint
	schedule       model_event_schedule.EventSchedule
	requirements   model_event_requirements.EventRequirements
	volunteerLimit model_volunteer_limit.VolunteerLimit
	now            time.Time
}

func (x Create) GetDescription() string {
	return x.description
}

func (x Create) GetImagePath() model_image_path.ImagePath {
	return x.imagePath
}

func (x Create) GetLocation() string {
	return x.location
}

func (x Create) GetNow() time.Time {
	return x.now
}

func (x Create) GetPoint() model_geo_point.GeoPoint {
	return x.point
}

func (x Create) GetRequirements() model_event_requirements.EventRequirements {
	return x.requirements
}

func (x Create) GetSchedule() model_event_schedule.EventSchedule {
	return x.schedule
}

func (x Create) GetTitle() model_event_title.EventTitle {
	return x.title
}

func (x Create) GetVolunteerLimit() model_volunteer_limit.VolunteerLimit {
	return x.volunteerLimit
}

func NewCreate(
	title model_event_title.EventTitle,
	description string,
	imagePath model_image_path.ImagePath,
	location string,
	point model_geo_point.GeoPoint,
	schedule model_event_schedule.EventSchedule,
	requirements model_event_requirements.EventRequirements,
	volunteerLimit model_volunteer_limit.VolunteerLimit,
	now time.Time,
) Create {
	return Create{
		title:          title,
		description:    description,
		imagePath:      imagePath,
		location:       location,
		point:          point,
		schedule:       schedule,
		requirements:   requirements,
		volunteerLimit: volunteerLimit,
		now:            now,
	}
}
