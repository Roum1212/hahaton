//go:build integration

package repository_cooperation_invitation_postgresql_test

import (
	aggregate_cooperation_invitation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation-invitation"
)

func (x *RepositoryTestSuite) TestRepository_List() {
	ctx := x.T().Context()
	values, err := x.repository.List(ctx)
	x.Require().NoError(err)
	x.Require().Empty(values)
	first, second := x.newValue(), x.newValue()
	x.Require().NoError(x.repository.Create(ctx, first))
	x.Require().NoError(x.repository.Create(ctx, second))
	values, err = x.repository.List(ctx)
	x.Require().NoError(err)
	x.Require().Len(values, 2)
	x.Require().ElementsMatch([]*aggregate_cooperation_invitation.CooperationInvitation{first, second}, values)
}
