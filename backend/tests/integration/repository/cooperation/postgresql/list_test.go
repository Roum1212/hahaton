//go:build integration

package repository_cooperation_postgresql_test

import (
	"context"
)

func (x *RepositoryTestSuite) TestRepository_List() {
	ctx := x.T().Context()
	values, err := x.repository.List(ctx)
	x.Require().NoError(err)
	x.Require().Empty(values)
	first, second := x.newValue(), x.newValue()
	x.Require().NoError(x.repository.Create(ctx, first, func(ctx context.Context) error {
		return nil
	}))
	x.Require().NoError(x.repository.Create(ctx, second, func(ctx context.Context) error {
		return nil
	}))
	values, err = x.repository.List(ctx)
	x.Require().NoError(err)
	x.Require().Len(values, 2)

	for _, value := range values {
		if value.GetVolunteerID() == first.GetVolunteerID() {
			assertCooperation(x.T(), first, value)
		} else {
			assertCooperation(x.T(), second, value)
		}
	}
}
