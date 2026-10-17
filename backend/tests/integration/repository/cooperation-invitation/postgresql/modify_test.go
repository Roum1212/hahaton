//go:build integration

package repository_cooperation_invitation_postgresql_test

import (
	"context"
	"errors"

	aggregate_cooperation_invitation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation-invitation"
)

func (x *RepositoryTestSuite) TestRepository_Modify() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))
	err := x.repository.Modify(
		ctx,
		value.GetInvitationID(),
		func(ctx context.Context, current *aggregate_cooperation_invitation.CooperationInvitation) error {
			x.changeValue(current)

			return nil
		},
	)
	x.Require().NoError(err)
	x.changeValue(value)
	got, err := x.repository.Load(ctx, value.GetInvitationID())
	x.Require().NoError(err)
	assertCooperationInvitation(x.T(), value, got)
}

func (x *RepositoryTestSuite) TestRepository_Modify_NotFound() {
	ctx := x.T().Context()
	value := x.newValue()
	err := x.repository.Modify(
		ctx,
		value.GetInvitationID(),
		func(ctx context.Context, current *aggregate_cooperation_invitation.CooperationInvitation) error {
			return nil
		},
	)
	x.Require().ErrorIs(err, aggregate_cooperation_invitation.ErrNotFound)
}

func (x *RepositoryTestSuite) TestRepository_Modify_Rollback() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))

	stop := errors.New("cancel modification")
	err := x.repository.Modify(
		ctx,
		value.GetInvitationID(),
		func(ctx context.Context, current *aggregate_cooperation_invitation.CooperationInvitation) error {
			x.changeValue(current)

			if err := x.repository.Save(ctx, current); err != nil {
				return err
			}

			return stop
		},
	)
	x.Require().ErrorIs(err, stop)
	got, err := x.repository.Load(ctx, value.GetInvitationID())
	x.Require().NoError(err)
	assertCooperationInvitation(x.T(), value, got)
}
