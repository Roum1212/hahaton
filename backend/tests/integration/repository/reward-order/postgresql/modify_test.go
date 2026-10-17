//go:build integration

package repository_reward_order_postgresql_test

import (
	"context"
	"errors"

	aggregate_reward_order "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order"
)

func (x *RepositoryTestSuite) TestRepository_Modify() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))
	err := x.repository.Modify(
		ctx,
		value.GetRewardOrderID(),
		func(ctx context.Context, current *aggregate_reward_order.RewardOrder) error {
			x.changeValue(current)

			return nil
		},
	)
	x.Require().NoError(err)
	x.changeValue(value)
	got, err := x.repository.Load(ctx, value.GetRewardOrderID())
	x.Require().NoError(err)
	assertRewardOrder(x.T(), value, got)
}

func (x *RepositoryTestSuite) TestRepository_Modify_NotFound() {
	ctx := x.T().Context()
	value := x.newValue()
	err := x.repository.Modify(
		ctx,
		value.GetRewardOrderID(),
		func(ctx context.Context, current *aggregate_reward_order.RewardOrder) error {
			return nil
		},
	)
	x.Require().ErrorIs(err, aggregate_reward_order.ErrNotFound)
}

func (x *RepositoryTestSuite) TestRepository_Modify_Rollback() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))

	stop := errors.New("cancel modification")
	err := x.repository.Modify(
		ctx,
		value.GetRewardOrderID(),
		func(ctx context.Context, current *aggregate_reward_order.RewardOrder) error {
			x.changeValue(current)

			if err := x.repository.Save(ctx, current); err != nil {
				return err
			}

			return stop
		},
	)
	x.Require().ErrorIs(err, stop)
	got, err := x.repository.Load(ctx, value.GetRewardOrderID())
	x.Require().NoError(err)
	assertRewardOrder(x.T(), value, got)
}
