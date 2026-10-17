//go:build integration

package repository_event_postgresql_test

import (
	"context"
	"errors"

	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
)

func (x *RepositoryTestSuite) TestRepository_Modify() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value, func(ctx context.Context) error {
		return nil
	}))
	err := x.repository.Modify(ctx, value.GetEventID(), func(ctx context.Context, current *aggregate_event.Event) error {
		x.changeValue(current)

		return nil
	})
	x.Require().NoError(err)
	x.changeValue(value)
	got, err := x.repository.Load(ctx, value.GetEventID())
	x.Require().NoError(err)
	assertEvent(x.T(), value, got)
}

func (x *RepositoryTestSuite) TestRepository_Modify_NotFound() {
	ctx := x.T().Context()
	value := x.newValue()
	err := x.repository.Modify(
		ctx,
		value.GetEventID(),
		func(ctx context.Context, current *aggregate_event.Event) error {
			return nil
		},
	)
	x.Require().ErrorIs(err, aggregate_event.ErrNotFound)
}

func (x *RepositoryTestSuite) TestRepository_Modify_Rollback() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value, func(ctx context.Context) error {
		return nil
	}))

	stop := errors.New("cancel modification")
	err := x.repository.Modify(ctx, value.GetEventID(), func(ctx context.Context, current *aggregate_event.Event) error {
		x.changeValue(current)

		if err := x.repository.Save(ctx, current); err != nil {
			return err
		}

		return stop
	})
	x.Require().ErrorIs(err, stop)
	got, err := x.repository.Load(ctx, value.GetEventID())
	x.Require().NoError(err)
	assertEvent(x.T(), value, got)
}
