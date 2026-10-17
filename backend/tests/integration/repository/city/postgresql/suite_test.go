//go:build integration

package repository_city_postgresql_test

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"

	client_postgresql "github.com/Roum1212/hahaton/internal/infra/client/postgresql"
	repository_city_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/city/postgresql"
	pkg_testcontainers_postgresql "github.com/Roum1212/hahaton/pkg/testcontainers/postgresql"
	repository_fixture "github.com/Roum1212/hahaton/tests/integration/repository/fixture"
)

type RepositoryTestSuite struct {
	suite.Suite

	pool       *pgxpool.Pool
	database   client_postgresql.Client
	repository *repository_city_postgresql.Catalog
}

func (x *RepositoryTestSuite) SetupSuite() {
	x.pool = pkg_testcontainers_postgresql.NewPostgreSQL(x.T())
	x.database = client_postgresql.NewClient(x.pool)
	repository, err := repository_city_postgresql.NewCatalog(x.T().Context(), x.database)
	x.Require().NoError(err)
	x.repository = repository
}

func (x *RepositoryTestSuite) SetupTest() {
	repository_fixture.Reset(x.T(), x.pool)
	repository, err := repository_city_postgresql.NewCatalog(x.T().Context(), x.database)
	x.Require().NoError(err)
	x.repository = repository
}

func TestRepositoryTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RepositoryTestSuite))
}
