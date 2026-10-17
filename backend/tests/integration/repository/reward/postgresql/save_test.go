//go:build integration

package repository_reward_postgresql_test

import (
	aggregate_reward "github.com/Roum1212/hahaton/internal/domain/aggregate/reward"
)

func (x *RepositoryTestSuite) TestRepository_Save() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))
	x.changeValue(value)
	x.Require().NoError(x.repository.Save(ctx, value))
	got, err := x.repository.Load(ctx, value.GetRewardID())
	x.Require().NoError(err)
	assertReward(x.T(), value, got)
}

func (x *RepositoryTestSuite) TestRepository_Save_NotFound() {
	x.Require().ErrorIs(x.repository.Save(x.T().Context(), x.newValue()), aggregate_reward.ErrNotFound)
}
