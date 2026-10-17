package service_update_reward

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_reward "github.com/Roum1212/hahaton/internal/domain/aggregate/reward"
	aggregate_reward_command "github.com/Roum1212/hahaton/internal/domain/aggregate/reward/command"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrNotFound         = errors.New("not found")
	ErrPermissionDenied = errors.New("permission denied")
)

type (
	ImageStorage interface {
		Exists(ctx context.Context, path model_image_path.ImagePath) (bool, error)
	}
	Service interface {
		UpdateReward(ctx context.Context, r Request) error
	}
	service struct {
		rewards      aggregate_reward.Repository
		coordinators aggregate_coordinator.Repository
		images       ImageStorage
	}
)

func (x service) UpdateReward(ctx context.Context, r Request) error {
	if err := r.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	if _, err := x.coordinators.Load(ctx, r.GetUserID()); err != nil {
		switch {
		case errors.Is(err, aggregate_coordinator.ErrNotFound):
			return ErrPermissionDenied
		default:
			return err
		}
	}

	if path, ok := r.GetImagePath().Get(); ok {
		if err := x.checkImage(ctx, r.GetUserID(), path); err != nil {
			return err
		}
	}

	return x.rewards.Modify(ctx, r.GetRewardID(), func(
		ctx context.Context,
		reward *aggregate_reward.Reward,
	) error {
		if reward.GetCoordinatorID() != r.GetUserID() {
			return ErrPermissionDenied
		}

		title, description := reward.GetTitle(), reward.GetDescription()
		if value, ok := r.GetTitle().Get(); ok {
			title = strings.TrimSpace(value)
		}

		if value, ok := r.GetDescription().Get(); ok {
			description = value
		}

		path := reward.GetImagePath()
		if value, ok := r.GetImagePath().Get(); ok {
			path = value
		}

		price, quantity := reward.GetPriceCoins(), reward.GetQuantityLimit()
		if value, ok := r.GetPriceCoins().Get(); ok {
			price = value
		}

		if value, ok := r.GetQuantityLimit().Get(); ok {
			quantity = value
		}

		return reward.Update(aggregate_reward_command.NewUpdate(
			title,
			description,
			path,
			price,
			quantity,
			time.Now().UTC(),
		))
	})
}

func (x service) checkImage(ctx context.Context, owner model_user_id.UserID, path model_image_path.ImagePath) error {
	if path.IsZero() {
		return nil
	}

	if path.GetUserID() != owner {
		return ErrPermissionDenied
	}

	exists, err := x.images.Exists(ctx, path)
	if err != nil {
		return err
	}

	if !exists {
		return ErrNotFound
	}

	return nil
}

func NewService(
	rewards aggregate_reward.Repository,
	coordinators aggregate_coordinator.Repository,
	images ImageStorage,
) Service {
	return service{
		rewards:      rewards,
		coordinators: coordinators,
		images:       images,
	}
}
