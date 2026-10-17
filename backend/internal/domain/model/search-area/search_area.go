package model_search_area

import (
	"errors"
	"math"

	model_geo_point "github.com/Roum1212/hahaton/internal/domain/model/geo-point"
)

const (
	halfAngleDivisor    = 2
	earthDiameterMeters = 2 * earthRadiusMeters
	MaxRadiusMeters     = 200000
	earthRadiusMeters   = 6371000
	degreesInHalfCircle = 180
)

type SearchArea struct {
	cityID       string
	center       model_geo_point.GeoPoint
	radiusMeters int32
}

func (x SearchArea) Contains(point model_geo_point.GeoPoint) bool {
	radians := math.Pi / degreesInHalfCircle
	latitudeDelta := (point.GetLatitude() - x.center.GetLatitude()) * radians
	longitudeDelta := (point.GetLongitude() - x.center.GetLongitude()) * radians
	halfLatitude := math.Sin(latitudeDelta / halfAngleDivisor)
	halfLongitude := math.Sin(longitudeDelta / halfAngleDivisor)
	h := halfLatitude*halfLatitude + math.Cos(x.center.GetLatitude()*radians)*math.Cos(
		point.GetLatitude()*radians,
	)*halfLongitude*halfLongitude
	distance := earthDiameterMeters * math.Asin(math.Sqrt(min(1, h)))

	return distance <= float64(x.radiusMeters)
}

func (x SearchArea) GetCityID() string {
	return x.cityID
}

func (x SearchArea) Validate() error {
	if len(x.cityID) != 0 {
		if x.radiusMeters != 0 {
			return errors.New("city and radius cannot be combined")
		}

		return nil
	}

	if x.radiusMeters < 1 || x.radiusMeters > MaxRadiusMeters {
		return errors.New("radius must be between 1 and 200000 meters")
	}

	return x.center.Validate()
}

func NewCitySearchArea(cityID string) (SearchArea, error) {
	x := SearchArea{cityID: cityID}
	if err := x.Validate(); err != nil {
		return SearchArea{}, err
	}

	return x, nil
}

func NewRadiusSearchArea(center model_geo_point.GeoPoint, radiusMeters int32) (SearchArea, error) {
	x := SearchArea{
		center:       center,
		radiusMeters: radiusMeters,
	}
	if err := x.Validate(); err != nil {
		return SearchArea{}, err
	}

	return x, nil
}
