package rpc_service_event_chat_rpc_list_my_event_chats

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	rpc_error "go.flida.dev/rpc-error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	event_chat_v1 "github.com/Roum1212/hahaton/genproto/event-chat/v1"
	rpc_service_event_chat_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/event-chat/serde"
	service_list_my_event_chats "github.com/Roum1212/hahaton/internal/app/service/list-my-event-chats"
	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_event_chat "github.com/Roum1212/hahaton/internal/domain/aggregate/event-chat"
	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_list_my_event_chats.Service
}

func (x RPC) ListMyEventChats(
	ctx context.Context,
	r *connect.Request[event_chat_v1.ListMyEventChatsRequest],
) (*connect.Response[event_chat_v1.ListMyEventChatsResponse], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	page, err := model_page.NewPage(r.Msg.GetLimit(), r.Msg.GetOffset())
	if err != nil {
		return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "pagination",
					Description: err.Error(),
				},
			},
		})
	}

	response, err := x.service.ListMyEventChats(ctx, service_list_my_event_chats.NewRequest(session.GetUserID(), page))
	if err != nil {
		switch {
		case errors.Is(err, service_list_my_event_chats.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_list_my_event_chats.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		case errors.Is(err, aggregate_event.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "event.v1.Event",
				ResourceName: "",
				Description:  "event not found",
			})
		case errors.Is(err, aggregate_coordinator.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "coordinator.v1.Coordinator",
				ResourceName: "",
				Description:  "coordinator not found",
			})
		case errors.Is(err, aggregate_cooperation.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "cooperation.v1.Cooperation",
				ResourceName: "",
				Description:  "cooperation not found",
			})
		case errors.Is(err, aggregate_event_chat.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "event_chat.v1.EventChat",
				ResourceName: "",
				Description:  "event chat not found",
			})
		default:
			return nil, rpc_error.NewInternalError()
		}
	}

	values := make([]*event_chat_v1.EventChat, 0, len(response.GetEvents()))
	for _, value := range response.GetEvents() {
		values = append(values, rpc_service_event_chat_serde.AggregateToPBEventChat(
			value,
			response.GetChats()[value.GetEventID()],
			response.GetBoundaries()[value.GetEventID()],
		))
	}

	return connect.NewResponse(&event_chat_v1.ListMyEventChatsResponse{
		Chats:      values,
		TotalCount: response.GetTotalCount(),
	}), nil
}

func NewRPC(service service_list_my_event_chats.Service) RPC {
	return RPC{service: service}
}
