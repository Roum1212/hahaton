//go:build integration

package repository_wallet_postgresql_test

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	aggregate_wallet "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet"
	client_postgresql "github.com/Roum1212/hahaton/internal/infra/client/postgresql"
	repository_wallet_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/wallet/postgresql"
	pkg_testcontainers_postgresql "github.com/Roum1212/hahaton/pkg/testcontainers/postgresql"
	repository_fixture "github.com/Roum1212/hahaton/tests/integration/repository/fixture"
)

type RepositoryTestSuite struct {
	suite.Suite

	pool       *pgxpool.Pool
	database   client_postgresql.Client
	repository aggregate_wallet.Repository
}

func (x *RepositoryTestSuite) SetupSuite() {
	x.pool = pkg_testcontainers_postgresql.NewPostgreSQL(x.T())
	x.database = client_postgresql.NewClient(x.pool)
	x.repository = repository_wallet_postgresql.NewWallets(x.database)
}

func (x *RepositoryTestSuite) SetupTest() {
	repository_fixture.Reset(x.T(), x.pool)
}

func (x *RepositoryTestSuite) newValue() *aggregate_wallet.Wallet {
	return repository_fixture.NewWallet(x.T(), repository_fixture.CreateVolunteer(x.T(), x.database).GetUserID())
}

func TestRepositoryTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RepositoryTestSuite))
}

func assertWallet(t *testing.T, expected, actual *aggregate_wallet.Wallet) {
	t.Helper()

	require.NotNil(t, actual)
	require.Equal(t, expected.GetCoinBalance(), actual.GetCoinBalance())
	require.Equal(t, expected.GetConfirmedMinutes(), actual.GetConfirmedMinutes())
	require.WithinDuration(t, expected.GetCreatedAt(), actual.GetCreatedAt(), 0)
	require.Equal(t, expected.GetLastSequence(), actual.GetLastSequence())
	require.WithinDuration(t, expected.GetUpdatedAt(), actual.GetUpdatedAt(), 0)
	require.Equal(t, expected.GetUserID(), actual.GetUserID())
}
