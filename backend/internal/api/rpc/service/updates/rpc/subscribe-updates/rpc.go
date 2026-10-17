package rpc_service_updates_rpc_subscribe_updates

import (
	"context"
	"errors"
	"net"
	"time"

	"connectrpc.com/connect"

	updates_v1 "github.com/Roum1212/hahaton/genproto/updates/v1"
	service_subscribe_updates "github.com/Roum1212/hahaton/internal/app/service/subscribe-updates"
	aggregate_session "github.com/Roum1212/hahaton/internal/domain/aggregate/session"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
	stream_controller "github.com/Roum1212/hahaton/internal/pkg/stream/controller"
	"github.com/Roum1212/hahaton/internal/pkg/updates"
)

const (
	heartbeatInterval = 10 * time.Second
	sendTimeout       = 5 * time.Second
)

type RPC struct {
	hub     *updates.Hub
	service service_subscribe_updates.Service
}

func (x RPC) SubscribeUpdates(
	ctx context.Context,
	r *connect.Request[updates_v1.SubscribeUpdatesRequest],
	stream *connect.ServerStream[updates_v1.SubscribeUpdatesResponse],
) error {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return rpc_error.NewUnauthenticatedError(nil)
	}

	changes, unsubscribe := x.hub.Subscribe(session.GetUserID())
	defer unsubscribe()

	stream.ResponseHeader().Set("X-Accel-Buffering", "no")

	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	kind := updates_v1.UpdateKind_UPDATE_KIND_READY

	for {
		if err := x.service.SubscribeUpdates(ctx,
			service_subscribe_updates.NewRequest(session.GetSessionID(), session.GetUserID()),
		); err != nil {
			switch {
			case errors.Is(err, aggregate_session.ErrInvalidCredentials):
				return rpc_error.NewUnauthenticatedError(nil)
			default:
				return rpc_error.NewInternalError()
			}
		}

		if err := x.send(ctx, stream, kind); err != nil {
			return err
		}

		kind = updates_v1.UpdateKind_UPDATE_KIND_CHANGED

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-changes:
		case <-ticker.C:
			kind = updates_v1.UpdateKind_UPDATE_KIND_HEARTBEAT
		}
	}
}

func (x RPC) send(
	ctx context.Context,
	stream *connect.ServerStream[updates_v1.SubscribeUpdatesResponse],
	kind updates_v1.UpdateKind,
) error {
	controller, ok := stream_controller.GetController(ctx)
	if ok {
		if err := controller.SetWriteDeadline(time.Now().Add(sendTimeout)); err != nil {
			return rpc_error.NewInternalError()
		}
	}

	if err := stream.Send(&updates_v1.SubscribeUpdatesResponse{Kind: kind}); err != nil {
		var networkError net.Error
		if errors.As(err, &networkError) && networkError.Timeout() {
			return rpc_error.NewResourceExhaustedError(nil)
		}

		return err
	}

	if ok {
		return controller.SetWriteDeadline(time.Time{})
	}

	return nil
}

func NewRPC(hub *updates.Hub, service service_subscribe_updates.Service) RPC {
	return RPC{hub: hub, service: service}
}
