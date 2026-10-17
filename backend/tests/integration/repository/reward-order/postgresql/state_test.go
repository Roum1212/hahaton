//go:build integration

package repository_reward_order_postgresql_test

import (
	aggregate_reward_order_command "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order/command"
	repository_fixture "github.com/Roum1212/hahaton/tests/integration/repository/fixture"
)

func (x *RepositoryTestSuite) TestRepository_Save_CancellationReason() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))
	x.Require().NoError(value.Cancel(aggregate_reward_order_command.NewCancel(
		"volunteer_canceled", repository_fixture.Now(),
	)))
	x.Require().NoError(x.repository.Save(ctx, value))
	got, err := x.repository.Load(ctx, value.GetRewardOrderID())
	x.Require().NoError(err)
	assertRewardOrder(x.T(), value, got)
}
