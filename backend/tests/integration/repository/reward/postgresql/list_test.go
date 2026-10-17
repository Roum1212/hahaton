//go:build integration

package repository_reward_postgresql_test

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

	for _, value := range values {
		if value.GetRewardID() == first.GetRewardID() {
			assertReward(x.T(), first, value)
		} else {
			assertReward(x.T(), second, value)
		}
	}
}
