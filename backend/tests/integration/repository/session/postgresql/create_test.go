//go:build integration

package repository_session_postgresql_test

import (
	"time"

	aggregate_session "github.com/Roum1212/hahaton/internal/domain/aggregate/session"
	aggregate_session_command "github.com/Roum1212/hahaton/internal/domain/aggregate/session/command"
	model_session_id "github.com/Roum1212/hahaton/internal/domain/model/session-id"
)

func (x *RepositoryTestSuite) TestRepository_CreateAndLoad() {
	id, err := model_session_id.GenerateSessionID()
	x.Require().NoError(err)

	now := time.Now().Truncate(time.Microsecond)
	session := aggregate_session.NewSession(id)
	x.Require().NoError(session.Open(aggregate_session_command.NewOpen(now, now.Add(time.Hour))))
	x.Require().NoError(x.repository.Create(x.T().Context(), session))

	got, err := x.repository.Load(x.T().Context(), id)
	x.Require().NoError(err)
	x.Require().Equal(session.GetSessionID(), got.GetSessionID())
	x.Require().Equal(session.GetUserID(), got.GetUserID())
	x.Require().True(got.IsAnonymous())
	x.Require().WithinDuration(session.GetCreatedAt(), got.GetCreatedAt(), 0)
	x.Require().WithinDuration(session.GetUpdatedAt(), got.GetUpdatedAt(), 0)
	x.Require().WithinDuration(session.GetExpiresAt(), got.GetExpiresAt(), 0)
}
