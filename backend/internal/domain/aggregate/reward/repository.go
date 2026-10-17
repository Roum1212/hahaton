package aggregate_reward

import (
	"context"
	"errors"

	model_reward_id "github.com/Roum1212/hahaton/internal/domain/model/reward-id"
)

var (
	ErrNotFound      = errors.New("reward not found")
	ErrAlreadyExists = errors.New("reward already exists")
)

type Repository interface {
	Create(ctx context.Context, value *Reward) error
	List(ctx context.Context) ([]*Reward, error)
	Load(ctx context.Context, rewardID model_reward_id.RewardID) (*Reward, error)
	Modify(
		ctx context.Context,
		rewardID model_reward_id.RewardID,
		change func(
			ctx context.Context,
			value *Reward,
		) error,
	) error
	Save(ctx context.Context, value *Reward) error
}
