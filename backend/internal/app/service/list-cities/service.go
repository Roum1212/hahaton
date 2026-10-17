package service_list_cities

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	model_city "github.com/Roum1212/hahaton/internal/domain/model/city"
	"github.com/Roum1212/hahaton/internal/pkg/page"
)

var ErrInvalidRequest = errors.New("invalid request")

type Catalog interface {
	ListCities(ctx context.Context) ([]model_city.City, error)
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

	values, err := x.catalog.ListCities(ctx)
	if err != nil {
		return Response{}, err
	}

	query, _ := r.GetQuery().Get()
	query = strings.ToLower(strings.TrimSpace(query))
	filtered := make([]model_city.City, 0)

	for _, city := range values {
		if strings.Contains(strings.ToLower(city.GetName()), query) {
			filtered = append(filtered, city)
		}
	}

	slices.SortFunc(filtered, func(
		a model_city.City,
		b model_city.City,
	) int {
		if value := strings.Compare(a.GetName(), b.GetName()); value != 0 {
			return value
		}

		if value := strings.Compare(a.GetRegionName(), b.GetRegionName()); value != 0 {
			return value
		}

		return strings.Compare(a.GetID(), b.GetID())
	})

	values, total := page.Slice(filtered, r.GetPage())

	return NewResponse(values, total, x.catalog.GetVersion()), nil
}

func NewService(catalog Catalog) Service {
	return service{
		catalog: catalog,
	}
}
