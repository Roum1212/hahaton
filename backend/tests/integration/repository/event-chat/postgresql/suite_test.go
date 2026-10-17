//go:build integration

package repository_event_chat_postgresql_test

import (
	"bytes"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	aggregate_event_chat "github.com/Roum1212/hahaton/internal/domain/aggregate/event-chat"
	client_encryption "github.com/Roum1212/hahaton/internal/infra/client/encryption"
	client_postgresql "github.com/Roum1212/hahaton/internal/infra/client/postgresql"
	repository_event_chat_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/event-chat/postgresql"
	pkg_testcontainers_postgresql "github.com/Roum1212/hahaton/pkg/testcontainers/postgresql"
	repository_fixture "github.com/Roum1212/hahaton/tests/integration/repository/fixture"
)

type RepositoryTestSuite struct {
	suite.Suite

	pool       *pgxpool.Pool
	database   client_postgresql.Client
	repository aggregate_event_chat.Repository
}

func (x *RepositoryTestSuite) SetupSuite() {
	x.pool = pkg_testcontainers_postgresql.NewPostgreSQL(x.T())
	x.database = client_postgresql.NewClient(x.pool)
	crypto, err := client_encryption.NewClient(bytes.Repeat([]byte{7}, 32))
	x.Require().NoError(err)
	x.repository = repository_event_chat_postgresql.NewEventChats(x.database, crypto)
}

func (x *RepositoryTestSuite) SetupTest() {
	repository_fixture.Reset(x.T(), x.pool)
}

func (x *RepositoryTestSuite) newValue() *aggregate_event_chat.EventChat {
	return aggregate_event_chat.NewEventChat(repository_fixture.CreateEvent(x.T(), x.database).GetEventID())
}

func TestRepositoryTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RepositoryTestSuite))
}

func assertEventChat(t *testing.T, expected, actual *aggregate_event_chat.EventChat) {
	t.Helper()

	require.NotNil(t, actual)
	require.Equal(t, expected.GetEventID(), actual.GetEventID())
	require.Equal(t, expected.GetLatestSequence(), actual.GetLatestSequence())
}
