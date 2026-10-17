//go:build integration

package repository_city_postgresql_test

import (
	"context"

	model_city "github.com/Roum1212/hahaton/internal/domain/model/city"
	model_geo_point "github.com/Roum1212/hahaton/internal/domain/model/geo-point"
	repository_city_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/city/postgresql"
)

func (x *RepositoryTestSuite) TestRepository_GetCity_NotFound() {
	city, found, err := x.repository.GetCity(x.T().Context(), "unknown-city")
	x.Require().NoError(err)
	x.Require().False(found)
	x.Require().Zero(city)
	x.Require().Empty(x.repository.GetVersion())
}

func (x *RepositoryTestSuite) TestRepository_ListCities() {
	ctx := x.T().Context()
	cities, err := x.repository.ListCities(ctx)
	x.Require().NoError(err)
	x.Require().Empty(cities)
	first := x.newCity("first", true)
	second := x.newCity("second", false)
	x.Require().NoError(x.repository.Save(ctx, []model_city.City{first, second}, "v1"))
	cities, err = x.repository.ListCities(ctx)
	x.Require().NoError(err)
	x.Require().ElementsMatch([]model_city.City{first, second}, cities)
}

func (x *RepositoryTestSuite) TestRepository_SaveAndReload() {
	ctx := x.T().Context()
	city := x.newCity("obninsk", true)
	x.Require().NoError(x.repository.Save(ctx, []model_city.City{city}, "v1"))
	got, found, err := x.repository.GetCity(ctx, city.GetID())
	x.Require().NoError(err)
	x.Require().True(found)
	x.Require().Equal(city, got)
	x.Require().Equal("v1", x.repository.GetVersion())
	reloaded, err := repository_city_postgresql.NewCatalog(ctx, x.database)
	x.Require().NoError(err)
	x.Require().Equal("v1", reloaded.GetVersion())
	got, found, err = reloaded.GetCity(ctx, city.GetID())
	x.Require().NoError(err)
	x.Require().True(found)
	x.Require().Equal(city, got)
}

func (x *RepositoryTestSuite) TestRepository_Save_CanceledContext() {
	ctx := x.T().Context()
	city := x.newCity("obninsk", false)
	x.Require().NoError(x.repository.Save(ctx, []model_city.City{city}, "v1"))
	canceled, cancel := context.WithCancel(ctx)
	cancel()

	err := x.repository.Save(canceled, []model_city.City{x.newCity("new", true)}, "v2")
	x.Require().ErrorIs(err, context.Canceled)
	x.Require().Equal("v1", x.repository.GetVersion())
	cities, err := x.repository.ListCities(ctx)
	x.Require().NoError(err)
	x.Require().Len(cities, 1)
	x.Require().Equal(city, cities[0])
}

func (x *RepositoryTestSuite) TestRepository_Save_UpdateAndClearGeometry() {
	ctx := x.T().Context()
	city := x.newCity("obninsk", true)
	x.Require().NoError(x.repository.Save(ctx, []model_city.City{city}, "v1"))
	updated, err := model_city.NewCity(city.GetID(), "Новое название", "Новая область", nil, nil)
	x.Require().NoError(err)
	x.Require().NoError(x.repository.Save(ctx, []model_city.City{updated}, "v2"))
	got, found, err := x.repository.GetCity(ctx, city.GetID())
	x.Require().NoError(err)
	x.Require().True(found)
	x.Require().Equal(updated, got)
	cities, err := x.repository.ListCities(ctx)
	x.Require().NoError(err)
	x.Require().Len(cities, 1)
	reloaded, err := repository_city_postgresql.NewCatalog(ctx, x.database)
	x.Require().NoError(err)
	x.Require().Equal("v2", reloaded.GetVersion())
}

func (x *RepositoryTestSuite) newCity(id string, withGeometry bool) model_city.City {
	if !withGeometry {
		city, err := model_city.NewCity(id, "Обнинск", "Калужская область", nil, nil)
		x.Require().NoError(err)

		return city
	}

	center, err := model_geo_point.NewGeoPoint(55.1, 36.6)
	x.Require().NoError(err)
	first, err := model_geo_point.NewGeoPoint(55, 36)
	x.Require().NoError(err)
	second, err := model_geo_point.NewGeoPoint(56, 36)
	x.Require().NoError(err)
	third, err := model_geo_point.NewGeoPoint(56, 37)
	x.Require().NoError(err)
	city, err := model_city.NewCity(id, "Обнинск", "Калужская область", &center, [][][]model_geo_point.GeoPoint{
		{{
			first,
			second,
			third,
			first,
		}},
	})
	x.Require().NoError(err)

	return city
}
