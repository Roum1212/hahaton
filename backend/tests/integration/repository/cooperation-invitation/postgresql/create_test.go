//go:build integration

package repository_cooperation_invitation_postgresql_test

import (
	aggregate_cooperation_invitation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation-invitation"
)

func (x *RepositoryTestSuite) TestRepository_CreateAndLoad() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))
	got, err := x.repository.Load(ctx, value.GetInvitationID())
	x.Require().NoError(err)
	assertCooperationInvitation(x.T(), value, got)
}

func (x *RepositoryTestSuite) TestRepository_Create_Duplicate() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))
	err := x.repository.Create(ctx, value)
	x.Require().ErrorIs(err, aggregate_cooperation_invitation.ErrAlreadyExists)
	got, err := x.repository.Load(ctx, value.GetInvitationID())
	x.Require().NoError(err)
	assertCooperationInvitation(x.T(), value, got)
}
