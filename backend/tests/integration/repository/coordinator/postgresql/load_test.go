//go:build integration

package repository_coordinator_postgresql_test

import (
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
)

func (x *RepositoryTestSuite) TestRepository_Load_NotFound() {
	value := x.newValue()
	got, err := x.repository.Load(x.T().Context(), value.GetUserID())
	x.Require().ErrorIs(err, aggregate_coordinator.ErrNotFound)
	x.Require().Nil(got)
}
