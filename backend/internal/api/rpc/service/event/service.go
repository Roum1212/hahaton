package rpc_service_event

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	event_v1 "github.com/Roum1212/hahaton/genproto/event/v1"
	event_v1connect "github.com/Roum1212/hahaton/genproto/event/v1/event_v1connect"
	rpc_service_event_rpc_cancel_event "github.com/Roum1212/hahaton/internal/api/rpc/service/event/rpc/cancel-event"
	rpc_service_event_rpc_complete_event "github.com/Roum1212/hahaton/internal/api/rpc/service/event/rpc/complete-event"
	rpc_service_event_rpc_create_event "github.com/Roum1212/hahaton/internal/api/rpc/service/event/rpc/create-event"
	rpc_service_event_rpc_get_event "github.com/Roum1212/hahaton/internal/api/rpc/service/event/rpc/get-event"
	rpc_service_event_rpc_list_events "github.com/Roum1212/hahaton/internal/api/rpc/service/event/rpc/list-events"
	rpc_service_event_rpc_update_event "github.com/Roum1212/hahaton/internal/api/rpc/service/event/rpc/update-event"
)

var _ event_v1connect.EventServiceHandler = Service{}

type Service struct {
	cancelEventRPC   rpc_service_event_rpc_cancel_event.RPC
	completeEventRPC rpc_service_event_rpc_complete_event.RPC
	createEventRPC   rpc_service_event_rpc_create_event.RPC
	getEventRPC      rpc_service_event_rpc_get_event.RPC
	listEventsRPC    rpc_service_event_rpc_list_events.RPC
	updateEventRPC   rpc_service_event_rpc_update_event.RPC
}

func (x Service) CancelEvent(
	ctx context.Context,
	r *connect.Request[event_v1.CancelEventRequest],
) (*connect.Response[emptypb.Empty], error) {
	return x.cancelEventRPC.CancelEvent(ctx, r)
}

func (x Service) CompleteEvent(
	ctx context.Context,
	r *connect.Request[event_v1.CompleteEventRequest],
) (*connect.Response[event_v1.CompleteEventResponse], error) {
	return x.completeEventRPC.CompleteEvent(ctx, r)
}

func (x Service) CreateEvent(
	ctx context.Context,
	r *connect.Request[event_v1.CreateEventRequest],
) (*connect.Response[event_v1.CreateEventResponse], error) {
	return x.createEventRPC.CreateEvent(ctx, r)
}

func (x Service) GetEvent(
	ctx context.Context,
	r *connect.Request[event_v1.GetEventRequest],
) (*connect.Response[event_v1.GetEventResponse], error) {
	return x.getEventRPC.GetEvent(ctx, r)
}

func (x Service) ListEvents(
	ctx context.Context,
	r *connect.Request[event_v1.ListEventsRequest],
) (*connect.Response[event_v1.ListEventsResponse], error) {
	return x.listEventsRPC.ListEvents(ctx, r)
}

func (x Service) UpdateEvent(
	ctx context.Context,
	r *connect.Request[event_v1.UpdateEventRequest],
) (*connect.Response[emptypb.Empty], error) {
	return x.updateEventRPC.UpdateEvent(ctx, r)
}

func NewService(
	cancelEventRPC rpc_service_event_rpc_cancel_event.RPC,
	completeEventRPC rpc_service_event_rpc_complete_event.RPC,
	createEventRPC rpc_service_event_rpc_create_event.RPC,
	getEventRPC rpc_service_event_rpc_get_event.RPC,
	listEventsRPC rpc_service_event_rpc_list_events.RPC,
	updateEventRPC rpc_service_event_rpc_update_event.RPC,
) Service {
	return Service{
		cancelEventRPC:   cancelEventRPC,
		completeEventRPC: completeEventRPC,
		createEventRPC:   createEventRPC,
		getEventRPC:      getEventRPC,
		listEventsRPC:    listEventsRPC,
		updateEventRPC:   updateEventRPC,
	}
}
