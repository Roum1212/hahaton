//go:build integration

package repository_session_postgresql_test

import (
	aggregate_session "github.com/Roum1212/hahaton/internal/domain/aggregate/session"
	repository_fixture "github.com/Roum1212/hahaton/tests/integration/repository/fixture"
)

func (x *RepositoryTestSuite) TestRepository_Delete() {
	ctx := x.T().Context()
	value := repository_fixture.NewSession(x.T())
	x.Require().NoError(x.repository.Create(ctx, value))
	x.Require().NoError(x.repository.Delete(ctx, value.GetSessionID()))
	got, err := x.repository.Load(ctx, value.GetSessionID())
	x.Require().ErrorIs(err, aggregate_session.ErrNotFound)
	x.Require().Nil(got)
	x.Require().NoError(x.repository.Delete(ctx, value.GetSessionID()))
}
