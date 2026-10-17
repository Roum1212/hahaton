package aggregate_reward_order

import (
	"context"
	"errors"

	model_reward_order_id "github.com/Roum1212/hahaton/internal/domain/model/reward-order-id"
)

var (
	ErrNotFound      = errors.New("reward order not found")
	ErrAlreadyExists = errors.New("reward order already exists")
)

type Repository interface {
	Create(ctx context.Context, value *RewardOrder) error
	List(ctx context.Context) ([]*RewardOrder, error)
	Load(ctx context.Context, rewardOrderID model_reward_order_id.RewardOrderID) (*RewardOrder, error)
	Modify(
		ctx context.Context,
		rewardOrderID model_reward_order_id.RewardOrderID,
		change func(
			ctx context.Context,
			value *RewardOrder,
		) error,
	) error
	Save(ctx context.Context, value *RewardOrder) error
}
