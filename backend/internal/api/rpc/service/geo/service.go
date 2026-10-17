package rpc_service_geo

import (
	"context"

	"connectrpc.com/connect"

	geo_v1 "github.com/Roum1212/hahaton/genproto/geo/v1"
	geo_v1connect "github.com/Roum1212/hahaton/genproto/geo/v1/geo_v1connect"
	rpc_service_geo_rpc_get_city "github.com/Roum1212/hahaton/internal/api/rpc/service/geo/rpc/get-city"
	rpc_service_geo_rpc_list_cities "github.com/Roum1212/hahaton/internal/api/rpc/service/geo/rpc/list-cities"
	rpc_service_geo_rpc_suggest_addresses "github.com/Roum1212/hahaton/internal/api/rpc/service/geo/rpc/suggest-addresses"
)

var _ geo_v1connect.GeoServiceHandler = Service{}

type Service struct {
	getCityRPC          rpc_service_geo_rpc_get_city.RPC
	listCitiesRPC       rpc_service_geo_rpc_list_cities.RPC
	suggestAddressesRPC rpc_service_geo_rpc_suggest_addresses.RPC
}

func (x Service) GetCity(
	ctx context.Context,
	r *connect.Request[geo_v1.GetCityRequest],
) (*connect.Response[geo_v1.GetCityResponse], error) {
	return x.getCityRPC.GetCity(ctx, r)
}

func (x Service) ListCities(
	ctx context.Context,
	r *connect.Request[geo_v1.ListCitiesRequest],
) (*connect.Response[geo_v1.ListCitiesResponse], error) {
	return x.listCitiesRPC.ListCities(ctx, r)
}

func (x Service) SuggestAddresses(
	ctx context.Context,
	r *connect.Request[geo_v1.SuggestAddressesRequest],
) (*connect.Response[geo_v1.SuggestAddressesResponse], error) {
	return x.suggestAddressesRPC.SuggestAddresses(ctx, r)
}

func NewService(
	getCityRPC rpc_service_geo_rpc_get_city.RPC,
	listCitiesRPC rpc_service_geo_rpc_list_cities.RPC,
	suggestAddressesRPC rpc_service_geo_rpc_suggest_addresses.RPC,
) Service {
	return Service{
		getCityRPC:          getCityRPC,
		listCitiesRPC:       listCitiesRPC,
		suggestAddressesRPC: suggestAddressesRPC,
	}
}
