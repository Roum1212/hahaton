//go:build integration

package repository_reward_postgresql_test

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	aggregate_reward "github.com/Roum1212/hahaton/internal/domain/aggregate/reward"
	aggregate_reward_command "github.com/Roum1212/hahaton/internal/domain/aggregate/reward/command"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	client_postgresql "github.com/Roum1212/hahaton/internal/infra/client/postgresql"
	repository_reward_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/reward/postgresql"
	pkg_testcontainers_postgresql "github.com/Roum1212/hahaton/pkg/testcontainers/postgresql"
	repository_fixture "github.com/Roum1212/hahaton/tests/integration/repository/fixture"
)

type RepositoryTestSuite struct {
	suite.Suite

	pool       *pgxpool.Pool
	database   client_postgresql.Client
	repository aggregate_reward.Repository
}

func (x *RepositoryTestSuite) SetupSuite() {
	x.pool = pkg_testcontainers_postgresql.NewPostgreSQL(x.T())
	x.database = client_postgresql.NewClient(x.pool)
	x.repository = repository_reward_postgresql.NewRewards(x.database)
}

func (x *RepositoryTestSuite) SetupTest() {
	repository_fixture.Reset(x.T(), x.pool)
}

func (x *RepositoryTestSuite) changeValue(value *aggregate_reward.Reward) {
	x.Require().NoError(
		value.Update(
			aggregate_reward_command.NewUpdate(
				"Новая награда",
				"Новое описание",
				model_image_path.ImagePath{},
				25,
				5,
				repository_fixture.Now().Add(time.Minute),
			),
		),
	)
}

func (x *RepositoryTestSuite) newValue() *aggregate_reward.Reward {
	return repository_fixture.NewReward(x.T(), repository_fixture.CreateCoordinator(x.T(), x.database).GetUserID())
}

func TestRepositoryTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RepositoryTestSuite))
}

func assertReward(t *testing.T, expected, actual *aggregate_reward.Reward) {
	t.Helper()

	require.NotNil(t, actual)
	require.Equal(t, expected.GetAvailableQuantity(), actual.GetAvailableQuantity())
	require.Equal(t, expected.GetCoordinatorID(), actual.GetCoordinatorID())
	require.WithinDuration(t, expected.GetCreatedAt(), actual.GetCreatedAt(), 0)
	require.Equal(t, expected.GetDeleted(), actual.GetDeleted())
	require.Equal(t, expected.GetDescription(), actual.GetDescription())
	require.Equal(t, expected.GetImagePath(), actual.GetImagePath())
	require.Equal(t, expected.GetIssuedQuantity(), actual.GetIssuedQuantity())
	require.Equal(t, expected.GetPriceCoins(), actual.GetPriceCoins())
	require.Equal(t, expected.GetQuantityLimit(), actual.GetQuantityLimit())
	require.Equal(t, expected.GetReservedQuantity(), actual.GetReservedQuantity())
	require.Equal(t, expected.GetRewardID(), actual.GetRewardID())
	require.Equal(t, expected.GetTitle(), actual.GetTitle())
	require.WithinDuration(t, expected.GetUpdatedAt(), actual.GetUpdatedAt(), 0)
}
