//go:build integration

package repository_reward_postgresql_test

import (
	aggregate_reward "github.com/Roum1212/hahaton/internal/domain/aggregate/reward"
)

func (x *RepositoryTestSuite) TestRepository_Load_NotFound() {
	value := x.newValue()
	got, err := x.repository.Load(x.T().Context(), value.GetRewardID())
	x.Require().ErrorIs(err, aggregate_reward.ErrNotFound)
	x.Require().Nil(got)
}
