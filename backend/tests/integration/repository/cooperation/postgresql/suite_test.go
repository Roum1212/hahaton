//go:build integration

package repository_cooperation_postgresql_test

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_cooperation_command "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation/command"
	client_postgresql "github.com/Roum1212/hahaton/internal/infra/client/postgresql"
	repository_cooperation_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/cooperation/postgresql"
	pkg_testcontainers_postgresql "github.com/Roum1212/hahaton/pkg/testcontainers/postgresql"
	repository_fixture "github.com/Roum1212/hahaton/tests/integration/repository/fixture"
)

type RepositoryTestSuite struct {
	suite.Suite

	pool       *pgxpool.Pool
	database   client_postgresql.Client
	repository aggregate_cooperation.Repository
}

func (x *RepositoryTestSuite) SetupSuite() {
	x.pool = pkg_testcontainers_postgresql.NewPostgreSQL(x.T())
	x.database = client_postgresql.NewClient(x.pool)
	x.repository = repository_cooperation_postgresql.NewCooperations(x.database)
}

func (x *RepositoryTestSuite) SetupTest() {
	repository_fixture.Reset(x.T(), x.pool)
}

func (x *RepositoryTestSuite) changeValue(value *aggregate_cooperation.Cooperation) {
	x.Require().NoError(
		value.Reject(
			aggregate_cooperation_command.NewReject(repository_fixture.Now().Add(
				time.Hour,
			), repository_fixture.Now().Add(time.Minute)),
		),
	)
}

func (x *RepositoryTestSuite) newValue() *aggregate_cooperation.Cooperation {
	return repository_fixture.NewCooperation(
		x.T(),
		repository_fixture.CreateCoordinator(x.T(), x.database).GetUserID(),
		repository_fixture.CreateVolunteer(x.T(), x.database).GetUserID(),
	)
}

func TestRepositoryTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RepositoryTestSuite))
}

func assertCooperation(t *testing.T, expected, actual *aggregate_cooperation.Cooperation) {
	t.Helper()

	require.NotNil(t, actual)
	require.Equal(t, expected.GetCoordinatorID(), actual.GetCoordinatorID())
	require.WithinDuration(t, expected.GetCreatedAt(), actual.GetCreatedAt(), 0)
	require.WithinDuration(t, expected.GetRetryAfter(), actual.GetRetryAfter(), 0)
	require.Equal(t, expected.GetStatus(), actual.GetStatus())
	require.WithinDuration(t, expected.GetUpdatedAt(), actual.GetUpdatedAt(), 0)
	require.Equal(t, expected.GetVolunteerID(), actual.GetVolunteerID())
}
