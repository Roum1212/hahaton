package rpc_service_event_chat_rpc_send_event_chat_message

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	rpc_error "go.flida.dev/rpc-error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	event_chat_v1 "github.com/Roum1212/hahaton/genproto/event-chat/v1"
	rpc_service_event_chat_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/event-chat/serde"
	service_send_event_chat_message "github.com/Roum1212/hahaton/internal/app/service/send-event-chat-message"
	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_event_chat "github.com/Roum1212/hahaton/internal/domain/aggregate/event-chat"
	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_send_event_chat_message.Service
}

func (x RPC) SendEventChatMessage(
	ctx context.Context,
	r *connect.Request[event_chat_v1.SendEventChatMessageRequest],
) (*connect.Response[event_chat_v1.SendEventChatMessageResponse], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

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

	text := r.Msg.GetText()
	requestID := r.Msg.GetRequestId()

	response, err := x.service.SendEventChatMessage(ctx, service_send_event_chat_message.NewRequest(
		session.GetUserID(),
		eventID,
		text,
		requestID,
	))
	if err != nil {
		switch {
		case errors.Is(err, service_send_event_chat_message.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_send_event_chat_message.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		case errors.Is(err, service_send_event_chat_message.ErrFailedPrecondition):
			return nil, rpc_error.NewFailedPreconditionError(nil)
		case errors.Is(err, service_send_event_chat_message.ErrAlreadyExists):
			return nil, rpc_error.NewAlreadyExistsError(nil)
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
		case errors.Is(err, aggregate_event_chat.ErrInvalidSequence):
			return nil, rpc_error.NewFailedPreconditionError(nil)
		case errors.Is(err, aggregate_event_chat.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "event_chat.v1.EventChat",
				ResourceName: "id:" + eventID.String(),
				Description:  "event chat not found",
			})
		case errors.Is(err, aggregate_volunteer.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "volunteer.v1.Volunteer",
				ResourceName: "id:" + session.GetUserID().String(),
				Description:  "volunteer not found",
			})
		default:
			return nil, rpc_error.NewInternalError()
		}
	}

	return connect.NewResponse(&event_chat_v1.SendEventChatMessageResponse{
		Message: rpc_service_event_chat_serde.MessageToPB(response.GetMessage()),
	}), nil
}

func NewRPC(service service_send_event_chat_message.Service) RPC {
	return RPC{service: service}
}
