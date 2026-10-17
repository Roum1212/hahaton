//go:build integration

package repository_volunteer_postgresql_test

import (
	"context"

	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
	model_identity "github.com/Roum1212/hahaton/internal/domain/model/identity"
)

func (x *RepositoryTestSuite) TestRepository_FindByIdentity() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value, func(ctx context.Context) error {
		return nil
	}))
	got, err := x.repository.FindByIdentity(ctx, value.GetIdentity())
	x.Require().NoError(err)
	assertVolunteer(x.T(), value, got)
}

func (x *RepositoryTestSuite) TestRepository_FindByIdentity_NotFound() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value, func(ctx context.Context) error {
		return nil
	}))

	for _, pair := range [][2]string{
		{
			value.GetIdentity().GetIssuer(),
			"another-subject",
		},
		{
			"https://other.flida.test",
			value.GetIdentity().GetSubject(),
		},
	} {
		identity, err := model_identity.NewIdentity(pair[0], pair[1])
		x.Require().NoError(err)
		got, err := x.repository.FindByIdentity(ctx, identity)
		x.Require().ErrorIs(err, aggregate_volunteer.ErrNotFound)
		x.Require().Nil(got)
	}
}
