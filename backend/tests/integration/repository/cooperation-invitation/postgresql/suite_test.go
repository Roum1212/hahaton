//go:build integration

package repository_cooperation_invitation_postgresql_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	aggregate_cooperation_invitation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation-invitation"
	aggregate_cooperation_invitation_command "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation-invitation/command"
	client_encryption "github.com/Roum1212/hahaton/internal/infra/client/encryption"
	client_postgresql "github.com/Roum1212/hahaton/internal/infra/client/postgresql"
	repository_cooperation_invitation_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/cooperation-invitation/postgresql"
	pkg_testcontainers_postgresql "github.com/Roum1212/hahaton/pkg/testcontainers/postgresql"
	repository_fixture "github.com/Roum1212/hahaton/tests/integration/repository/fixture"
)

type RepositoryTestSuite struct {
	suite.Suite

	pool       *pgxpool.Pool
	database   client_postgresql.Client
	repository aggregate_cooperation_invitation.Repository
}

func (x *RepositoryTestSuite) SetupSuite() {
	x.pool = pkg_testcontainers_postgresql.NewPostgreSQL(x.T())
	x.database = client_postgresql.NewClient(x.pool)
	x.repository = repository_cooperation_invitation_postgresql.NewCooperationInvitations(x.database, crypto(x.T()))
}

func (x *RepositoryTestSuite) SetupTest() {
	repository_fixture.Reset(x.T(), x.pool)
}

func (x *RepositoryTestSuite) changeValue(value *aggregate_cooperation_invitation.CooperationInvitation) {
	value.Create(aggregate_cooperation_invitation_command.NewCreate(
		"updated-token",
		time.Time{},
		repository_fixture.Now(),
	))
}

func (x *RepositoryTestSuite) newValue() *aggregate_cooperation_invitation.CooperationInvitation {
	return repository_fixture.NewInvitation(x.T(), repository_fixture.CreateCoordinator(x.T(), x.database).GetUserID())
}

func TestRepositoryTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RepositoryTestSuite))
}

func assertCooperationInvitation(
	t *testing.T,
	expected *aggregate_cooperation_invitation.CooperationInvitation,
	actual *aggregate_cooperation_invitation.CooperationInvitation,
) {
	t.Helper()

	require.NotNil(t, actual)
	require.Equal(t, expected.GetCoordinatorID(), actual.GetCoordinatorID())
	require.WithinDuration(t, expected.GetCreatedAt(), actual.GetCreatedAt(), 0)
	require.WithinDuration(t, expected.GetExpiresAt(), actual.GetExpiresAt(), 0)
	require.Equal(t, expected.GetInvitationID(), actual.GetInvitationID())
	require.Equal(t, expected.GetToken(), actual.GetToken())
}

func crypto(t *testing.T) client_encryption.Client {
	t.Helper()

	client, err := client_encryption.NewClient(bytes.Repeat([]byte{7}, 32))
	require.NoError(t, err)

	return client
}
