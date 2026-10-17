package service_list_rewards

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_reward "github.com/Roum1212/hahaton/internal/domain/aggregate/reward"
	model_cooperation_status "github.com/Roum1212/hahaton/internal/domain/model/cooperation-status"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	"github.com/Roum1212/hahaton/internal/pkg/page"
)

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrPermissionDenied = errors.New("permission denied")
)

type (
	Service interface {
		ListRewards(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		rewards      aggregate_reward.Repository
		coordinators aggregate_coordinator.Repository
		cooperations aggregate_cooperation.Repository
	}
)

func (x service) ListRewards(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	allowed, err := x.allowedCoordinators(ctx, r)
	if err != nil {
		return Response{}, err
	}

	rewards, err := x.rewards.List(ctx)
	if err != nil {
		return Response{}, err
	}

	filtered := make([]*aggregate_reward.Reward, 0)

	for _, reward := range rewards {
		if reward.GetDeleted() || !allowed[reward.GetCoordinatorID()] {
			continue
		}

		if value, ok := r.GetCoordinatorID().Get(); ok && reward.GetCoordinatorID() != value {
			continue
		}

		filtered = append(filtered, reward)
	}

	slices.SortFunc(filtered, func(
		a *aggregate_reward.Reward,
		b *aggregate_reward.Reward,
	) int {
		if result := strings.Compare(a.GetTitle(), b.GetTitle()); result != 0 {
			return result
		}

		return strings.Compare(a.GetRewardID().String(), b.GetRewardID().String())
	})
	values, total := page.Slice(filtered, r.GetPage())

	return NewResponse(values, total), nil
}

func (x service) allowedCoordinators(ctx context.Context, r Request) (map[model_user_id.UserID]bool, error) {
	allowed := make(map[model_user_id.UserID]bool)

	if r.GetScope() == scopeOwned {
		if _, err := x.coordinators.Load(ctx, r.GetUserID()); err != nil {
			switch {
			case errors.Is(err, aggregate_coordinator.ErrNotFound):
				return nil, ErrPermissionDenied
			default:
				return nil, err
			}
		}

		allowed[r.GetUserID()] = true

		return allowed, nil
	}

	cooperations, err := x.cooperations.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, value := range cooperations {
		if value.GetVolunteerID() == r.GetUserID() && value.GetStatus() == model_cooperation_status.Active {
			allowed[value.GetCoordinatorID()] = true
		}
	}

	return allowed, nil
}

func NewService(
	rewards aggregate_reward.Repository,
	coordinators aggregate_coordinator.Repository,
	cooperations aggregate_cooperation.Repository,
) Service {
	return service{
		rewards:      rewards,
		coordinators: coordinators,
		cooperations: cooperations,
	}
}
