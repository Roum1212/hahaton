package service_get_city

import (
	"errors"
	"strings"

	model_city "github.com/Roum1212/hahaton/internal/domain/model/city"
)

type Request struct {
	cityID string
}

func (x Request) GetCityID() string {
	return x.cityID
}

func (x Request) Validate() error {
	if len(strings.TrimSpace(x.cityID)) == 0 {
		return errors.New("city id is required")
	}

	return nil
}

func NewRequest(cityID string) Request {
	return Request{
		cityID: cityID,
	}
}

type Response struct {
	city           model_city.City
	catalogVersion string
}

func (x Response) GetCatalogVersion() string {
	return x.catalogVersion
}

func (x Response) GetCity() model_city.City {
	return x.city
}

func NewResponse(city model_city.City, catalogVersion string) Response {
	return Response{
		city:           city,
		catalogVersion: catalogVersion,
	}
}
