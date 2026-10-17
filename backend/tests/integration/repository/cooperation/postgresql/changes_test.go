//go:build integration

package repository_cooperation_postgresql_test

import (
	"context"
	"errors"
	"time"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	model_cooperation_change "github.com/Roum1212/hahaton/internal/domain/model/cooperation-change"
	repository_fixture "github.com/Roum1212/hahaton/tests/integration/repository/fixture"
)

func (x *RepositoryTestSuite) TestRepository_AppendChange() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value, func(ctx context.Context) error {
		return nil
	}))
	change, err := model_cooperation_change.NewCooperationChange(
		value.GetCoordinatorID(), value.GetVolunteerID(), value.GetStatus(), value.GetVolunteerID(),
		time.Time{}, repository_fixture.Now(),
	)
	x.Require().NoError(err)
	x.Require().NoError(x.repository.AppendChange(ctx, change))

	var (
		coordinator, volunteer, status, actor string
		retryAfter                            *time.Time
		createdAt                             time.Time
	)

	err = x.pool.QueryRow(ctx, `SELECT coordinator_id, volunteer_id, status, actor_id, retry_after, created_at
FROM cooperation_changes WHERE coordinator_id = $1 AND volunteer_id = $2`,
		value.GetCoordinatorID().String(), value.GetVolunteerID().String(),
	).Scan(&coordinator, &volunteer, &status, &actor, &retryAfter, &createdAt)
	x.Require().NoError(err)
	x.Require().Equal(change.GetCoordinatorID().String(), coordinator)
	x.Require().Equal(change.GetVolunteerID().String(), volunteer)
	x.Require().Equal(change.GetStatus().String(), status)
	x.Require().Equal(change.GetActorID().String(), actor)
	x.Require().Nil(retryAfter)
	x.Require().WithinDuration(change.GetCreatedAt(), createdAt, 0)
}

func (x *RepositoryTestSuite) TestRepository_AppendChange_Rollback() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value, func(ctx context.Context) error {
		return nil
	}))

	stop := errors.New("cancel cooperation journal")
	err := x.repository.Modify(ctx, value.GetCoordinatorID(), value.GetVolunteerID(), func(
		ctx context.Context,
		current *aggregate_cooperation.Cooperation,
	) error {
		x.changeValue(current)

		change, err := model_cooperation_change.NewCooperationChange(
			current.GetCoordinatorID(), current.GetVolunteerID(), current.GetStatus(), current.GetCoordinatorID(),
			current.GetRetryAfter(), repository_fixture.Now(),
		)
		if err != nil {
			return err
		}

		if err := x.repository.AppendChange(ctx, change); err != nil {
			return err
		}

		return stop
	})
	x.Require().ErrorIs(err, stop)

	var count int

	x.Require().NoError(x.pool.QueryRow(ctx, "SELECT count(*) FROM cooperation_changes").Scan(&count))
	x.Require().Zero(count)
	got, err := x.repository.Load(ctx, value.GetCoordinatorID(), value.GetVolunteerID())
	x.Require().NoError(err)
	assertCooperation(x.T(), value, got)
}
