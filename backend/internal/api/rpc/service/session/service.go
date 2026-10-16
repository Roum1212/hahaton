package rpc_service_session

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	session_v1 "github.com/Roum1212/hahaton/genproto/session/v1"
	"github.com/Roum1212/hahaton/genproto/session/v1/session_v1connect"
	"github.com/Roum1212/hahaton/internal/api/rpc/service/session/rpc/authenticate-session"
	"github.com/Roum1212/hahaton/internal/api/rpc/service/session/rpc/close-session"
	"github.com/Roum1212/hahaton/internal/api/rpc/service/session/rpc/get-session"
	"github.com/Roum1212/hahaton/internal/api/rpc/service/session/rpc/open-session"
)

var _ session_v1connect.SessionServiceHandler = Service{}

type Service struct {
	authenticateSessionRPC rpc_service_session_rpc_authenticate_session.RPC
	closeSessionRPC        rpc_service_session_rpc_close_session.RPC
	getSessionRPC          rpc_service_session_rpc_get_session.RPC
	openSessionRPC         rpc_service_session_rpc_open_session.RPC
}

func (x Service) AuthenticateSession(
	ctx context.Context,
	r *connect.Request[session_v1.AuthenticateSessionRequest],
) (*connect.Response[emptypb.Empty], error) {
	return x.authenticateSessionRPC.AuthenticateSession(ctx, r)
}

func (x Service) CloseSession(
	ctx context.Context,
	r *connect.Request[emptypb.Empty],
) (*connect.Response[emptypb.Empty], error) {
	return x.closeSessionRPC.CloseSession(ctx, r)
}

func (x Service) GetSession(
	ctx context.Context,
	r *connect.Request[emptypb.Empty],
) (*connect.Response[session_v1.GetSessionResponse], error) {
	return x.getSessionRPC.GetSession(ctx, r)
}

func (x Service) OpenSession(
	ctx context.Context,
	r *connect.Request[emptypb.Empty],
) (*connect.Response[emptypb.Empty], error) {
	return x.openSessionRPC.OpenSession(ctx, r)
}

func NewService(
	authenticateSessionRPC rpc_service_session_rpc_authenticate_session.RPC,
	closeSessionRPC rpc_service_session_rpc_close_session.RPC,
	getSessionRPC rpc_service_session_rpc_get_session.RPC,
	openSessionRPC rpc_service_session_rpc_open_session.RPC,
) Service {
	return Service{
		authenticateSessionRPC: authenticateSessionRPC,
		closeSessionRPC:        closeSessionRPC,
		getSessionRPC:          getSessionRPC,
		openSessionRPC:         openSessionRPC,
	}
}
