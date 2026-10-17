//go:build integration

package repository_wallet_postgresql_test

import (
	"context"
	"errors"
	"time"

	aggregate_wallet "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet"
	repository_fixture "github.com/Roum1212/hahaton/tests/integration/repository/fixture"
)

func (x *RepositoryTestSuite) TestRepository_CreateAndLoad() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))
	got, err := x.repository.Load(ctx, value.GetUserID())
	x.Require().NoError(err)
	assertWallet(x.T(), value, got)
	x.Require().Error(x.repository.Create(ctx, value))
}

func (x *RepositoryTestSuite) TestRepository_Load_NotFound() {
	got, err := x.repository.Load(x.T().Context(), x.newValue().GetUserID())
	x.Require().ErrorIs(err, aggregate_wallet.ErrNotFound)
	x.Require().Nil(got)
}

func (x *RepositoryTestSuite) TestRepository_Modify_NotFound() {
	err := x.repository.Modify(x.T().Context(), x.newValue().GetUserID(), func(
		ctx context.Context,
		value *aggregate_wallet.Wallet,
	) error {
		return nil
	})
	x.Require().ErrorIs(err, aggregate_wallet.ErrNotFound)
}

func (x *RepositoryTestSuite) TestRepository_Modify_Rollback() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))

	stop := errors.New("cancel wallet change")
	err := x.repository.Modify(ctx, value.GetUserID(), func(
		ctx context.Context,
		current *aggregate_wallet.Wallet,
	) error {
		updated := aggregate_wallet.NewWalletFromFields(
			current.GetUserID(), 30, 180, 1, current.GetCreatedAt(), repository_fixture.Now().Add(time.Minute),
		)
		if err := x.repository.Save(ctx, updated); err != nil {
			return err
		}

		return stop
	})
	x.Require().ErrorIs(err, stop)
	got, err := x.repository.Load(ctx, value.GetUserID())
	x.Require().NoError(err)
	assertWallet(x.T(), value, got)
}

func (x *RepositoryTestSuite) TestRepository_Save() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))
	updated := aggregate_wallet.NewWalletFromFields(
		value.GetUserID(), 30, 180, 1, value.GetCreatedAt(), repository_fixture.Now().Add(time.Minute),
	)
	x.Require().NoError(x.repository.Save(ctx, updated))
	got, err := x.repository.Load(ctx, value.GetUserID())
	x.Require().NoError(err)
	assertWallet(x.T(), updated, got)
}

func (x *RepositoryTestSuite) TestRepository_Save_NotFound() {
	x.Require().ErrorIs(x.repository.Save(x.T().Context(), x.newValue()), aggregate_wallet.ErrNotFound)
}
