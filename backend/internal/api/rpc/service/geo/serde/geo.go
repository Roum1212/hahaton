package rpc_service_geo_serde

import (
	"errors"

	geo_v1 "github.com/Roum1212/hahaton/genproto/geo/v1"
	model_address_suggestion "github.com/Roum1212/hahaton/internal/domain/model/address-suggestion"
	model_city "github.com/Roum1212/hahaton/internal/domain/model/city"
	model_geo_point "github.com/Roum1212/hahaton/internal/domain/model/geo-point"
	model_search_area "github.com/Roum1212/hahaton/internal/domain/model/search-area"
)

func CityToPB(city model_city.City) *geo_v1.City {
	result := &geo_v1.City{
		Id:          city.GetID(),
		Name:        city.GetName(),
		RegionName:  city.GetRegionName(),
		HasBoundary: city.HasBoundary(),
	}
	if point := city.GetCenter(); point != nil {
		result.Center = PointToPB(*point)
	}

	return result
}

func PointFromPB(point *geo_v1.GeoPoint) (model_geo_point.GeoPoint, error) {
	if point == nil || point.Latitude == nil || point.Longitude == nil {
		return model_geo_point.GeoPoint{}, errors.New("both point coordinates are required")
	}

	return model_geo_point.NewGeoPoint(point.GetLatitude(), point.GetLongitude())
}

func PointToPB(point model_geo_point.GeoPoint) *geo_v1.GeoPoint {
	latitude, longitude := point.GetLatitude(), point.GetLongitude()

	return &geo_v1.GeoPoint{
		Latitude:  &latitude,
		Longitude: &longitude,
	}
}

func SearchAreaFromPB(area *geo_v1.SearchArea) (model_search_area.SearchArea, error) {
	if area == nil {
		return model_search_area.SearchArea{}, errors.New("search area is required")
	}

	switch value := area.GetArea().(type) {
	case *geo_v1.SearchArea_CityId:
		return model_search_area.NewCitySearchArea(value.CityId)
	case *geo_v1.SearchArea_Radius:
		if value.Radius == nil {
			return model_search_area.SearchArea{}, errors.New("radius area is required")
		}

		point, err := PointFromPB(value.Radius.GetCenter())
		if err != nil {
			return model_search_area.SearchArea{}, err
		}

		return model_search_area.NewRadiusSearchArea(point, value.Radius.GetRadiusMeters())
	default:
		return model_search_area.SearchArea{}, errors.New("one search area is required")
	}
}

func SuggestionToPB(value model_address_suggestion.AddressSuggestion) *geo_v1.AddressSuggestion {
	data := &geo_v1.AddressSuggestionData{
		FiasId:             value.GetFiasID(),
		CityFiasId:         value.GetCityFiasID(),
		SettlementFiasId:   value.GetSettlementFiasID(),
		CityWithType:       value.GetCityWithType(),
		RegionWithType:     value.GetRegionWithType(),
		SettlementWithType: value.GetSettlementWithType(),
		QcGeo:              value.GetQcGeo(),
	}
	if point := value.GetPoint(); point != nil {
		data.Point = PointToPB(*point)
	}

	return &geo_v1.AddressSuggestion{
		Value:             value.GetValue(),
		UnrestrictedValue: value.GetUnrestrictedValue(),
		Data:              data,
	}
}
