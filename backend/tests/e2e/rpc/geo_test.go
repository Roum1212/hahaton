//go:build e2e

package rpc

import (
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"

	geo_v1 "github.com/Roum1212/hahaton/genproto/geo/v1"
)

func (x *RPCTestSuite) TestCityCatalogAndProviderErrors() {
	ctx := x.T().Context()
	cities, err := x.clients.geo.ListCities(
		ctx,
		request(actor{}, &geo_v1.ListCitiesRequest{Query: proto.String("  ОБНИНСК ")}),
	)
	x.Require().NoError(err)
	x.Require().Len(cities.Msg.GetCities(), 1)
	x.Require().EqualValues(1, cities.Msg.GetTotalCount())
	x.Require().Equal("e2e", cities.Msg.GetCatalogVersion())
	city, err := x.clients.geo.GetCity(ctx, request(actor{}, &geo_v1.GetCityRequest{Id: "obninsk"}))
	x.Require().NoError(err)
	x.Require().Equal("Обнинск", city.Msg.GetCity().GetName())
	x.Require().Equal(55.1, city.Msg.GetCity().GetCenter().GetLatitude())
	_, err = x.clients.geo.GetCity(ctx, request(actor{}, &geo_v1.GetCityRequest{Id: "missing"}))
	callError(x.T(), err, connect.CodeNotFound)
	a := x.login()
	query := &geo_v1.SuggestAddressesRequest{
		Query:     " Обнинск ",
		FromBound: &geo_v1.AddressBound{Value: "city"},
		ToBound:   &geo_v1.AddressBound{Value: "settlement"},
	}
	suggestions, err := x.clients.geo.SuggestAddresses(ctx, request(a, query))
	x.Require().NoError(err)
	x.Require().Len(suggestions.Msg.GetSuggestions(), 1)
	x.Require().Equal("obninsk", suggestions.Msg.GetSuggestions()[0].GetData().GetFiasId())
	x.Require().Equal(36.6, suggestions.Msg.GetSuggestions()[0].GetData().GetPoint().GetLongitude())

	x.provider.addressStatus.Store(http.StatusTooManyRequests)
	defer x.provider.addressStatus.Store(http.StatusOK)

	_, err = x.clients.geo.SuggestAddresses(ctx, request(a, query))
	callError(x.T(), err, connect.CodeResourceExhausted)
	x.provider.addressStatus.Store(http.StatusInternalServerError)
	_, err = x.clients.geo.SuggestAddresses(ctx, request(a, query))
	callError(x.T(), err, connect.CodeUnavailable)
	_, err = x.clients.geo.SuggestAddresses(ctx, request(a, &geo_v1.SuggestAddressesRequest{
		Query:     "Обнинск",
		FromBound: &geo_v1.AddressBound{Value: "house"},
		ToBound:   &geo_v1.AddressBound{Value: "city"},
	}))
	callError(x.T(), err, connect.CodeInvalidArgument)
}
