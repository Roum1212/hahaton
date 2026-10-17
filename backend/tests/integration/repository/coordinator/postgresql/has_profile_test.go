//go:build integration

package repository_coordinator_postgresql_test

func (x *RepositoryTestSuite) TestRepository_HasProfile() {
	ctx := x.T().Context()
	value := x.newValue()
	exists, err := x.repository.HasProfile(ctx, value.GetUserID())
	x.Require().NoError(err)
	x.Require().False(exists)
	x.Require().NoError(x.repository.Create(ctx, value))
	exists, err = x.repository.HasProfile(ctx, value.GetUserID())
	x.Require().NoError(err)
	x.Require().True(exists)
}
