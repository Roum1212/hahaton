package service_list_events

import (
	"errors"
	"fmt"
	"time"

	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_event_status "github.com/Roum1212/hahaton/internal/domain/model/event-status"
	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
	model_search_area "github.com/Roum1212/hahaton/internal/domain/model/search-area"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

const (
	scopeCatalog int32 = 1
	scopeOwned   int32 = 2
)

type Request struct {
	userID model_user_id.UserID
	scope  int32

	coordinatorID utils_value.Optional[model_user_id.UserID]
	searchArea    utils_value.Optional[model_search_area.SearchArea]
	status        utils_value.Optional[model_event_status.EventStatus]

	startsAtFrom utils_value.Optional[time.Time]
	startsAtTo   utils_value.Optional[time.Time]

	endsAtFrom utils_value.Optional[time.Time]
	endsAtTo   utils_value.Optional[time.Time]

	page model_page.Page
}

func (x Request) GetCoordinatorID() utils_value.Optional[model_user_id.UserID] {
	return x.coordinatorID
}

func (x Request) GetEndsAtFrom() utils_value.Optional[time.Time] {
	return x.endsAtFrom
}

func (x Request) GetEndsAtTo() utils_value.Optional[time.Time] {
	return x.endsAtTo
}

func (x Request) GetPage() model_page.Page {
	return x.page
}

func (x Request) GetScope() int32 {
	return x.scope
}

func (x Request) GetSearchArea() utils_value.Optional[model_search_area.SearchArea] {
	return x.searchArea
}

func (x Request) GetStartsAtFrom() utils_value.Optional[time.Time] {
	return x.startsAtFrom
}

func (x Request) GetStartsAtTo() utils_value.Optional[time.Time] {
	return x.startsAtTo
}

func (x Request) GetStatus() utils_value.Optional[model_event_status.EventStatus] {
	return x.status
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if value, ok := x.coordinatorID.Get(); ok {
		if err := value.Validate(); err != nil {
			return err
		}
	}

	if value, ok := x.searchArea.Get(); ok {
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

	if x.scope != scopeCatalog && x.scope != scopeOwned {
		return errors.New("invalid event list scope")
	}

	if _, ok := x.searchArea.Get(); x.scope == scopeCatalog && !ok {
		return errors.New("catalog search area is required")
	}

	startsAtFrom, hasStartsAtFrom := x.startsAtFrom.Get()

	startsAtTo, hasStartsAtTo := x.startsAtTo.Get()
	if hasStartsAtFrom && hasStartsAtTo && !startsAtFrom.Before(startsAtTo) {
		return errors.New("invalid start interval")
	}

	endsAtFrom, hasEndsAtFrom := x.endsAtFrom.Get()

	endsAtTo, hasEndsAtTo := x.endsAtTo.Get()
	if hasEndsAtFrom && hasEndsAtTo && !endsAtFrom.Before(endsAtTo) {
		return errors.New("invalid end interval")
	}

	return nil
}

type RequestOption func(x *Request)

func WithCoordinatorID(value model_user_id.UserID) RequestOption {
	return func(x *Request) {
		x.coordinatorID = utils_value.NewOptional(value)
	}
}

func WithEndsAtFrom(value time.Time) RequestOption {
	return func(x *Request) {
		x.endsAtFrom = utils_value.NewOptional(value)
	}
}

func WithEndsAtTo(value time.Time) RequestOption {
	return func(x *Request) {
		x.endsAtTo = utils_value.NewOptional(value)
	}
}

func WithSearchArea(value model_search_area.SearchArea) RequestOption {
	return func(x *Request) {
		x.searchArea = utils_value.NewOptional(value)
	}
}

func WithStartsAtFrom(value time.Time) RequestOption {
	return func(x *Request) {
		x.startsAtFrom = utils_value.NewOptional(value)
	}
}

func WithStartsAtTo(value time.Time) RequestOption {
	return func(x *Request) {
		x.startsAtTo = utils_value.NewOptional(value)
	}
}

func WithStatus(value model_event_status.EventStatus) RequestOption {
	return func(x *Request) {
		x.status = utils_value.NewOptional(value)
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
	events         []*aggregate_event.Event
	approvedCounts map[model_event_id.EventID]int32
	pendingCounts  map[model_event_id.EventID]int32
	totalCount     int64
}

func (x Response) GetApprovedCounts() map[model_event_id.EventID]int32 {
	return x.approvedCounts
}

func (x Response) GetEvents() []*aggregate_event.Event {
	return x.events
}

func (x Response) GetPendingCounts() map[model_event_id.EventID]int32 {
	return x.pendingCounts
}

func (x Response) GetTotalCount() int64 {
	return x.totalCount
}

func NewResponse(
	events []*aggregate_event.Event,
	approvedCounts map[model_event_id.EventID]int32,
	pendingCounts map[model_event_id.EventID]int32,
	totalCount int64,
) Response {
	return Response{
		events:         events,
		approvedCounts: approvedCounts,
		pendingCounts:  pendingCounts,
		totalCount:     totalCount,
	}
}
