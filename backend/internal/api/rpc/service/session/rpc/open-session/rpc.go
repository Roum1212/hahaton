package rpc_service_session_rpc_open_session

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"go.flida.dev/cookies/pkg/cookies"
	rpc_error "go.flida.dev/rpc-error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"

	service_open_session "github.com/Roum1212/hahaton/internal/app/service/open-session"
)

type RPC struct {
	cookiesService cookies.Service
	service        service_open_session.Service
}

func (x RPC) OpenSession(
	ctx context.Context,
	r *connect.Request[emptypb.Empty],
) (*connect.Response[emptypb.Empty], error) {
	encryptedSessionToken, err := x.service.OpenSession(ctx, service_open_session.NewRequest())
	if err != nil {
		switch {
		case errors.Is(err, service_open_session.ErrInvalidRequest):
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
	x.cookiesService.SetCookie(response.Header(), encryptedSessionToken.String())

	return response, nil
}

func NewRPC(cookiesService cookies.Service, service service_open_session.Service) RPC {
	return RPC{
		cookiesService: cookiesService,
		service:        service,
	}
}
