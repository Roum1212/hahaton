package service_get_reward_order

import (
	"context"
	"errors"
	"fmt"

	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_reward_order "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order"
)

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrPermissionDenied = errors.New("permission denied")
)

type (
	Service interface {
		GetRewardOrder(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		orders       aggregate_reward_order.Repository
		coordinators aggregate_coordinator.Repository
	}
)

func (x service) GetRewardOrder(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	order, err := x.orders.Load(ctx, r.GetRewardOrderID())
	if err != nil {
		return Response{}, err
	}

	if order.GetVolunteerID() != r.GetUserID() {
		if order.GetCoordinatorID() != r.GetUserID() {
			return Response{}, ErrPermissionDenied
		}

		if _, err := x.coordinators.Load(ctx, r.GetUserID()); err != nil {
			switch {
			case errors.Is(err, aggregate_coordinator.ErrNotFound):
				return Response{}, ErrPermissionDenied
			default:
				return Response{}, err
			}
		}
	}

	return NewResponse(order), nil
}

func NewService(orders aggregate_reward_order.Repository, coordinators aggregate_coordinator.Repository) Service {
	return service{
		orders:       orders,
		coordinators: coordinators,
	}
}
