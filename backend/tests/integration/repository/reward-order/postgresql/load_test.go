//go:build integration

package repository_reward_order_postgresql_test

import (
	aggregate_reward_order "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order"
)

func (x *RepositoryTestSuite) TestRepository_Load_NotFound() {
	value := x.newValue()
	got, err := x.repository.Load(x.T().Context(), value.GetRewardOrderID())
	x.Require().ErrorIs(err, aggregate_reward_order.ErrNotFound)
	x.Require().Nil(got)
}
