//go:build integration

package repository_cooperation_postgresql_test

import (
	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
)

func (x *RepositoryTestSuite) TestRepository_Load_NotFound() {
	value := x.newValue()
	got, err := x.repository.Load(x.T().Context(), value.GetCoordinatorID(), value.GetVolunteerID())
	x.Require().ErrorIs(err, aggregate_cooperation.ErrNotFound)
	x.Require().Nil(got)
}
