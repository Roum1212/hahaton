package repository_city_postgresql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	model_city "github.com/Roum1212/hahaton/internal/domain/model/city"
	model_geo_point "github.com/Roum1212/hahaton/internal/domain/model/geo-point"
)

const (
	idColumn = "id"

	catalogVersionColumn = "catalog_version"
	latitudeColumn       = "latitude"
	longitudeColumn      = "longitude"
	nameColumn           = "name"
	polygonsColumn       = "polygons"
	regionNameColumn     = "region_name"
)

const coordinateDimension = 2

type Database interface {
	Exec(ctx context.Context, query string, arguments ...any) (pgconn.CommandTag, error)
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
	Query(ctx context.Context, query string, arguments ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, arguments ...any) pgx.Row
}

type Catalog struct {
	database Database
	queries  squirrel.StatementBuilderType
	version  string
}

func (x *Catalog) GetCity(ctx context.Context, id string) (model_city.City, bool, error) {
	query, arguments, err := x.queries.Select(
		idColumn,
		nameColumn,
		regionNameColumn,
		latitudeColumn,
		longitudeColumn,
		polygonsColumn,
	).
		From("cities").Where(squirrel.Eq{idColumn: id}).ToSql()
	if err != nil {
		return model_city.City{}, false, err
	}

	city, err := scanCity(x.database.QueryRow(ctx, query, arguments...))
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return model_city.City{}, false, nil
	case err != nil:
		return model_city.City{}, false, err
	}

	return city, true, nil
}

func (x *Catalog) GetVersion() string {
	return x.version
}

func (x *Catalog) ListCities(ctx context.Context) ([]model_city.City, error) {
	query, arguments, err := x.queries.Select(
		idColumn,
		nameColumn,
		regionNameColumn,
		latitudeColumn,
		longitudeColumn,
		polygonsColumn,
	).
		From("cities").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := x.database.Query(ctx, query, arguments...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cities := make([]model_city.City, 0)

	for rows.Next() {
		city, err := scanCity(rows)
		if err != nil {
			return nil, err
		}

		cities = append(cities, city)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cities, nil
}

func (x *Catalog) Save(ctx context.Context, cities []model_city.City, version string) error {
	if err := x.database.InTx(ctx, func(ctx context.Context) error {
		for _, city := range cities {
			if err := x.saveCity(ctx, city, version); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	x.version = version

	return nil
}

func decodePolygons(source [][][][]float64) ([][][]model_geo_point.GeoPoint, error) {
	polygons := make([][][]model_geo_point.GeoPoint, len(source))
	for index, polygon := range source {
		polygons[index] = make([][]model_geo_point.GeoPoint, len(polygon))
		for ringIndex, ring := range polygon {
			for _, pair := range ring {
				if len(pair) != coordinateDimension {
					return nil, errors.New("boundary coordinate must contain longitude and latitude")
				}

				point, err := model_geo_point.NewGeoPoint(pair[1], pair[0])
				if err != nil {
					return nil, fmt.Errorf("restore point: %w", err)
				}

				polygons[index][ringIndex] = append(polygons[index][ringIndex], point)
			}
		}
	}

	return polygons, nil
}

func (x *Catalog) saveCity(ctx context.Context, city model_city.City, version string) error {
	polygons := make([][][][]float64, 0, len(city.GetPolygons()))
	for _, polygon := range city.GetPolygons() {
		rings := make([][][]float64, 0, len(polygon))
		for _, ring := range polygon {
			points := make([][]float64, 0, len(ring))
			for _, point := range ring {
				points = append(points, []float64{point.GetLongitude(), point.GetLatitude()})
			}

			rings = append(rings, points)
		}

		polygons = append(polygons, rings)
	}

	data, err := json.Marshal(polygons)
	if err != nil {
		return err
	}

	var latitude, longitude *float64

	if center := city.GetCenter(); center != nil {
		lat, lon := center.GetLatitude(), center.GetLongitude()
		latitude, longitude = &lat, &lon
	}

	columns := []string{
		nameColumn,
		regionNameColumn,
		latitudeColumn,
		longitudeColumn,
		polygonsColumn,
		catalogVersionColumn,
	}

	assignments := make([]string, len(columns))
	for index, column := range columns {
		assignments[index] = column + " = EXCLUDED." + column
	}

	conflict := "ON CONFLICT (" + idColumn + ") DO UPDATE SET " + strings.Join(assignments, ", ")

	query, arguments, err := x.queries.Insert("cities").SetMap(map[string]any{
		idColumn:             city.GetID(),
		nameColumn:           city.GetName(),
		regionNameColumn:     city.GetRegionName(),
		latitudeColumn:       latitude,
		longitudeColumn:      longitude,
		polygonsColumn:       string(data),
		catalogVersionColumn: version,
	}).Suffix(conflict).ToSql()
	if err != nil {
		return err
	}

	_, err = x.database.Exec(ctx, query, arguments...)

	return err
}

func scanCity(row pgx.Row) (model_city.City, error) {
	var (
		id, name, region    string
		latitude, longitude *float64
		data                []byte
	)

	if err := row.Scan(&id, &name, &region, &latitude, &longitude, &data); err != nil {
		return model_city.City{}, err
	}

	var center *model_geo_point.GeoPoint

	if latitude != nil && longitude != nil {
		point, err := model_geo_point.NewGeoPoint(*latitude, *longitude)
		if err != nil {
			return model_city.City{}, fmt.Errorf("restore point: %w", err)
		}

		center = &point
	}

	var source [][][][]float64
	if err := json.Unmarshal(data, &source); err != nil {
		return model_city.City{}, err
	}

	polygons, err := decodePolygons(source)
	if err != nil {
		return model_city.City{}, err
	}

	return model_city.NewCity(id, name, region, center, polygons)
}

func NewCatalog(ctx context.Context, database Database) (*Catalog, error) {
	x := &Catalog{
		database: database,
		queries:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}

	query, arguments, err := x.queries.Select(catalogVersionColumn).From("cities").Limit(1).ToSql()
	if err != nil {
		return nil, err
	}

	err = database.QueryRow(ctx, query, arguments...).Scan(&x.version)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return x, nil
	case err != nil:
		return nil, err
	}

	return x, nil
}
