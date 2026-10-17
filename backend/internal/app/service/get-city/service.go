package service_get_city

import (
	"context"
	"errors"
	"fmt"

	model_city "github.com/Roum1212/hahaton/internal/domain/model/city"
)

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrNotFound       = errors.New("not found")
)

type Catalog interface {
	GetCity(ctx context.Context, id string) (model_city.City, bool, error)
	GetVersion() string
}
type (
	Service interface {
		GetCity(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		catalog Catalog
	}
)

func (x service) GetCity(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	city, found, err := x.catalog.GetCity(ctx, r.GetCityID())
	if err != nil {
		return Response{}, err
	}

	if !found {
		return Response{}, ErrNotFound
	}

	return NewResponse(city, x.catalog.GetVersion()), nil
}

func NewService(catalog Catalog) Service {
	return service{
		catalog: catalog,
	}
}
