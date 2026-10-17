//go:build integration

package repository_reward_order_postgresql_test

import (
	aggregate_reward_order "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order"
)

func (x *RepositoryTestSuite) TestRepository_CreateAndLoad() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))
	got, err := x.repository.Load(ctx, value.GetRewardOrderID())
	x.Require().NoError(err)
	assertRewardOrder(x.T(), value, got)
}

func (x *RepositoryTestSuite) TestRepository_Create_Duplicate() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))
	err := x.repository.Create(ctx, value)
	x.Require().ErrorIs(err, aggregate_reward_order.ErrAlreadyExists)
	got, err := x.repository.Load(ctx, value.GetRewardOrderID())
	x.Require().NoError(err)
	assertRewardOrder(x.T(), value, got)
}
