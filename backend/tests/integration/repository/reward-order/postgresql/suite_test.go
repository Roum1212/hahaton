//go:build integration

package repository_reward_order_postgresql_test

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	aggregate_reward_order "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order"
	aggregate_reward_order_command "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order/command"
	client_postgresql "github.com/Roum1212/hahaton/internal/infra/client/postgresql"
	repository_reward_order_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/reward-order/postgresql"
	pkg_testcontainers_postgresql "github.com/Roum1212/hahaton/pkg/testcontainers/postgresql"
	repository_fixture "github.com/Roum1212/hahaton/tests/integration/repository/fixture"
)

type RepositoryTestSuite struct {
	suite.Suite

	pool       *pgxpool.Pool
	database   client_postgresql.Client
	repository aggregate_reward_order.Repository
}

func (x *RepositoryTestSuite) SetupSuite() {
	x.pool = pkg_testcontainers_postgresql.NewPostgreSQL(x.T())
	x.database = client_postgresql.NewClient(x.pool)
	x.repository = repository_reward_order_postgresql.NewRewardOrders(x.database)
}

func (x *RepositoryTestSuite) SetupTest() {
	repository_fixture.Reset(x.T(), x.pool)
}

func (x *RepositoryTestSuite) changeValue(value *aggregate_reward_order.RewardOrder) {
	x.Require().NoError(value.Issue(aggregate_reward_order_command.NewIssue(repository_fixture.Now().Add(time.Minute))))
}

func (x *RepositoryTestSuite) newValue() *aggregate_reward_order.RewardOrder {
	return repository_fixture.NewOrder(
		x.T(),
		repository_fixture.CreateReward(x.T(), x.database),
		repository_fixture.CreateVolunteer(x.T(), x.database).GetUserID(),
	)
}

func TestRepositoryTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RepositoryTestSuite))
}

func assertRewardOrder(t *testing.T, expected, actual *aggregate_reward_order.RewardOrder) {
	t.Helper()

	require.NotNil(t, actual)
	require.Equal(t, expected.GetCancellationReason(), actual.GetCancellationReason())
	require.Equal(t, expected.GetCoordinatorID(), actual.GetCoordinatorID())
	require.WithinDuration(t, expected.GetCreatedAt(), actual.GetCreatedAt(), 0)
	require.Equal(t, expected.GetPriceCoins(), actual.GetPriceCoins())
	require.Equal(t, expected.GetRequestID(), actual.GetRequestID())
	require.Equal(t, expected.GetRewardID(), actual.GetRewardID())
	require.Equal(t, expected.GetRewardOrderID(), actual.GetRewardOrderID())
	require.Equal(t, expected.GetRewardTitle(), actual.GetRewardTitle())
	require.Equal(t, expected.GetStatus(), actual.GetStatus())
	require.WithinDuration(t, expected.GetUpdatedAt(), actual.GetUpdatedAt(), 0)
	require.Equal(t, expected.GetVolunteerID(), actual.GetVolunteerID())
}
