package rpc_service_geo_rpc_get_city

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	geo_v1 "github.com/Roum1212/hahaton/genproto/geo/v1"
	rpc_service_geo_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/geo/serde"
	service_get_city "github.com/Roum1212/hahaton/internal/app/service/get-city"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
)

type RPC struct {
	service service_get_city.Service
}

func (x RPC) GetCity(
	ctx context.Context,
	r *connect.Request[geo_v1.GetCityRequest],
) (*connect.Response[geo_v1.GetCityResponse], error) {
	cityID := r.Msg.GetId()

	response, err := x.service.GetCity(ctx, service_get_city.NewRequest(cityID))
	if err != nil {
		switch {
		case errors.Is(err, service_get_city.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_get_city.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "geo.v1.City",
				ResourceName: "id:" + cityID,
				Description:  "city not found",
			})
		default:
			return nil, rpc_error.NewInternalError()
		}
	}

	return connect.NewResponse(&geo_v1.GetCityResponse{
		City:           rpc_service_geo_serde.CityToPB(response.GetCity()),
		CatalogVersion: response.GetCatalogVersion(),
	}), nil
}

func NewRPC(service service_get_city.Service) RPC {
	return RPC{service: service}
}
