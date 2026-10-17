package client_city_catalog

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	model_city "github.com/Roum1212/hahaton/internal/domain/model/city"
	model_geo_point "github.com/Roum1212/hahaton/internal/domain/model/geo-point"
)

const coordinateDimension = 2

type Client struct {
	cities  map[string]model_city.City
	version string
}

func (x Client) GetCity(ctx context.Context, id string) (model_city.City, bool, error) {
	if err := ctx.Err(); err != nil {
		return model_city.City{}, false, err
	}

	value, ok := x.cities[id]

	return value, ok, nil
}

func (x Client) GetVersion() string {
	return x.version
}

func (x Client) ListCities(ctx context.Context) ([]model_city.City, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	values := make([]model_city.City, 0, len(x.cities))
	for _, value := range x.cities {
		values = append(values, value)
	}

	return values, nil
}

func readBoundaries(path string) (map[string][][][]model_geo_point.GeoPoint, error) {
	result := make(map[string][][][]model_geo_point.GeoPoint)
	if len(path) == 0 {
		return result, nil
	}

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	var document map[string]json.RawMessage
	if err := json.Unmarshal(data, &document); err != nil {
		return nil, err
	}

	for id, raw := range document {
		var source [][][][]float64
		if err := json.Unmarshal(raw, &source); err != nil {
			return nil, err
		}

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
						return nil, err
					}

					polygons[index][ringIndex] = append(polygons[index][ringIndex], point)
				}
			}
		}

		result[id] = polygons
	}

	return result, nil
}

func NewClient(csvPath, boundaryPath, version string) (_ Client, err error) {
	file, err := os.Open(filepath.Clean(csvPath))
	if err != nil {
		return Client{}, err
	}

	defer func() {
		closeError := file.Close()
		if err == nil {
			err = closeError
		}
	}()

	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err != nil {
		return Client{}, err
	}

	columns := make(map[string]int, len(header))
	for index, name := range header {
		columns[name] = index
	}

	for _, name := range []string{
		"fias_id",
		"city",
		"region",
		"geo_lat",
		"geo_lon",
	} {
		if _, ok := columns[name]; !ok {
			return Client{}, fmt.Errorf("missing city column: %s", name)
		}
	}

	boundaries, err := readBoundaries(boundaryPath)
	if err != nil {
		return Client{}, err
	}

	cities := make(map[string]model_city.City)

	for {
		record, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return Client{}, err
		}

		id, name := record[columns["fias_id"]], record[columns["city"]]
		if len(id) == 0 || len(name) == 0 {
			continue
		}

		var center *model_geo_point.GeoPoint

		latitude, latitudeError := strconv.ParseFloat(record[columns["geo_lat"]], 64)

		longitude, longitudeError := strconv.ParseFloat(record[columns["geo_lon"]], 64)
		if latitudeError == nil && longitudeError == nil {
			value, err := model_geo_point.NewGeoPoint(latitude, longitude)
			if err == nil {
				center = &value
			}
		}

		city, err := model_city.NewCity(id, name, record[columns["region"]], center, boundaries[id])
		if err != nil {
			return Client{}, err
		}

		cities[id] = city
	}

	if len(cities) == 0 {
		return Client{}, errors.New("city catalog is empty")
	}

	return Client{
		cities:  cities,
		version: version,
	}, nil
}
