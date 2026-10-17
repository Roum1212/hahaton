package rpc_service_event_chat_rpc_list_event_chat_messages

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	event_chat_v1 "github.com/Roum1212/hahaton/genproto/event-chat/v1"
	rpc_service_event_chat_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/event-chat/serde"
	service_list_event_chat_messages "github.com/Roum1212/hahaton/internal/app/service/list-event-chat-messages"
	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_event_chat "github.com/Roum1212/hahaton/internal/domain/aggregate/event-chat"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_list_event_chat_messages.Service
}

func (x RPC) ListEventChatMessages(
	ctx context.Context,
	r *connect.Request[event_chat_v1.ListEventChatMessagesRequest],
) (*connect.Response[event_chat_v1.ListEventChatMessagesResponse], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	var opts []service_list_event_chat_messages.RequestOption

	eventID, err := model_event_id.NewEventID(r.Msg.GetEventId())
	if err != nil {
		return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "event_id",
					Description: err.Error(),
				},
			},
		})
	}

	if value := r.Msg.SnapshotSequence; value != nil {
		opts = append(opts, service_list_event_chat_messages.WithSnapshotSequence(*value))
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

	response, err := x.service.ListEventChatMessages(
		ctx,
		service_list_event_chat_messages.NewRequest(session.GetUserID(), eventID, page, opts...),
	)
	if err != nil {
		switch {
		case errors.Is(err, service_list_event_chat_messages.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_list_event_chat_messages.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		case errors.Is(err, aggregate_event.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "event.v1.Event",
				ResourceName: "id:" + eventID.String(),
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
				ResourceName: "id:" + eventID.String(),
				Description:  "event chat not found",
			})
		default:
			return nil, rpc_error.NewInternalError()
		}
	}

	values := make([]*event_chat_v1.EventChatMessage, 0, len(response.GetMessages()))
	for _, value := range response.GetMessages() {
		values = append(values, rpc_service_event_chat_serde.MessageToPB(value))
	}

	return connect.NewResponse(&event_chat_v1.ListEventChatMessagesResponse{
		Messages:             values,
		TotalCount:           response.GetTotalCount(),
		SnapshotSequence:     response.GetSnapshotSequence(),
		VisibleAfterSequence: response.GetVisibleAfterSequence(),
	}), nil
}

func NewRPC(service service_list_event_chat_messages.Service) RPC {
	return RPC{service: service}
}
