package service_update_event

import (
	"fmt"
	"time"

	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_event_requirements "github.com/Roum1212/hahaton/internal/domain/model/event-requirements"
	model_event_title "github.com/Roum1212/hahaton/internal/domain/model/event-title"
	model_geo_point "github.com/Roum1212/hahaton/internal/domain/model/geo-point"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	model_volunteer_limit "github.com/Roum1212/hahaton/internal/domain/model/volunteer-limit"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

type Request struct {
	userID  model_user_id.UserID
	eventID model_event_id.EventID

	description utils_value.Optional[string]
	imagePath   utils_value.Optional[model_image_path.ImagePath]
	title       utils_value.Optional[model_event_title.EventTitle]

	location utils_value.Optional[string]
	point    utils_value.Optional[model_geo_point.GeoPoint]

	startsAt utils_value.Optional[time.Time]
	endsAt   utils_value.Optional[time.Time]

	requirements   utils_value.Optional[model_event_requirements.EventRequirements]
	volunteerLimit utils_value.Optional[model_volunteer_limit.VolunteerLimit]
}

func (x Request) GetDescription() utils_value.Optional[string] {
	return x.description
}

func (x Request) GetEndsAt() utils_value.Optional[time.Time] {
	return x.endsAt
}

func (x Request) GetEventID() model_event_id.EventID {
	return x.eventID
}

func (x Request) GetImagePath() utils_value.Optional[model_image_path.ImagePath] {
	return x.imagePath
}

func (x Request) GetLocation() utils_value.Optional[string] {
	return x.location
}

func (x Request) GetPoint() utils_value.Optional[model_geo_point.GeoPoint] {
	return x.point
}

func (x Request) GetRequirements() utils_value.Optional[model_event_requirements.EventRequirements] {
	return x.requirements
}

func (x Request) GetStartsAt() utils_value.Optional[time.Time] {
	return x.startsAt
}

func (x Request) GetTitle() utils_value.Optional[model_event_title.EventTitle] {
	return x.title
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) GetVolunteerLimit() utils_value.Optional[model_volunteer_limit.VolunteerLimit] {
	return x.volunteerLimit
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if err := x.eventID.Validate(); err != nil {
		return fmt.Errorf("event id: %w", err)
	}

	if value, ok := x.title.Get(); ok {
		if err := value.Validate(); err != nil {
			return err
		}
	}

	if value, ok := x.point.Get(); ok {
		if err := value.Validate(); err != nil {
			return err
		}
	}

	if value, ok := x.requirements.Get(); ok {
		if err := value.Validate(); err != nil {
			return err
		}
	}

	if value, ok := x.volunteerLimit.Get(); ok {
		if err := value.Validate(); err != nil {
			return err
		}
	}

	if value, ok := x.imagePath.Get(); ok {
		if err := value.Validate(); err != nil {
			return fmt.Errorf("image path: %w", err)
		}
	}

	return nil
}

type RequestOption func(x *Request)

func WithDescription(value string) RequestOption {
	return func(x *Request) {
		x.description = utils_value.NewOptional(value)
	}
}

func WithEndsAt(value time.Time) RequestOption {
	return func(x *Request) {
		x.endsAt = utils_value.NewOptional(value)
	}
}

func WithImagePath(value model_image_path.ImagePath) RequestOption {
	return func(x *Request) {
		x.imagePath = utils_value.NewOptional(value)
	}
}

func WithLocation(value string) RequestOption {
	return func(x *Request) {
		x.location = utils_value.NewOptional(value)
	}
}

func WithPoint(value model_geo_point.GeoPoint) RequestOption {
	return func(x *Request) {
		x.point = utils_value.NewOptional(value)
	}
}

func WithRequirements(value model_event_requirements.EventRequirements) RequestOption {
	return func(x *Request) {
		x.requirements = utils_value.NewOptional(value)
	}
}

func WithStartsAt(value time.Time) RequestOption {
	return func(x *Request) {
		x.startsAt = utils_value.NewOptional(value)
	}
}

func WithTitle(value model_event_title.EventTitle) RequestOption {
	return func(x *Request) {
		x.title = utils_value.NewOptional(value)
	}
}

func WithVolunteerLimit(value model_volunteer_limit.VolunteerLimit) RequestOption {
	return func(x *Request) {
		x.volunteerLimit = utils_value.NewOptional(value)
	}
}

func NewRequest(userID model_user_id.UserID, eventID model_event_id.EventID, opts ...RequestOption) Request {
	x := Request{
		userID:  userID,
		eventID: eventID,
	}
	for _, option := range opts {
		option(&x)
	}

	return x
}
