package service_list_reward_orders

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_reward_order "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order"
	"github.com/Roum1212/hahaton/internal/pkg/page"
)

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrPermissionDenied = errors.New("permission denied")
)

type (
	Service interface {
		ListRewardOrders(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		orders       aggregate_reward_order.Repository
		coordinators aggregate_coordinator.Repository
	}
)

func (x service) ListRewardOrders(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	if r.GetScope() == scopeIncoming {
		if _, err := x.coordinators.Load(ctx, r.GetUserID()); err != nil {
			switch {
			case errors.Is(err, aggregate_coordinator.ErrNotFound):
				return Response{}, ErrPermissionDenied
			default:
				return Response{}, err
			}
		}
	}

	orders, err := x.orders.List(ctx)
	if err != nil {
		return Response{}, err
	}

	filtered := make([]*aggregate_reward_order.RewardOrder, 0)

	for _, order := range orders {
		if r.GetScope() == scopeMyPurchases && order.GetVolunteerID() != r.GetUserID() {
			continue
		}

		if r.GetScope() == scopeIncoming && order.GetCoordinatorID() != r.GetUserID() {
			continue
		}

		if value, ok := r.GetCoordinatorID().Get(); ok && order.GetCoordinatorID() != value {
			continue
		}

		if value, ok := r.GetRewardID().Get(); ok && order.GetRewardID() != value {
			continue
		}

		if value, ok := r.GetVolunteerID().Get(); ok && order.GetVolunteerID() != value {
			continue
		}

		if value, ok := r.GetStatus().Get(); ok && order.GetStatus() != value {
			continue
		}

		filtered = append(filtered, order)
	}

	slices.SortFunc(filtered, func(
		a *aggregate_reward_order.RewardOrder,
		b *aggregate_reward_order.RewardOrder,
	) int {
		if result := b.GetCreatedAt().Compare(a.GetCreatedAt()); result != 0 {
			return result
		}

		return strings.Compare(b.GetRewardOrderID().String(), a.GetRewardOrderID().String())
	})
	values, total := page.Slice(filtered, r.GetPage())

	return NewResponse(values, total), nil
}

func NewService(orders aggregate_reward_order.Repository, coordinators aggregate_coordinator.Repository) Service {
	return service{
		orders:       orders,
		coordinators: coordinators,
	}
}
