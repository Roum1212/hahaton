//go:build integration

package repository_event_postgresql_test

import (
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
)

func (x *RepositoryTestSuite) TestRepository_Load_NotFound() {
	value := x.newValue()
	got, err := x.repository.Load(x.T().Context(), value.GetEventID())
	x.Require().ErrorIs(err, aggregate_event.ErrNotFound)
	x.Require().Nil(got)
}
