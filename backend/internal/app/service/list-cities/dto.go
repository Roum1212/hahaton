package service_list_cities

import (
	model_city "github.com/Roum1212/hahaton/internal/domain/model/city"
	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

type Request struct {
	query utils_value.Optional[string]

	page model_page.Page
}

func (x Request) GetPage() model_page.Page {
	return x.page
}

func (x Request) GetQuery() utils_value.Optional[string] {
	return x.query
}

func (x Request) Validate() error {
	return x.page.Validate()
}

type RequestOption func(x *Request)

func WithQuery(value string) RequestOption {
	return func(x *Request) {
		x.query = utils_value.NewOptional(value)
	}
}

func NewRequest(page model_page.Page, opts ...RequestOption) Request {
	x := Request{
		page: page,
	}
	for _, option := range opts {
		option(&x)
	}

	return x
}

type Response struct {
	cities         []model_city.City
	totalCount     int64
	catalogVersion string
}

func (x Response) GetCatalogVersion() string {
	return x.catalogVersion
}

func (x Response) GetCities() []model_city.City {
	return x.cities
}

func (x Response) GetTotalCount() int64 {
	return x.totalCount
}

func NewResponse(cities []model_city.City, totalCount int64, catalogVersion string) Response {
	return Response{
		cities:         cities,
		totalCount:     totalCount,
		catalogVersion: catalogVersion,
	}
}
