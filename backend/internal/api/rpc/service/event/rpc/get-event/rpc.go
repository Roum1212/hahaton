package rpc_service_event_rpc_get_event

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	event_v1 "github.com/Roum1212/hahaton/genproto/event/v1"
	rpc_service_event_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/event/serde"
	service_get_event "github.com/Roum1212/hahaton/internal/app/service/get-event"
	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_get_event.Service
}

func (x RPC) GetEvent(
	ctx context.Context,
	r *connect.Request[event_v1.GetEventRequest],
) (*connect.Response[event_v1.GetEventResponse], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	eventID, err := model_event_id.NewEventID(r.Msg.GetId())
	if err != nil {
		return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "id",
					Description: err.Error(),
				},
			},
		})
	}

	response, err := x.service.GetEvent(ctx, service_get_event.NewRequest(session.GetUserID(), eventID))
	if err != nil {
		switch {
		case errors.Is(err, service_get_event.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_get_event.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		case errors.Is(err, aggregate_event.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "event.v1.Event",
				ResourceName: "id:" + eventID.String(),
				Description:  "event not found",
			})
		case errors.Is(err, aggregate_cooperation.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "cooperation.v1.Cooperation",
				ResourceName: "",
				Description:  "cooperation not found",
			})
		default:
			return nil, rpc_error.NewInternalError()
		}
	}

	return connect.NewResponse(&event_v1.GetEventResponse{Event: rpc_service_event_serde.AggregateToPBEvent(
		response.GetEvent(),
		response.GetApprovedCount(),
		response.GetPendingCount(),
	)}), nil
}

func NewRPC(service service_get_event.Service) RPC {
	return RPC{service: service}
}
