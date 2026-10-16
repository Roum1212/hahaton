package rpc_service_session_rpc_authenticate_session

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"go.flida.dev/cookies/pkg/cookies"
	rpc_error "go.flida.dev/rpc-error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"

	session_v1 "github.com/Roum1212/hahaton/genproto/session/v1"
	service_authenticate_session "github.com/Roum1212/hahaton/internal/app/service/authenticate-session"
	aggregate_session "github.com/Roum1212/hahaton/internal/domain/aggregate/session"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	cookiesService cookies.Service
	service        service_authenticate_session.Service
}

func (x RPC) AuthenticateSession(
	ctx context.Context,
	r *connect.Request[session_v1.AuthenticateSessionRequest],
) (*connect.Response[emptypb.Empty], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	encryptedSessionToken, err := x.service.AuthenticateSession(
		ctx,
		service_authenticate_session.NewRequest(session.GetSessionID(), r.Msg.GetIdToken()),
	)
	if err != nil {
		switch {
		case errors.Is(err, service_authenticate_session.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "id_token",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, aggregate_session.ErrNotFound):
			return nil, rpc_error.NewUnauthenticatedError(nil)
		case errors.Is(err, aggregate_session.ErrInvalidCredentials):
			return nil, rpc_error.NewUnauthenticatedError(nil)
		case errors.Is(err, service_authenticate_session.ErrInvalidIDToken):
			return nil, rpc_error.NewUnauthenticatedError(nil)
		case errors.Is(err, aggregate_session.ErrAlreadyAuthenticated):
			return nil, rpc_error.NewFailedPreconditionError(&errdetails.PreconditionFailure{
				Violations: []*errdetails.PreconditionFailure_Violation{
					{
						Type:        "session",
						Description: "session is already authenticated",
					},
				},
			})
		case errors.Is(err, service_authenticate_session.ErrIncompleteProfile):
			return nil, rpc_error.NewFailedPreconditionError(&errdetails.PreconditionFailure{
				Violations: []*errdetails.PreconditionFailure_Violation{
					{
						Type:        "profile",
						Description: "identity profile is incomplete",
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

func NewRPC(cookiesService cookies.Service, service service_authenticate_session.Service) RPC {
	return RPC{
		cookiesService: cookiesService,
		service:        service,
	}
}
