//go:build integration

package repository_session_postgresql_test

import (
	"context"
	"errors"

	aggregate_session "github.com/Roum1212/hahaton/internal/domain/aggregate/session"
	repository_fixture "github.com/Roum1212/hahaton/tests/integration/repository/fixture"
)

func (x *RepositoryTestSuite) TestRepository_Replace() {
	ctx := x.T().Context()
	original, replacement := repository_fixture.NewSession(x.T()), repository_fixture.NewSession(x.T())
	x.Require().NoError(x.repository.Create(ctx, original))
	x.Require().NoError(x.repository.Replace(
		ctx,
		original.GetSessionID(),
		func(ctx context.Context, current *aggregate_session.Session) (*aggregate_session.Session, error) {
			assertSession(x.T(), original, current)

			return replacement, nil
		},
	))
	old, err := x.repository.Load(ctx, original.GetSessionID())
	x.Require().ErrorIs(err, aggregate_session.ErrNotFound)
	x.Require().Nil(old)
	got, err := x.repository.Load(ctx, replacement.GetSessionID())
	x.Require().NoError(err)
	assertSession(x.T(), replacement, got)
}

func (x *RepositoryTestSuite) TestRepository_Replace_CallbackRollback() {
	ctx := x.T().Context()
	original := repository_fixture.NewSession(x.T())
	x.Require().NoError(x.repository.Create(ctx, original))

	stop := errors.New("cancel replacement")
	err := x.repository.Replace(
		ctx,
		original.GetSessionID(),
		func(ctx context.Context, current *aggregate_session.Session) (*aggregate_session.Session, error) {
			return nil, stop
		},
	)
	x.Require().ErrorIs(err, stop)
	got, err := x.repository.Load(ctx, original.GetSessionID())
	x.Require().NoError(err)
	assertSession(x.T(), original, got)
}

func (x *RepositoryTestSuite) TestRepository_Replace_InsertRollback() {
	ctx := x.T().Context()
	original, existing := repository_fixture.NewSession(x.T()), repository_fixture.NewSession(x.T())
	x.Require().NoError(x.repository.Create(ctx, original))
	x.Require().NoError(x.repository.Create(ctx, existing))
	err := x.repository.Replace(
		ctx,
		original.GetSessionID(),
		func(ctx context.Context, current *aggregate_session.Session) (*aggregate_session.Session, error) {
			return existing, nil
		},
	)
	x.Require().Error(err)
	got, err := x.repository.Load(ctx, original.GetSessionID())
	x.Require().NoError(err)
	assertSession(x.T(), original, got)
	got, err = x.repository.Load(ctx, existing.GetSessionID())
	x.Require().NoError(err)
	assertSession(x.T(), existing, got)
}

func (x *RepositoryTestSuite) TestRepository_Replace_NotFound() {
	ctx := x.T().Context()
	value := repository_fixture.NewSession(x.T())
	err := x.repository.Replace(
		ctx,
		value.GetSessionID(),
		func(ctx context.Context, current *aggregate_session.Session) (*aggregate_session.Session, error) {
			return value, nil
		},
	)
	x.Require().ErrorIs(err, aggregate_session.ErrNotFound)
}
