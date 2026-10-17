//go:build integration

package repository_event_postgresql_test

import (
	"context"
	"errors"

	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
)

func (x *RepositoryTestSuite) TestRepository_CreateAndLoad() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value, func(ctx context.Context) error {
		return nil
	}))
	got, err := x.repository.Load(ctx, value.GetEventID())
	x.Require().NoError(err)
	assertEvent(x.T(), value, got)
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
	x.Require().ErrorIs(err, aggregate_event.ErrAlreadyExists)
	got, err := x.repository.Load(ctx, value.GetEventID())
	x.Require().NoError(err)
	assertEvent(x.T(), value, got)
}

func (x *RepositoryTestSuite) TestRepository_Create_Rollback() {
	ctx := x.T().Context()
	value := x.newValue()
	stop := errors.New("cancel creation")
	err := x.repository.Create(ctx, value, func(ctx context.Context) error {
		return stop
	})
	x.Require().ErrorIs(err, stop)
	got, err := x.repository.Load(ctx, value.GetEventID())
	x.Require().ErrorIs(err, aggregate_event.ErrNotFound)
	x.Require().Nil(got)
}
