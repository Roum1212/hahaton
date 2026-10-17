package model_geo_point

import (
	"errors"
	"math"
)

const (
	maxLatitude  = 90
	maxLongitude = 180
)

type GeoPoint struct {
	latitude  float64
	longitude float64
}

func (x GeoPoint) GetLatitude() float64 {
	return x.latitude
}

func (x GeoPoint) GetLongitude() float64 {
	return x.longitude
}

func (x GeoPoint) Validate() error {
	if math.IsNaN(x.latitude) || math.IsInf(x.latitude, 0) || math.Abs(x.latitude) > maxLatitude {
		return errors.New("latitude must be finite and between -90 and 90")
	}

	if math.IsNaN(x.longitude) || math.IsInf(x.longitude, 0) || math.Abs(x.longitude) > maxLongitude {
		return errors.New("longitude must be finite and between -180 and 180")
	}

	return nil
}

func NewGeoPoint(latitude, longitude float64) (GeoPoint, error) {
	x := GeoPoint{
		latitude:  latitude,
		longitude: longitude,
	}
	if err := x.Validate(); err != nil {
		return GeoPoint{}, err
	}

	return x, nil
}
