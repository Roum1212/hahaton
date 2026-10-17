//go:build integration

package repository_event_postgresql_test

import (
	"context"
	"time"

	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	repository_fixture "github.com/Roum1212/hahaton/tests/integration/repository/fixture"
)

func (x *RepositoryTestSuite) TestRepository_Modify_IndependentEvents() {
	ctx, cancel := context.WithTimeout(x.T().Context(), 5*time.Second)
	defer cancel()

	first := x.newValue()

	second := repository_fixture.NewEvent(x.T(), first.GetCoordinatorID())
	for _, value := range []*aggregate_event.Event{first, second} {
		x.Require().NoError(x.repository.Create(ctx, value, func(ctx context.Context) error { return nil }))
	}

	entered, release := make(chan struct{}), make(chan struct{})
	firstDone := make(chan error, 1)

	go func() {
		firstDone <- x.repository.Modify(
			ctx,
			first.GetEventID(),
			func(ctx context.Context, current *aggregate_event.Event) error {
				close(entered)

				select {
				case <-release:
					return nil
				case <-ctx.Done():
					return ctx.Err()
				}
			},
		)
	}()

	select {
	case <-entered:
	case <-ctx.Done():
		x.Require().NoError(ctx.Err())
	}

	secondCtx, stop := context.WithTimeout(ctx, time.Second)
	defer stop()

	err := x.repository.Modify(
		secondCtx,
		second.GetEventID(),
		func(ctx context.Context, current *aggregate_event.Event) error {
			return nil
		},
	)

	close(release)
	x.Require().NoError(<-firstDone)
	x.Require().NoError(err, "independent events of the same coordinator must not block each other")
}
