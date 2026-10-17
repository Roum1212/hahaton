package service_list_cities

import (
	"context"
	"errors"
	"fmt"
	"strings"

	model_city "github.com/Roum1212/hahaton/internal/domain/model/city"
)

var ErrInvalidRequest = errors.New("invalid request")

type Catalog interface {
	FindCities(ctx context.Context, query string, count int32) ([]model_city.City, error)
	GetVersion() string
}
type (
	Service interface {
		ListCities(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		catalog Catalog
	}
)

func (x service) ListCities(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	query, _ := r.GetQuery().Get()
	query = strings.TrimSpace(query)

	if len(query) == 0 {
		return NewResponse(nil, 0, x.catalog.GetVersion()), nil
	}

	values, err := x.catalog.FindCities(ctx, query, r.GetPage().GetLimit())
	if err != nil {
		return Response{}, err
	}

	return NewResponse(values, int64(len(values)), x.catalog.GetVersion()), nil
}

func NewService(catalog Catalog) Service {
	return service{
		catalog: catalog,
	}
}
