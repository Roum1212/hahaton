//go:build integration

package repository_session_postgresql_test

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	aggregate_session "github.com/Roum1212/hahaton/internal/domain/aggregate/session"
	client_postgresql "github.com/Roum1212/hahaton/internal/infra/client/postgresql"
	repository_session_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/session/postgresql"
	pkg_testcontainers_postgresql "github.com/Roum1212/hahaton/pkg/testcontainers/postgresql"
	repository_fixture "github.com/Roum1212/hahaton/tests/integration/repository/fixture"
)

type RepositoryTestSuite struct {
	suite.Suite

	pool       *pgxpool.Pool
	repository aggregate_session.Repository
}

func (x *RepositoryTestSuite) SetupSuite() {
	x.pool = pkg_testcontainers_postgresql.NewPostgreSQL(x.T())
	x.repository = repository_session_postgresql.NewSessions(client_postgresql.NewClient(x.pool))
}

func (x *RepositoryTestSuite) SetupTest() {
	repository_fixture.Reset(x.T(), x.pool)
}

func TestRepositoryTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(RepositoryTestSuite))
}

func assertSession(t *testing.T, expected, actual *aggregate_session.Session) {
	t.Helper()

	require.NotNil(t, actual)
	require.Equal(t, expected.GetSessionID(), actual.GetSessionID())
	require.Equal(t, expected.GetUserID(), actual.GetUserID())
	require.WithinDuration(t, expected.GetCreatedAt(), actual.GetCreatedAt(), 0)
	require.WithinDuration(t, expected.GetUpdatedAt(), actual.GetUpdatedAt(), 0)
	require.WithinDuration(t, expected.GetExpiresAt(), actual.GetExpiresAt(), 0)
}
