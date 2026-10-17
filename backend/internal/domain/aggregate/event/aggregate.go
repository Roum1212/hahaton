package aggregate_event

import (
	"errors"
	"time"

	aggregate_event_command "github.com/Roum1212/hahaton/internal/domain/aggregate/event/command"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_event_requirements "github.com/Roum1212/hahaton/internal/domain/model/event-requirements"
	model_event_schedule "github.com/Roum1212/hahaton/internal/domain/model/event-schedule"
	model_event_status "github.com/Roum1212/hahaton/internal/domain/model/event-status"
	model_event_title "github.com/Roum1212/hahaton/internal/domain/model/event-title"
	model_geo_point "github.com/Roum1212/hahaton/internal/domain/model/geo-point"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	model_volunteer_limit "github.com/Roum1212/hahaton/internal/domain/model/volunteer-limit"
)

var (
	ErrInvalidSchedule = errors.New("event must start in the future")
	ErrInvalidState    = errors.New("invalid event state")
	ErrNotEditable     = errors.New("event cannot be edited")
	ErrNotFinished     = errors.New("event has not ended")
)

type Event struct {
	eventID             model_event_id.EventID
	coordinatorID       model_user_id.UserID
	title               model_event_title.EventTitle
	description         string
	imagePath           model_image_path.ImagePath
	location            string
	point               model_geo_point.GeoPoint
	schedule            model_event_schedule.EventSchedule
	requirements        model_event_requirements.EventRequirements
	volunteerLimit      model_volunteer_limit.VolunteerLimit
	status              model_event_status.EventStatus
	completionRequestID string
	resultFingerprint   string
	createdAt           time.Time
	updatedAt           time.Time
}

func (x *Event) GetCompletionRequestID() string {
	return x.completionRequestID
}

func (x *Event) GetCoordinatorID() model_user_id.UserID {
	return x.coordinatorID
}

func (x *Event) GetCreatedAt() time.Time {
	return x.createdAt
}

func (x *Event) GetDescription() string {
	return x.description
}

func (x *Event) GetEventID() model_event_id.EventID {
	return x.eventID
}

func (x *Event) GetImagePath() model_image_path.ImagePath {
	return x.imagePath
}

func (x *Event) GetLocation() string {
	return x.location
}

func (x *Event) GetPoint() model_geo_point.GeoPoint {
	return x.point
}

func (x *Event) GetRequirements() model_event_requirements.EventRequirements {
	return x.requirements
}

func (x *Event) GetResultFingerprint() string {
	return x.resultFingerprint
}

func (x *Event) GetSchedule() model_event_schedule.EventSchedule {
	return x.schedule
}

func (x *Event) GetStatus() model_event_status.EventStatus {
	return x.status
}

func (x *Event) GetTitle() model_event_title.EventTitle {
	return x.title
}

func (x *Event) GetUpdatedAt() time.Time {
	return x.updatedAt
}

func (x *Event) GetVolunteerLimit() model_volunteer_limit.VolunteerLimit {
	return x.volunteerLimit
}

func (x *Event) Cancel(command aggregate_event_command.Cancel) error {
	if x.status != model_event_status.Published || !command.GetNow().Before(x.schedule.GetStartsAt()) {
		return ErrInvalidState
	}

	x.status = model_event_status.Canceled
	x.updatedAt = command.GetNow()

	return nil
}

func (x *Event) Complete(command aggregate_event_command.Complete) error {
	if x.status != model_event_status.Published {
		return ErrInvalidState
	}

	if !command.IsForced() && command.GetNow().Before(x.schedule.GetEndsAt()) {
		return ErrNotFinished
	}

	x.status = model_event_status.Completed
	x.completionRequestID = command.GetRequestID()
	x.resultFingerprint = command.GetResultFingerprint()
	x.updatedAt = command.GetNow()

	return nil
}

func (x *Event) Create(command aggregate_event_command.Create) error {
	if !x.createdAt.IsZero() {
		return ErrInvalidState
	}

	if !command.GetSchedule().GetStartsAt().After(command.GetNow()) {
		return ErrInvalidSchedule
	}

	x.title = command.GetTitle()
	x.description = command.GetDescription()
	x.imagePath = command.GetImagePath()
	x.location = command.GetLocation()
	x.point = command.GetPoint()
	x.schedule = command.GetSchedule()
	x.requirements = command.GetRequirements()
	x.volunteerLimit = command.GetVolunteerLimit()
	x.status = model_event_status.Published
	x.createdAt = command.GetNow()
	x.updatedAt = command.GetNow()

	return nil
}

func (x *Event) Update(command aggregate_event_command.Update) error {
	if x.status != model_event_status.Published ||
		!command.GetNow().Before(x.schedule.GetStartsAt()) || command.GetHasActiveParticipations() {
		return ErrNotEditable
	}

	if !command.GetSchedule().GetStartsAt().After(command.GetNow()) {
		return ErrInvalidSchedule
	}

	x.title = command.GetTitle()
	x.description = command.GetDescription()
	x.imagePath = command.GetImagePath()
	x.location = command.GetLocation()
	x.point = command.GetPoint()
	x.schedule = command.GetSchedule()
	x.requirements = command.GetRequirements()
	x.volunteerLimit = command.GetVolunteerLimit()
	x.updatedAt = command.GetNow()

	return nil
}

func NewEvent(eventID model_event_id.EventID, coordinatorID model_user_id.UserID) *Event {
	return &Event{
		eventID:       eventID,
		coordinatorID: coordinatorID,
	}
}

func NewEventFromFields(
	eventID model_event_id.EventID,
	coordinatorID model_user_id.UserID,
	title model_event_title.EventTitle,
	description string,
	imagePath model_image_path.ImagePath,
	location string,
	point model_geo_point.GeoPoint,
	schedule model_event_schedule.EventSchedule,
	requirements model_event_requirements.EventRequirements,
	volunteerLimit model_volunteer_limit.VolunteerLimit,
	status model_event_status.EventStatus,
	completionRequestID string,
	resultFingerprint string,
	createdAt time.Time,
	updatedAt time.Time,
) *Event {
	return &Event{
		eventID:        eventID,
		coordinatorID:  coordinatorID,
		title:          title,
		description:    description,
		imagePath:      imagePath,
		location:       location,
		point:          point,
		schedule:       schedule,
		requirements:   requirements,
		volunteerLimit: volunteerLimit,

		status:              status,
		completionRequestID: completionRequestID,
		resultFingerprint:   resultFingerprint,
		createdAt:           createdAt,
		updatedAt:           updatedAt,
	}
}
