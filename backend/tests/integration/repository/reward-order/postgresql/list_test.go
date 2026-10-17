//go:build integration

package repository_reward_order_postgresql_test

import (
	aggregate_reward_order "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order"
)

func (x *RepositoryTestSuite) TestRepository_List() {
	ctx := x.T().Context()
	values, err := x.repository.List(ctx)
	x.Require().NoError(err)
	x.Require().Empty(values)
	first, second := x.newValue(), x.newValue()
	x.Require().NoError(x.repository.Create(ctx, first))
	x.Require().NoError(x.repository.Create(ctx, second))
	values, err = x.repository.List(ctx)
	x.Require().NoError(err)
	x.Require().Len(values, 2)
	x.Require().ElementsMatch([]*aggregate_reward_order.RewardOrder{first, second}, values)
}
