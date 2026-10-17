//go:build integration

package repository_cooperation_invitation_postgresql_test

import (
	aggregate_cooperation_invitation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation-invitation"
)

func (x *RepositoryTestSuite) TestRepository_Load_NotFound() {
	value := x.newValue()
	got, err := x.repository.Load(x.T().Context(), value.GetInvitationID())
	x.Require().ErrorIs(err, aggregate_cooperation_invitation.ErrNotFound)
	x.Require().Nil(got)
}
