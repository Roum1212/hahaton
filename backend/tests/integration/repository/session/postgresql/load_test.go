//go:build integration

package repository_session_postgresql_test

import (
	"context"

	aggregate_session "github.com/Roum1212/hahaton/internal/domain/aggregate/session"
	repository_fixture "github.com/Roum1212/hahaton/tests/integration/repository/fixture"
)

func (x *RepositoryTestSuite) TestRepository_Load_CanceledContext() {
	ctx, cancel := context.WithCancel(x.T().Context())
	cancel()

	value, err := x.repository.Load(ctx, repository_fixture.NewSession(x.T()).GetSessionID())
	x.Require().ErrorIs(err, context.Canceled)
	x.Require().Nil(value)
}

func (x *RepositoryTestSuite) TestRepository_Load_NotFound() {
	value, err := x.repository.Load(x.T().Context(), repository_fixture.NewSession(x.T()).GetSessionID())
	x.Require().ErrorIs(err, aggregate_session.ErrNotFound)
	x.Require().Nil(value)
}
