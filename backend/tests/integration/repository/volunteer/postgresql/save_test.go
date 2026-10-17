//go:build integration

package repository_volunteer_postgresql_test

import (
	"context"

	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
)

func (x *RepositoryTestSuite) TestRepository_Save() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value, func(ctx context.Context) error {
		return nil
	}))
	x.changeValue(value)
	x.Require().NoError(x.repository.Save(ctx, value))
	got, err := x.repository.Load(ctx, value.GetUserID())
	x.Require().NoError(err)
	assertVolunteer(x.T(), value, got)
}

func (x *RepositoryTestSuite) TestRepository_Save_NotFound() {
	x.Require().ErrorIs(x.repository.Save(x.T().Context(), x.newValue()), aggregate_volunteer.ErrNotFound)
}
