package service_list_cooperations

import (
	"errors"
	"fmt"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	model_cooperation_status "github.com/Roum1212/hahaton/internal/domain/model/cooperation-status"
	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

const (
	scopeAsVolunteer   int32 = 1
	scopeAsCoordinator int32 = 2
)

type Request struct {
	userID model_user_id.UserID
	scope  int32

	status utils_value.Optional[model_cooperation_status.CooperationStatus]

	page model_page.Page
}

func (x Request) GetPage() model_page.Page {
	return x.page
}

func (x Request) GetScope() int32 {
	return x.scope
}

func (x Request) GetStatus() utils_value.Optional[model_cooperation_status.CooperationStatus] {
	return x.status
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if value, ok := x.status.Get(); ok {
		if err := value.Validate(); err != nil {
			return err
		}
	}

	if err := x.page.Validate(); err != nil {
		return fmt.Errorf("page: %w", err)
	}

	if x.scope != scopeAsVolunteer && x.scope != scopeAsCoordinator {
		return errors.New("invalid cooperation list scope")
	}

	return nil
}

type RequestOption func(x *Request)

func WithStatus(value model_cooperation_status.CooperationStatus) RequestOption {
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
	cooperations []*aggregate_cooperation.Cooperation
	totalCount   int64
}

func (x Response) GetCooperations() []*aggregate_cooperation.Cooperation {
	return x.cooperations
}

func (x Response) GetTotalCount() int64 {
	return x.totalCount
}

func NewResponse(cooperations []*aggregate_cooperation.Cooperation, totalCount int64) Response {
	return Response{
		cooperations: cooperations,
		totalCount:   totalCount,
	}
}
