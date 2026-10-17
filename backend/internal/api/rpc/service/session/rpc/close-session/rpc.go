package rpc_service_session_rpc_close_session

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"

	service_close_session "github.com/Roum1212/hahaton/internal/app/service/close-session"
	"github.com/Roum1212/hahaton/internal/pkg/cookies"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
)

type RPC struct {
	expiredCookiesService cookies.Service
	service               service_close_session.Service
}

func (x RPC) CloseSession(
	ctx context.Context,
	r *connect.Request[emptypb.Empty],
) (*connect.Response[emptypb.Empty], error) {
	if err := x.service.CloseSession(ctx, service_close_session.NewRequest()); err != nil {
		switch {
		case errors.Is(err, service_close_session.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		default:
			return nil, rpc_error.NewInternalError()
		}
	}

	response := connect.NewResponse(&emptypb.Empty{})
	x.expiredCookiesService.SetCookie(response.Header(), "")

	return response, nil
}

func NewRPC(expiredCookiesService cookies.Service, service service_close_session.Service) RPC {
	return RPC{
		expiredCookiesService: expiredCookiesService,
		service:               service,
	}
}
