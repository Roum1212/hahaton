package rpc_service_updates

import (
	"context"

	"connectrpc.com/connect"

	updates_v1 "github.com/Roum1212/hahaton/genproto/updates/v1"
	updates_v1connect "github.com/Roum1212/hahaton/genproto/updates/v1/updates_v1connect"
	rpc_service_updates_rpc_subscribe_updates "github.com/Roum1212/hahaton/internal/api/rpc/service/updates/rpc/subscribe-updates"
)

var _ updates_v1connect.UpdatesServiceHandler = Service{}

type Service struct {
	subscribeUpdatesRPC rpc_service_updates_rpc_subscribe_updates.RPC
}

func (x Service) SubscribeUpdates(
	ctx context.Context,
	r *connect.Request[updates_v1.SubscribeUpdatesRequest],
	stream *connect.ServerStream[updates_v1.SubscribeUpdatesResponse],
) error {
	return x.subscribeUpdatesRPC.SubscribeUpdates(ctx, r, stream)
}

func NewService(subscribeUpdatesRPC rpc_service_updates_rpc_subscribe_updates.RPC) Service {
	return Service{subscribeUpdatesRPC: subscribeUpdatesRPC}
}
