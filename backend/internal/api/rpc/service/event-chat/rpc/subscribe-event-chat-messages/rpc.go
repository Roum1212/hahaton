package rpc_service_event_chat_rpc_subscribe_event_chat_messages

import (
	"context"
	"errors"
	"net"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	event_chat_v1 "github.com/Roum1212/hahaton/genproto/event-chat/v1"
	rpc_service_event_chat_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/event-chat/serde"
	service_subscribe_event_chat_messages "github.com/Roum1212/hahaton/internal/app/service/subscribe-event-chat-messages"
	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_event_chat "github.com/Roum1212/hahaton/internal/domain/aggregate/event-chat"
	aggregate_session "github.com/Roum1212/hahaton/internal/domain/aggregate/session"
	model_chat_message "github.com/Roum1212/hahaton/internal/domain/model/chat-message"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
	stream_controller "github.com/Roum1212/hahaton/internal/pkg/stream/controller"
)

const (
	pollInterval = 500 * time.Millisecond
	sendTimeout  = 5 * time.Second
)

type RPC struct {
	service service_subscribe_event_chat_messages.Service
}

func (x RPC) SubscribeEventChatMessages(
	ctx context.Context,
	r *connect.Request[event_chat_v1.SubscribeEventChatMessagesRequest],
	stream *connect.ServerStream[event_chat_v1.SubscribeEventChatMessagesResponse],
) error {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return rpc_error.NewUnauthenticatedError(nil)
	}

	eventID, err := model_event_id.NewEventID(r.Msg.GetEventId())
	if err != nil {
		return rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "event_id",
					Description: err.Error(),
				},
			},
		})
	}

	afterSequence := r.Msg.GetAfterSequence()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	var initialBoundary *uint64

	for {
		response, err := x.service.SubscribeEventChatMessages(ctx, service_subscribe_event_chat_messages.NewRequest(
			session.GetUserID(),
			eventID,
			session.GetSessionID(),
			afterSequence,
		))
		if err != nil {
			switch {
			case errors.Is(err, service_subscribe_event_chat_messages.ErrInvalidRequest):
				return rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequest_FieldViolation{
						{
							Field:       "request",
							Description: err.Error(),
						},
					},
				})
			case errors.Is(err, service_subscribe_event_chat_messages.ErrPermissionDenied):
				return rpc_error.NewPermissionDeniedError(nil)
			case errors.Is(err, aggregate_event.ErrNotFound):
				return rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
					ResourceType: "event.v1.Event",
					ResourceName: "id:" + eventID.String(),
					Description:  "event not found",
				})
			case errors.Is(err, aggregate_coordinator.ErrNotFound):
				return rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
					ResourceType: "coordinator.v1.Coordinator",
					ResourceName: "",
					Description:  "coordinator not found",
				})
			case errors.Is(err, aggregate_cooperation.ErrNotFound):
				return rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
					ResourceType: "cooperation.v1.Cooperation",
					ResourceName: "",
					Description:  "cooperation not found",
				})
			case errors.Is(err, aggregate_event_chat.ErrNotFound):
				return rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
					ResourceType: "event_chat.v1.EventChat",
					ResourceName: "id:" + eventID.String(),
					Description:  "event chat not found",
				})
			case errors.Is(err, aggregate_session.ErrInvalidCredentials):
				return rpc_error.NewUnauthenticatedError(nil)
			default:
				return rpc_error.NewInternalError()
			}
		}

		if initialBoundary != nil &&
			*initialBoundary != response.GetVisibleAfterSequence() {
			return rpc_error.NewPermissionDeniedError(nil)
		}

		boundary := response.GetVisibleAfterSequence()
		initialBoundary = &boundary

		for _, message := range response.GetMessages() {
			if err := x.send(ctx, stream, message); err != nil {
				return err
			}

			afterSequence = message.GetSequence()
		}

		if response.GetReadOnly() {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func (x RPC) send(
	ctx context.Context,
	stream *connect.ServerStream[event_chat_v1.SubscribeEventChatMessagesResponse],
	message model_chat_message.ChatMessage,
) error {
	controller, hasController := stream_controller.GetController(ctx)
	if hasController {
		if err := controller.SetWriteDeadline(time.Now().Add(sendTimeout)); err != nil {
			return rpc_error.NewInternalError()
		}
	}

	if err := stream.Send(&event_chat_v1.SubscribeEventChatMessagesResponse{
		Message: rpc_service_event_chat_serde.MessageToPB(message),
	}); err != nil {
		var networkError net.Error
		if errors.As(err, &networkError) && networkError.Timeout() {
			return rpc_error.NewResourceExhaustedError(nil)
		}

		return err
	}

	if hasController {
		if err := controller.SetWriteDeadline(time.Time{}); err != nil {
			return rpc_error.NewInternalError()
		}
	}

	return nil
}

func NewRPC(service service_subscribe_event_chat_messages.Service) RPC {
	return RPC{service: service}
}
