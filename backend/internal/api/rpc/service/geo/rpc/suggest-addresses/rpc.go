package rpc_service_geo_rpc_suggest_addresses

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	rpc_error "go.flida.dev/rpc-error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	geo_v1 "github.com/Roum1212/hahaton/genproto/geo/v1"
	rpc_service_geo_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/geo/serde"
	service_suggest_addresses "github.com/Roum1212/hahaton/internal/app/service/suggest-addresses"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_suggest_addresses.Service
}

func (x RPC) SuggestAddresses(
	ctx context.Context,
	r *connect.Request[geo_v1.SuggestAddressesRequest],
) (*connect.Response[geo_v1.SuggestAddressesResponse], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	query := r.Msg.GetQuery()
	fromBound := r.Msg.GetFromBound().GetValue()
	toBound := r.Msg.GetToBound().GetValue()
	count := r.Msg.GetCount()

	response, err := x.service.SuggestAddresses(ctx, service_suggest_addresses.NewRequest(
		session.GetUserID(),
		query,
		fromBound,
		toBound,
		count,
	))
	if err != nil {
		switch {
		case errors.Is(err, service_suggest_addresses.ErrInvalidRequest):
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

	values := make([]*geo_v1.AddressSuggestion, 0, len(response.GetSuggestions()))
	for _, value := range response.GetSuggestions() {
		values = append(values, rpc_service_geo_serde.SuggestionToPB(value))
	}

	return connect.NewResponse(&geo_v1.SuggestAddressesResponse{Suggestions: values}), nil
}

func NewRPC(service service_suggest_addresses.Service) RPC {
	return RPC{service: service}
}
