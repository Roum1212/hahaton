package model_city

import (
	"errors"
	"math"
	"strings"

	model_geo_point "github.com/Roum1212/hahaton/internal/domain/model/geo-point"
)

const edgeTolerance = 1e-10

type City struct {
	id         string
	name       string
	regionName string
	center     *model_geo_point.GeoPoint
	polygons   [][][]model_geo_point.GeoPoint
}

func (x City) Contains(point model_geo_point.GeoPoint) bool {
	for _, polygon := range x.polygons {
		if len(polygon) == 0 {
			continue
		}

		inside, boundary := insideRing(point, polygon[0])
		if boundary {
			return true
		}

		if !inside {
			continue
		}

		hole := false

		for _, ring := range polygon[1:] {
			inHole, onEdge := insideRing(point, ring)
			if onEdge {
				return true
			}

			if inHole {
				hole = true

				break
			}
		}

		if !hole {
			return true
		}
	}

	return false
}

func (x City) GetCenter() *model_geo_point.GeoPoint {
	if x.center == nil {
		return nil
	}

	center := *x.center

	return &center
}

func (x City) GetID() string {
	return x.id
}

func (x City) GetName() string {
	return x.name
}

func (x City) GetPolygons() [][][]model_geo_point.GeoPoint {
	polygons := make([][][]model_geo_point.GeoPoint, len(x.polygons))
	for index, polygon := range x.polygons {
		polygons[index] = make([][]model_geo_point.GeoPoint, len(polygon))
		for ringIndex, ring := range polygon {
			polygons[index][ringIndex] = append([]model_geo_point.GeoPoint(nil), ring...)
		}
	}

	return polygons
}

func (x City) GetRegionName() string {
	return x.regionName
}

func (x City) HasBoundary() bool {
	return len(x.polygons) != 0
}

func (x City) Validate() error {
	if len(strings.TrimSpace(x.id)) == 0 || len(strings.TrimSpace(x.name)) == 0 {
		return errors.New("city id and name are required")
	}

	if x.center != nil {
		if err := x.center.Validate(); err != nil {
			return err
		}
	}

	for _, polygon := range x.polygons {
		if len(polygon) == 0 {
			return errors.New("polygon must contain an outer ring")
		}

		for _, ring := range polygon {
			if len(ring) < 4 || ring[0] != ring[len(ring)-1] {
				return errors.New("city ring must be closed and have at least four points")
			}

			for _, point := range ring {
				if err := point.Validate(); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func insideRing(point model_geo_point.GeoPoint, ring []model_geo_point.GeoPoint) (bool, bool) {
	inside := false

	for index, current := range ring {
		previous := ring[(index+len(ring)-1)%len(ring)]
		ax, ay := previous.GetLongitude(), previous.GetLatitude()
		bx, by := current.GetLongitude(), current.GetLatitude()
		px, py := point.GetLongitude(), point.GetLatitude()

		cross := (px-ax)*(by-ay) - (py-ay)*(bx-ax)
		if math.Abs(cross) <= edgeTolerance && px >= min(ax, bx)-edgeTolerance && px <= max(ax, bx)+edgeTolerance &&
			py >= min(ay, by)-edgeTolerance &&
			py <= max(ay, by)+edgeTolerance {
			return true, true
		}

		if (ay > py) != (by > py) && px < (bx-ax)*(py-ay)/(by-ay)+ax {
			inside = !inside
		}
	}

	return inside, false
}

func NewCity(
	id string,
	name string,
	regionName string,
	center *model_geo_point.GeoPoint,
	polygons [][][]model_geo_point.GeoPoint,
) (City, error) {
	copyPolygons := make([][][]model_geo_point.GeoPoint, len(polygons))
	for index, polygon := range polygons {
		copyPolygons[index] = make([][]model_geo_point.GeoPoint, len(polygon))
		for ringIndex, ring := range polygon {
			copyPolygons[index][ringIndex] = append([]model_geo_point.GeoPoint(nil), ring...)
		}
	}

	if center != nil {
		value := *center
		center = &value
	}

	x := City{
		id:         id,
		name:       name,
		regionName: regionName,
		center:     center,
		polygons:   copyPolygons,
	}
	if err := x.Validate(); err != nil {
		return City{}, err
	}

	return x, nil
}
