//go:build integration

package repository_volunteer_postgresql_test

import (
	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
)

func (x *RepositoryTestSuite) TestRepository_Load_NotFound() {
	value := x.newValue()
	got, err := x.repository.Load(x.T().Context(), value.GetUserID())
	x.Require().ErrorIs(err, aggregate_volunteer.ErrNotFound)
	x.Require().Nil(got)
}
