//go:build integration

package repository_reward_postgresql_test

import (
	"context"
	"errors"

	aggregate_reward "github.com/Roum1212/hahaton/internal/domain/aggregate/reward"
)

func (x *RepositoryTestSuite) TestRepository_Modify() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))
	err := x.repository.Modify(
		ctx,
		value.GetRewardID(),
		func(ctx context.Context, current *aggregate_reward.Reward) error {
			x.changeValue(current)

			return nil
		},
	)
	x.Require().NoError(err)
	x.changeValue(value)
	got, err := x.repository.Load(ctx, value.GetRewardID())
	x.Require().NoError(err)
	assertReward(x.T(), value, got)
}

func (x *RepositoryTestSuite) TestRepository_Modify_NotFound() {
	ctx := x.T().Context()
	value := x.newValue()
	err := x.repository.Modify(
		ctx,
		value.GetRewardID(),
		func(ctx context.Context, current *aggregate_reward.Reward) error {
			return nil
		},
	)
	x.Require().ErrorIs(err, aggregate_reward.ErrNotFound)
}

func (x *RepositoryTestSuite) TestRepository_Modify_Rollback() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))

	stop := errors.New("cancel modification")
	err := x.repository.Modify(
		ctx,
		value.GetRewardID(),
		func(ctx context.Context, current *aggregate_reward.Reward) error {
			x.changeValue(current)

			if err := x.repository.Save(ctx, current); err != nil {
				return err
			}

			return stop
		},
	)
	x.Require().ErrorIs(err, stop)
	got, err := x.repository.Load(ctx, value.GetRewardID())
	x.Require().NoError(err)
	assertReward(x.T(), value, got)
}
