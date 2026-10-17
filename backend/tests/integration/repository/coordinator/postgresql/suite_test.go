//go:build integration

package repository_coordinator_postgresql_test

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_coordinator_command "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator/command"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	client_postgresql "github.com/Roum1212/hahaton/internal/infra/client/postgresql"
	repository_coordinator_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/coordinator/postgresql"
	pkg_testcontainers_postgresql "github.com/Roum1212/hahaton/pkg/testcontainers/postgresql"
	repository_fixture "github.com/Roum1212/hahaton/tests/integration/repository/fixture"
)

type RepositoryTestSuite struct {
	suite.Suite

	pool       *pgxpool.Pool
	database   client_postgresql.Client
	repository repository_coordinator_postgresql.Coordinators
}

func (x *RepositoryTestSuite) SetupSuite() {
	x.pool = pkg_testcontainers_postgresql.NewPostgreSQL(x.T())
	x.database = client_postgresql.NewClient(x.pool)
	x.repository = repository_coordinator_postgresql.NewCoordinators(x.database)
}

func (x *RepositoryTestSuite) SetupTest() {
	repository_fixture.Reset(x.T(), x.pool)
}

func (x *RepositoryTestSuite) changeValue(value *aggregate_coordinator.Coordinator) {
	value.Update(
		aggregate_coordinator_command.NewUpdate(
			"Обновлено",
			model_image_path.ImagePath{},
			120,
			repository_fixture.Now().Add(time.Minute),
		),
	)
}

func (x *RepositoryTestSuite) newValue() *aggregate_coordinator.Coordinator {
	return repository_fixture.NewCoordinator(x.T(), repository_fixture.CreateVolunteer(x.T(), x.database).GetUserID())
}

func TestRepositoryTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RepositoryTestSuite))
}

func assertCoordinator(t *testing.T, expected, actual *aggregate_coordinator.Coordinator) {
	t.Helper()

	require.NotNil(t, actual)
	require.Equal(t, expected.GetAvatarPath(), actual.GetAvatarPath())
	require.Equal(t, expected.GetCooldownSeconds(), actual.GetCooldownSeconds())
	require.WithinDuration(t, expected.GetCreatedAt(), actual.GetCreatedAt(), 0)
	require.Equal(t, expected.GetDescription(), actual.GetDescription())
	require.WithinDuration(t, expected.GetUpdatedAt(), actual.GetUpdatedAt(), 0)
	require.Equal(t, expected.GetUserID(), actual.GetUserID())
}
