package service_list_reward_orders

import (
	"errors"
	"fmt"

	aggregate_reward_order "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order"
	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
	model_reward_id "github.com/Roum1212/hahaton/internal/domain/model/reward-id"
	model_reward_order_status "github.com/Roum1212/hahaton/internal/domain/model/reward-order-status"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

const (
	scopeMyPurchases int32 = 1
	scopeIncoming    int32 = 2
)

type Request struct {
	userID model_user_id.UserID
	scope  int32

	coordinatorID utils_value.Optional[model_user_id.UserID]
	rewardID      utils_value.Optional[model_reward_id.RewardID]
	status        utils_value.Optional[model_reward_order_status.RewardOrderStatus]
	volunteerID   utils_value.Optional[model_user_id.UserID]

	page model_page.Page
}

func (x Request) GetCoordinatorID() utils_value.Optional[model_user_id.UserID] {
	return x.coordinatorID
}

func (x Request) GetPage() model_page.Page {
	return x.page
}

func (x Request) GetRewardID() utils_value.Optional[model_reward_id.RewardID] {
	return x.rewardID
}

func (x Request) GetScope() int32 {
	return x.scope
}

func (x Request) GetStatus() utils_value.Optional[model_reward_order_status.RewardOrderStatus] {
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

	if value, ok := x.coordinatorID.Get(); ok {
		if err := value.Validate(); err != nil {
			return err
		}
	}

	if value, ok := x.rewardID.Get(); ok {
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

	if x.scope != scopeMyPurchases && x.scope != scopeIncoming {
		return errors.New("invalid reward order list scope")
	}

	return nil
}

type RequestOption func(x *Request)

func WithCoordinatorID(value model_user_id.UserID) RequestOption {
	return func(x *Request) {
		x.coordinatorID = utils_value.NewOptional(value)
	}
}

func WithRewardID(value model_reward_id.RewardID) RequestOption {
	return func(x *Request) {
		x.rewardID = utils_value.NewOptional(value)
	}
}

func WithStatus(value model_reward_order_status.RewardOrderStatus) RequestOption {
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
	orders     []*aggregate_reward_order.RewardOrder
	totalCount int64
}

func (x Response) GetOrders() []*aggregate_reward_order.RewardOrder {
	return x.orders
}

func (x Response) GetTotalCount() int64 {
	return x.totalCount
}

func NewResponse(orders []*aggregate_reward_order.RewardOrder, totalCount int64) Response {
	return Response{
		orders:     orders,
		totalCount: totalCount,
	}
}
