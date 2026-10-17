//go:build integration

package repository_participation_postgresql_test

import (
	model_actor_kind "github.com/Roum1212/hahaton/internal/domain/model/actor-kind"
	model_operation_id "github.com/Roum1212/hahaton/internal/domain/model/operation-id"
	model_participation_change "github.com/Roum1212/hahaton/internal/domain/model/participation-change"
	model_participation_status "github.com/Roum1212/hahaton/internal/domain/model/participation-status"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
	repository_fixture "github.com/Roum1212/hahaton/tests/integration/repository/fixture"
)

func (x *RepositoryTestSuite) TestRepository_Changes() {
	ctx := x.T().Context()
	changes, err := x.repository.ListChanges(ctx)
	x.Require().NoError(err)
	x.Require().Empty(changes)
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))

	expectedSequences := []uint64{1, 2, 3}

	for index, actor := range []struct {
		kind         model_actor_kind.ActorKind
		id           model_user_id.UserID
		previous     utils_value.Optional[model_participation_status.ParticipationStatus]
		status       model_participation_status.ParticipationStatus
		cancellation string
		rejection    string
	}{
		{
			model_actor_kind.Volunteer,
			value.GetVolunteerID(),
			utils_value.Optional[model_participation_status.ParticipationStatus]{},
			model_participation_status.Pending,
			"",
			"",
		},
		{
			model_actor_kind.Coordinator,
			value.GetCoordinatorID(),
			utils_value.NewOptional(model_participation_status.Pending),
			model_participation_status.Rejected,
			"",
			"Причина",
		},
		{
			model_actor_kind.System,
			model_user_id.UserID{},
			utils_value.NewOptional(model_participation_status.Pending),
			model_participation_status.Canceled,
			"event_canceled",
			"",
		},
	} {
		id, err := model_operation_id.GenerateOperationID()
		x.Require().NoError(err)
		change, err := model_participation_change.NewParticipationChange(
			id, value.GetParticipationID(), actor.previous, actor.status, actor.id, actor.kind,
			actor.cancellation, actor.rejection, 0, repository_fixture.Now(),
		)
		x.Require().NoError(err)
		x.Require().NoError(x.repository.AppendChange(ctx, change))
		changes, err = x.repository.ListChanges(ctx)
		x.Require().NoError(err)
		x.Require().Len(changes, index+1)

		found := false

		for _, got := range changes {
			if got.GetChangeID() != id {
				continue
			}

			found = true

			x.Require().Equal(change.GetParticipationID(), got.GetParticipationID())
			x.Require().Equal(change.GetPreviousStatus(), got.GetPreviousStatus())
			x.Require().Equal(change.GetStatus(), got.GetStatus())
			x.Require().Equal(change.GetActorID(), got.GetActorID())
			x.Require().Equal(change.GetActorKind(), got.GetActorKind())
			x.Require().Equal(change.GetCancellationReason(), got.GetCancellationReason())
			x.Require().Equal(change.GetRejectionReason(), got.GetRejectionReason())
			x.Require().Equal(expectedSequences[index], got.GetSequence())
			x.Require().WithinDuration(change.GetCreatedAt(), got.GetCreatedAt(), 0)
		}

		x.Require().True(found)
	}
}
