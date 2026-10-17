package rpc_service_geo_rpc_list_cities

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	geo_v1 "github.com/Roum1212/hahaton/genproto/geo/v1"
	rpc_service_geo_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/geo/serde"
	service_list_cities "github.com/Roum1212/hahaton/internal/app/service/list-cities"
	service_suggest_addresses "github.com/Roum1212/hahaton/internal/app/service/suggest-addresses"
	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
)

const (
	defaultCityLimit = 10
	maxCityLimit     = service_list_cities.MaxCitySuggestions
)

type RPC struct {
	service service_list_cities.Service
}

func (x RPC) ListCities(
	ctx context.Context,
	r *connect.Request[geo_v1.ListCitiesRequest],
) (*connect.Response[geo_v1.ListCitiesResponse], error) {
	var opts []service_list_cities.RequestOption

	if value := r.Msg.Query; value != nil {
		opts = append(opts, service_list_cities.WithQuery(*value))
	}

	paging, err := model_page.NewPageWithLimits(r.Msg.GetLimit(), r.Msg.GetOffset(), defaultCityLimit, maxCityLimit)
	if err != nil {
		return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "pagination",
					Description: err.Error(),
				},
			},
		})
	}

	response, err := x.service.ListCities(ctx, service_list_cities.NewRequest(paging, opts...))
	if err != nil {
		switch {
		case errors.Is(err, service_list_cities.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_suggest_addresses.ErrUnavailable):
			return nil, rpc_error.NewUnavailableError()
		case errors.Is(err, service_suggest_addresses.ErrQuota):
			return nil, rpc_error.NewResourceExhaustedError(nil)
		default:
			return nil, rpc_error.NewInternalError()
		}
	}

	values := make([]*geo_v1.City, 0, len(response.GetCities()))
	for _, value := range response.GetCities() {
		values = append(values, rpc_service_geo_serde.CityToPB(value))
	}

	return connect.NewResponse(&geo_v1.ListCitiesResponse{
		Cities:         values,
		TotalCount:     response.GetTotalCount(),
		CatalogVersion: response.GetCatalogVersion(),
	}), nil
}

func NewRPC(service service_list_cities.Service) RPC {
	return RPC{service: service}
}
