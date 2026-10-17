//go:build integration

package repository_event_postgresql_test

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_event_command "github.com/Roum1212/hahaton/internal/domain/aggregate/event/command"
	client_postgresql "github.com/Roum1212/hahaton/internal/infra/client/postgresql"
	repository_event_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/event/postgresql"
	pkg_testcontainers_postgresql "github.com/Roum1212/hahaton/pkg/testcontainers/postgresql"
	repository_fixture "github.com/Roum1212/hahaton/tests/integration/repository/fixture"
)

type RepositoryTestSuite struct {
	suite.Suite

	pool       *pgxpool.Pool
	database   client_postgresql.Client
	repository aggregate_event.Repository
}

func (x *RepositoryTestSuite) SetupSuite() {
	x.pool = pkg_testcontainers_postgresql.NewPostgreSQL(x.T())
	x.database = client_postgresql.NewClient(x.pool)
	x.repository = repository_event_postgresql.NewEvents(x.database)
}

func (x *RepositoryTestSuite) SetupTest() {
	repository_fixture.Reset(x.T(), x.pool)
}

func (x *RepositoryTestSuite) changeValue(value *aggregate_event.Event) {
	x.Require().NoError(value.Cancel(aggregate_event_command.NewCancel(repository_fixture.Now().Add(time.Minute))))
}

func (x *RepositoryTestSuite) newValue() *aggregate_event.Event {
	return repository_fixture.NewEvent(x.T(), repository_fixture.CreateCoordinator(x.T(), x.database).GetUserID())
}

func TestRepositoryTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RepositoryTestSuite))
}

func assertEvent(t *testing.T, expected, actual *aggregate_event.Event) {
	t.Helper()

	require.NotNil(t, actual)
	require.Equal(t, expected.GetCompletionRequestID(), actual.GetCompletionRequestID())
	require.Equal(t, expected.GetCoordinatorID(), actual.GetCoordinatorID())
	require.WithinDuration(t, expected.GetCreatedAt(), actual.GetCreatedAt(), 0)
	require.Equal(t, expected.GetDescription(), actual.GetDescription())
	require.Equal(t, expected.GetEventID(), actual.GetEventID())
	require.Equal(t, expected.GetImagePath(), actual.GetImagePath())
	require.Equal(t, expected.GetLocation(), actual.GetLocation())
	require.Equal(t, expected.GetPoint(), actual.GetPoint())
	require.ElementsMatch(t, expected.GetRequirements().GetValues(), actual.GetRequirements().GetValues())
	require.Equal(t, expected.GetResultFingerprint(), actual.GetResultFingerprint())
	require.Equal(t, expected.GetSchedule(), actual.GetSchedule())
	require.Equal(t, expected.GetStatus(), actual.GetStatus())
	require.Equal(t, expected.GetTitle(), actual.GetTitle())
	require.WithinDuration(t, expected.GetUpdatedAt(), actual.GetUpdatedAt(), 0)
	require.Equal(t, expected.GetVolunteerLimit(), actual.GetVolunteerLimit())
}
