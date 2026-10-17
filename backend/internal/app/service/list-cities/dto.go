package service_list_cities

import (
	"errors"
	"strings"
	"unicode/utf8"

	model_city "github.com/Roum1212/hahaton/internal/domain/model/city"
	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

const (
	MaxCitySuggestions = 20
	maxQueryLength     = 300
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
	if err := x.page.Validate(); err != nil {
		return err
	}

	if x.page.GetLimit() > MaxCitySuggestions || x.page.GetOffset() != 0 {
		return errors.New("city suggestions require limit up to 20 and zero offset")
	}

	query, _ := x.query.Get()
	if utf8.RuneCountInString(strings.TrimSpace(query)) > maxQueryLength {
		return errors.New("city query must contain at most 300 characters")
	}

	return nil
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
