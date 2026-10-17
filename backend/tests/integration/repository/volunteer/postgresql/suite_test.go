//go:build integration

package repository_volunteer_postgresql_test

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
	aggregate_volunteer_command "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer/command"
	model_date_of_birth "github.com/Roum1212/hahaton/internal/domain/model/date-of-birth"
	model_full_name "github.com/Roum1212/hahaton/internal/domain/model/full-name"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	client_postgresql "github.com/Roum1212/hahaton/internal/infra/client/postgresql"
	repository_volunteer_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/volunteer/postgresql"
	pkg_testcontainers_postgresql "github.com/Roum1212/hahaton/pkg/testcontainers/postgresql"
	repository_fixture "github.com/Roum1212/hahaton/tests/integration/repository/fixture"
)

type RepositoryTestSuite struct {
	suite.Suite

	pool       *pgxpool.Pool
	database   client_postgresql.Client
	repository aggregate_volunteer.Repository
}

func (x *RepositoryTestSuite) SetupSuite() {
	x.pool = pkg_testcontainers_postgresql.NewPostgreSQL(x.T())
	x.database = client_postgresql.NewClient(x.pool)
	x.repository = repository_volunteer_postgresql.NewVolunteers(x.database)
}

func (x *RepositoryTestSuite) SetupTest() {
	repository_fixture.Reset(x.T(), x.pool)
}

func (x *RepositoryTestSuite) changeValue(value *aggregate_volunteer.Volunteer) {
	name, err := model_full_name.NewFullName("Пётр", "Петров", "")
	x.Require().NoError(err)
	date, err := model_date_of_birth.NewDateOfBirth(1999, time.December, 31)
	x.Require().NoError(err)
	value.Update(aggregate_volunteer_command.NewUpdate(
		name,
		date,
		model_image_path.ImagePath{},
		repository_fixture.Now().Add(time.Minute),
	))
}

func (x *RepositoryTestSuite) newValue() *aggregate_volunteer.Volunteer {
	return repository_fixture.NewVolunteer(x.T())
}

func TestRepositoryTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RepositoryTestSuite))
}

func assertVolunteer(t *testing.T, expected, actual *aggregate_volunteer.Volunteer) {
	t.Helper()

	require.NotNil(t, actual)
	require.Equal(t, expected.GetAvatarPath(), actual.GetAvatarPath())
	require.WithinDuration(t, expected.GetCreatedAt(), actual.GetCreatedAt(), 0)
	require.Equal(t, expected.GetDateOfBirth(), actual.GetDateOfBirth())
	require.Equal(t, expected.GetEmailAddress(), actual.GetEmailAddress())
	require.Equal(t, expected.GetFullName(), actual.GetFullName())
	require.Equal(t, expected.GetIdentity(), actual.GetIdentity())
	require.Equal(t, expected.GetPhoneNumber(), actual.GetPhoneNumber())
	require.WithinDuration(t, expected.GetUpdatedAt(), actual.GetUpdatedAt(), 0)
	require.Equal(t, expected.GetUserID(), actual.GetUserID())
}
