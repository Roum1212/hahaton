//go:build integration

package repository_participation_postgresql_test

import (
	"time"

	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	aggregate_participation_command "github.com/Roum1212/hahaton/internal/domain/aggregate/participation/command"
	model_event_result "github.com/Roum1212/hahaton/internal/domain/model/event-result"
	model_operation_id "github.com/Roum1212/hahaton/internal/domain/model/operation-id"
	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
	repository_fixture "github.com/Roum1212/hahaton/tests/integration/repository/fixture"
)

func (x *RepositoryTestSuite) TestRepository_CreateAndLoad() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))
	got, err := x.repository.Load(ctx, value.GetParticipationID())
	x.Require().NoError(err)
	assertParticipation(x.T(), value, got)
}

func (x *RepositoryTestSuite) TestRepository_Create_DuplicatePair() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))

	id, err := model_participation_id.GenerateParticipationID()
	x.Require().NoError(err)

	duplicate := aggregate_participation.NewParticipation(
		id,
		value.GetEventID(),
		value.GetCoordinatorID(),
		value.GetVolunteerID(),
	)
	x.Require().NoError(duplicate.Request(aggregate_participation_command.NewRequest(repository_fixture.Now())))
	x.Require().ErrorIs(x.repository.Create(ctx, duplicate), aggregate_participation.ErrAlreadyExists)
}

func (x *RepositoryTestSuite) TestRepository_List() {
	ctx := x.T().Context()
	values, err := x.repository.List(ctx)
	x.Require().NoError(err)
	x.Require().Empty(values)
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))
	values, err = x.repository.List(ctx)
	x.Require().NoError(err)
	x.Require().Len(values, 1)
	assertParticipation(x.T(), value, values[0])
}

func (x *RepositoryTestSuite) TestRepository_Load_NotFound() {
	got, err := x.repository.Load(x.T().Context(), x.newValue().GetParticipationID())
	x.Require().ErrorIs(err, aggregate_participation.ErrNotFound)
	x.Require().Nil(got)
}

func (x *RepositoryTestSuite) TestRepository_Save_Absent() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))
	x.Require().NoError(value.Approve(aggregate_participation_command.NewApprove(0, repository_fixture.Now())))
	result, err := model_event_result.NewEventResult(value.GetParticipationID(), 0, true)
	x.Require().NoError(err)
	x.Require().NoError(value.Complete(aggregate_participation_command.NewComplete(
		result,
		0,
		model_operation_id.OperationID{},
		repository_fixture.Now(),
	)))
	x.Require().NoError(x.repository.Save(ctx, value))
	got, err := x.repository.Load(ctx, value.GetParticipationID())
	x.Require().NoError(err)
	assertParticipation(x.T(), value, got)
}

func (x *RepositoryTestSuite) TestRepository_Save_Confirmed() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))
	x.Require().NoError(value.Approve(aggregate_participation_command.NewApprove(7, repository_fixture.Now())))
	result, err := model_event_result.NewEventResult(value.GetParticipationID(), 180, false)
	x.Require().NoError(err)
	operationID, err := model_operation_id.GenerateOperationID()
	x.Require().NoError(err)
	x.Require().NoError(value.Complete(aggregate_participation_command.NewComplete(
		result,
		30,
		operationID,
		repository_fixture.Now().Add(time.Minute),
	)))
	x.Require().NoError(x.repository.Save(ctx, value))
	got, err := x.repository.Load(ctx, value.GetParticipationID())
	x.Require().NoError(err)
	assertParticipation(x.T(), value, got)
}

func (x *RepositoryTestSuite) TestRepository_Save_NotFound() {
	x.Require().ErrorIs(x.repository.Save(x.T().Context(), x.newValue()), aggregate_participation.ErrNotFound)
}
