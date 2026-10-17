package service_list_participations

import (
	"errors"
	"fmt"

	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
	model_participation_status "github.com/Roum1212/hahaton/internal/domain/model/participation-status"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

const (
	scopeMyApplications int32 = 1
	scopeMyEvents       int32 = 2
)

type Request struct {
	userID model_user_id.UserID
	scope  int32

	eventID     utils_value.Optional[model_event_id.EventID]
	status      utils_value.Optional[model_participation_status.ParticipationStatus]
	volunteerID utils_value.Optional[model_user_id.UserID]

	page model_page.Page
}

func (x Request) GetEventID() utils_value.Optional[model_event_id.EventID] {
	return x.eventID
}

func (x Request) GetPage() model_page.Page {
	return x.page
}

func (x Request) GetScope() int32 {
	return x.scope
}

func (x Request) GetStatus() utils_value.Optional[model_participation_status.ParticipationStatus] {
	return x.status
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) GetVolunteerID() utils_value.Optional[model_user_id.UserID] {
	return x.volunteerID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if value, ok := x.eventID.Get(); ok {
		if err := value.Validate(); err != nil {
			return err
		}
	}

	if value, ok := x.volunteerID.Get(); ok {
		if err := value.Validate(); err != nil {
			return err
		}
	}

	if value, ok := x.status.Get(); ok {
		if err := value.Validate(); err != nil {
			return err
		}
	}

	if err := x.page.Validate(); err != nil {
		return fmt.Errorf("page: %w", err)
	}

	if x.scope != scopeMyApplications && x.scope != scopeMyEvents {
		return errors.New("invalid participation list scope")
	}

	return nil
}

type RequestOption func(x *Request)

func WithEventID(value model_event_id.EventID) RequestOption {
	return func(x *Request) {
		x.eventID = utils_value.NewOptional(value)
	}
}

func WithStatus(value model_participation_status.ParticipationStatus) RequestOption {
	return func(x *Request) {
		x.status = utils_value.NewOptional(value)
	}
}

func WithVolunteerID(value model_user_id.UserID) RequestOption {
	return func(x *Request) {
		x.volunteerID = utils_value.NewOptional(value)
	}
}

func NewRequest(userID model_user_id.UserID, scope int32, page model_page.Page, opts ...RequestOption) Request {
	x := Request{
		userID: userID,
		scope:  scope,
		page:   page,
	}
	for _, option := range opts {
		option(&x)
	}

	return x
}

type Response struct {
	participations []*aggregate_participation.Participation
	totalCount     int64
}

func (x Response) GetParticipations() []*aggregate_participation.Participation {
	return x.participations
}

func (x Response) GetTotalCount() int64 {
	return x.totalCount
}

func NewResponse(participations []*aggregate_participation.Participation, totalCount int64) Response {
	return Response{
		participations: participations,
		totalCount:     totalCount,
	}
}
