package service_list_rewards

import (
	"errors"
	"fmt"

	aggregate_reward "github.com/Roum1212/hahaton/internal/domain/aggregate/reward"
	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
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

	page model_page.Page
}

func (x Request) GetCoordinatorID() utils_value.Optional[model_user_id.UserID] {
	return x.coordinatorID
}

func (x Request) GetPage() model_page.Page {
	return x.page
}

func (x Request) GetScope() int32 {
	return x.scope
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

	if err := x.page.Validate(); err != nil {
		return fmt.Errorf("page: %w", err)
	}

	if x.scope != scopeCatalog && x.scope != scopeOwned {
		return errors.New("invalid reward list scope")
	}

	return nil
}

type RequestOption func(x *Request)

func WithCoordinatorID(value model_user_id.UserID) RequestOption {
	return func(x *Request) {
		x.coordinatorID = utils_value.NewOptional(value)
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
	rewards    []*aggregate_reward.Reward
	totalCount int64
}

func (x Response) GetRewards() []*aggregate_reward.Reward {
	return x.rewards
}

func (x Response) GetTotalCount() int64 {
	return x.totalCount
}

func NewResponse(rewards []*aggregate_reward.Reward, totalCount int64) Response {
	return Response{
		rewards:    rewards,
		totalCount: totalCount,
	}
}
