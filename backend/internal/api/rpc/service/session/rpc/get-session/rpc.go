package rpc_service_session_rpc_get_session

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"

	session_v1 "github.com/Roum1212/hahaton/genproto/session/v1"
	"github.com/Roum1212/hahaton/internal/api/rpc/service/session/serde"
	service_get_session "github.com/Roum1212/hahaton/internal/app/service/get-session"
	aggregate_session "github.com/Roum1212/hahaton/internal/domain/aggregate/session"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
)

type RPC struct {
	service service_get_session.Service
}

func (x RPC) GetSession(
	ctx context.Context,
	r *connect.Request[emptypb.Empty],
) (*connect.Response[session_v1.GetSessionResponse], error) {
	response, err := x.service.GetSession(ctx, service_get_session.NewRequest())
	if err != nil {
		switch {
		case errors.Is(err, service_get_session.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, aggregate_session.ErrInvalidCredentials):
			return nil, rpc_error.NewUnauthenticatedError(nil)
		default:
			return nil, rpc_error.NewInternalError()
		}
	}

	return connect.NewResponse(&session_v1.GetSessionResponse{
		Session: rpc_service_session_serde.AggregateToPBSession(
			response.GetSession(),
			response.HasCoordinatorProfile(),
		),
	}), nil
}

func NewRPC(service service_get_session.Service) RPC {
	return RPC{service: service}
}
