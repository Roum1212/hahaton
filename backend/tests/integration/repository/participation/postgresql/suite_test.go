//go:build integration

package repository_participation_postgresql_test

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	client_postgresql "github.com/Roum1212/hahaton/internal/infra/client/postgresql"
	repository_participation_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/participation/postgresql"
	pkg_testcontainers_postgresql "github.com/Roum1212/hahaton/pkg/testcontainers/postgresql"
	repository_fixture "github.com/Roum1212/hahaton/tests/integration/repository/fixture"
)

type RepositoryTestSuite struct {
	suite.Suite

	pool       *pgxpool.Pool
	database   client_postgresql.Client
	repository aggregate_participation.Repository
}

func (x *RepositoryTestSuite) SetupSuite() {
	x.pool = pkg_testcontainers_postgresql.NewPostgreSQL(x.T())
	x.database = client_postgresql.NewClient(x.pool)
	x.repository = repository_participation_postgresql.NewParticipations(x.database)
}

func (x *RepositoryTestSuite) SetupTest() {
	repository_fixture.Reset(x.T(), x.pool)
}

func (x *RepositoryTestSuite) newValue() *aggregate_participation.Participation {
	return repository_fixture.NewParticipation(x.T(), repository_fixture.CreateEvent(
		x.T(),
		x.database,
	), repository_fixture.CreateVolunteer(x.T(), x.database).GetUserID())
}

func TestRepositoryTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RepositoryTestSuite))
}

func assertParticipation(t *testing.T, expected, actual *aggregate_participation.Participation) {
	t.Helper()

	require.NotNil(t, actual)
	require.Equal(t, expected.GetAccrualOperationID(), actual.GetAccrualOperationID())
	require.Equal(t, expected.GetAwardedCoins(), actual.GetAwardedCoins())
	require.Equal(t, expected.GetAwardedMinutes(), actual.GetAwardedMinutes())
	require.Equal(t, expected.GetCancellationReason(), actual.GetCancellationReason())
	require.Equal(t, expected.GetCoordinatorID(), actual.GetCoordinatorID())
	require.WithinDuration(t, expected.GetCreatedAt(), actual.GetCreatedAt(), 0)
	require.Equal(t, expected.GetEventID(), actual.GetEventID())
	require.Equal(t, expected.GetHasChatAccess(), actual.GetHasChatAccess())
	require.Equal(t, expected.GetParticipationID(), actual.GetParticipationID())
	require.Equal(t, expected.GetRejectionReason(), actual.GetRejectionReason())
	require.Equal(t, expected.GetResult(), actual.GetResult())
	require.Equal(t, expected.GetStatus(), actual.GetStatus())
	require.WithinDuration(t, expected.GetUpdatedAt(), actual.GetUpdatedAt(), 0)
	require.Equal(t, expected.GetVisibleAfterSequence(), actual.GetVisibleAfterSequence())
	require.Equal(t, expected.GetVolunteerID(), actual.GetVolunteerID())
}
