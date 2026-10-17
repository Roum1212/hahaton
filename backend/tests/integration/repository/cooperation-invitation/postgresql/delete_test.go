//go:build integration

package repository_cooperation_invitation_postgresql_test

import (
	aggregate_cooperation_invitation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation-invitation"
)

func (x *RepositoryTestSuite) TestRepository_Delete() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))
	x.Require().NoError(x.repository.Delete(ctx, value.GetInvitationID()))
	got, err := x.repository.Load(ctx, value.GetInvitationID())
	x.Require().ErrorIs(err, aggregate_cooperation_invitation.ErrNotFound)
	x.Require().Nil(got)
	x.Require().NoError(x.repository.Delete(ctx, value.GetInvitationID()))
}

func (x *RepositoryTestSuite) TestRepository_TokenEncrypted() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))

	var stored string

	err := x.pool.QueryRow(
		ctx,
		"SELECT token FROM cooperation_invitations WHERE id = $1",
		value.GetInvitationID().String(),
	).
		Scan(&stored)
	x.Require().NoError(err)
	x.Require().NotEmpty(stored)
	x.Require().NotEqual(value.GetToken(), stored)
	got, err := x.repository.Load(ctx, value.GetInvitationID())
	x.Require().NoError(err)
	x.Require().Equal(value.GetToken(), got.GetToken())
}
