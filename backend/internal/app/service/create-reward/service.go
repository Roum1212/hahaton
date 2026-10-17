package service_create_reward

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
	model_reward_id "github.com/Roum1212/hahaton/internal/domain/model/reward-id"
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
		CreateReward(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		rewards      aggregate_reward.Repository
		coordinators aggregate_coordinator.Repository
		images       ImageStorage
	}
)

func (x service) CreateReward(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	if _, err := x.coordinators.Load(ctx, r.GetUserID()); err != nil {
		switch {
		case errors.Is(err, aggregate_coordinator.ErrNotFound):
			return Response{}, ErrPermissionDenied
		default:
			return Response{}, err
		}
	}

	path := r.GetImagePath()
	if err := x.checkImage(ctx, r.GetUserID(), path); err != nil {
		return Response{}, err
	}

	id, err := model_reward_id.GenerateRewardID()
	if err != nil {
		return Response{}, err
	}

	reward := aggregate_reward.NewReward(id, r.GetUserID())
	if err := reward.Create(aggregate_reward_command.NewCreate(
		strings.TrimSpace(r.GetTitle()),
		r.GetDescription(),
		path,
		r.GetPriceCoins(),
		r.GetQuantityLimit(),
		time.Now().UTC(),
	)); err != nil {
		return Response{}, err
	}

	if err := x.rewards.Create(ctx, reward); err != nil {
		return Response{}, err
	}

	return NewResponse(id), nil
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
