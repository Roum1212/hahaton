//go:build integration

package repository_reward_order_postgresql_test

import (
	aggregate_reward_order "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order"
)

func (x *RepositoryTestSuite) TestRepository_Save() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))
	x.changeValue(value)
	x.Require().NoError(x.repository.Save(ctx, value))
	got, err := x.repository.Load(ctx, value.GetRewardOrderID())
	x.Require().NoError(err)
	assertRewardOrder(x.T(), value, got)
}

func (x *RepositoryTestSuite) TestRepository_Save_NotFound() {
	x.Require().ErrorIs(x.repository.Save(x.T().Context(), x.newValue()), aggregate_reward_order.ErrNotFound)
}
