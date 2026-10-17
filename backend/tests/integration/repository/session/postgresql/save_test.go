//go:build integration

package repository_session_postgresql_test

import (
	"time"

	aggregate_session "github.com/Roum1212/hahaton/internal/domain/aggregate/session"
	aggregate_session_command "github.com/Roum1212/hahaton/internal/domain/aggregate/session/command"
	repository_fixture "github.com/Roum1212/hahaton/tests/integration/repository/fixture"
)

func (x *RepositoryTestSuite) TestRepository_Save_Authentication() {
	ctx := x.T().Context()
	value := repository_fixture.NewSession(x.T())
	x.Require().NoError(x.repository.Create(ctx, value))
	user := repository_fixture.NewVolunteer(x.T()).GetUserID()
	x.Require().NoError(value.Authenticate(aggregate_session_command.NewAuthenticate(
		user, repository_fixture.Now().Add(time.Minute), repository_fixture.Now().Add(2*time.Hour),
	)))
	x.Require().NoError(x.repository.Save(ctx, value))
	got, err := x.repository.Load(ctx, value.GetSessionID())
	x.Require().NoError(err)
	assertSession(x.T(), value, got)
	x.Require().False(got.IsAnonymous())
}

func (x *RepositoryTestSuite) TestRepository_Save_NotFound() {
	x.Require().ErrorIs(x.repository.Save(x.T().Context(), repository_fixture.NewSession(
		x.T(),
	)), aggregate_session.ErrNotFound)
}
