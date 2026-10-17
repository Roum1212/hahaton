//go:build integration

package repository_coordinator_postgresql_test

import (
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
)

func (x *RepositoryTestSuite) TestRepository_CreateAndLoad() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))
	got, err := x.repository.Load(ctx, value.GetUserID())
	x.Require().NoError(err)
	assertCoordinator(x.T(), value, got)
}

func (x *RepositoryTestSuite) TestRepository_Create_Duplicate() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))
	err := x.repository.Create(ctx, value)
	x.Require().ErrorIs(err, aggregate_coordinator.ErrAlreadyExists)
	got, err := x.repository.Load(ctx, value.GetUserID())
	x.Require().NoError(err)
	assertCoordinator(x.T(), value, got)
}
