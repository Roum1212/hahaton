//go:build integration

package repository_volunteer_postgresql_test

import (
	"context"
	"errors"

	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
)

func (x *RepositoryTestSuite) TestRepository_CreateAndLoad() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value, func(ctx context.Context) error {
		return nil
	}))
	got, err := x.repository.Load(ctx, value.GetUserID())
	x.Require().NoError(err)
	assertVolunteer(x.T(), value, got)
}

func (x *RepositoryTestSuite) TestRepository_Create_Duplicate() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value, func(ctx context.Context) error {
		return nil
	}))
	err := x.repository.Create(ctx, value, func(ctx context.Context) error {
		return nil
	})
	x.Require().Error(err)
	got, err := x.repository.Load(ctx, value.GetUserID())
	x.Require().NoError(err)
	assertVolunteer(x.T(), value, got)
}

func (x *RepositoryTestSuite) TestRepository_Create_Rollback() {
	ctx := x.T().Context()
	value := x.newValue()
	stop := errors.New("cancel creation")
	err := x.repository.Create(ctx, value, func(ctx context.Context) error {
		return stop
	})
	x.Require().ErrorIs(err, stop)
	got, err := x.repository.Load(ctx, value.GetUserID())
	x.Require().ErrorIs(err, aggregate_volunteer.ErrNotFound)
	x.Require().Nil(got)
}
