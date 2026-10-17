package service_get_reward

import (
	"context"
	"errors"
	"fmt"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_reward "github.com/Roum1212/hahaton/internal/domain/aggregate/reward"
	model_cooperation_status "github.com/Roum1212/hahaton/internal/domain/model/cooperation-status"
)

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrPermissionDenied = errors.New("permission denied")
)

type (
	Service interface {
		GetReward(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		rewards      aggregate_reward.Repository
		cooperations aggregate_cooperation.Repository
	}
)

func (x service) GetReward(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	reward, err := x.rewards.Load(ctx, r.GetRewardID())
	if err != nil {
		return Response{}, err
	}

	if reward.GetDeleted() {
		return Response{}, aggregate_reward.ErrNotFound
	}

	if reward.GetCoordinatorID() != r.GetUserID() {
		cooperation, err := x.cooperations.Load(ctx, reward.GetCoordinatorID(), r.GetUserID())
		if err != nil {
			switch {
			case errors.Is(err, aggregate_cooperation.ErrNotFound):
				return Response{}, ErrPermissionDenied
			default:
				return Response{}, err
			}
		}

		if cooperation.GetStatus() != model_cooperation_status.Active {
			return Response{}, ErrPermissionDenied
		}
	}

	return NewResponse(reward), nil
}

func NewService(rewards aggregate_reward.Repository, cooperations aggregate_cooperation.Repository) Service {
	return service{
		rewards:      rewards,
		cooperations: cooperations,
	}
}
