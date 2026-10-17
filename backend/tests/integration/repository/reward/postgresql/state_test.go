//go:build integration

package repository_reward_postgresql_test

import (
	aggregate_reward_command "github.com/Roum1212/hahaton/internal/domain/aggregate/reward/command"
	repository_fixture "github.com/Roum1212/hahaton/tests/integration/repository/fixture"
)

func (x *RepositoryTestSuite) TestRepository_Save_StockAndDeletion() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))
	x.Require().NoError(value.Reserve(aggregate_reward_command.NewReserve(20, repository_fixture.Now())))
	x.Require().NoError(value.Reserve(aggregate_reward_command.NewReserve(20, repository_fixture.Now())))
	x.Require().NoError(value.Issue(aggregate_reward_command.NewIssue(repository_fixture.Now())))
	x.Require().NoError(x.repository.Save(ctx, value))
	got, err := x.repository.Load(ctx, value.GetRewardID())
	x.Require().NoError(err)
	assertReward(x.T(), value, got)
	x.Require().Equal(int32(1), got.GetReservedQuantity())
	x.Require().Equal(int32(1), got.GetIssuedQuantity())
	x.Require().NoError(value.Release(aggregate_reward_command.NewRelease(repository_fixture.Now())))
	x.Require().NoError(value.Delete(aggregate_reward_command.NewDelete(repository_fixture.Now())))
	x.Require().NoError(x.repository.Save(ctx, value))
	got, err = x.repository.Load(ctx, value.GetRewardID())
	x.Require().NoError(err)
	assertReward(x.T(), value, got)
	x.Require().True(got.GetDeleted())
}
