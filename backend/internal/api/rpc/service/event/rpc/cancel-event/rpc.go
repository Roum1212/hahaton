package rpc_service_event_rpc_cancel_event

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"

	event_v1 "github.com/Roum1212/hahaton/genproto/event/v1"
	service_cancel_event "github.com/Roum1212/hahaton/internal/app/service/cancel-event"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_cancel_event.Service
}

func (x RPC) CancelEvent(
	ctx context.Context,
	r *connect.Request[event_v1.CancelEventRequest],
) (*connect.Response[emptypb.Empty], error) {
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

	if err := x.service.CancelEvent(ctx, service_cancel_event.NewRequest(session.GetUserID(), eventID)); err != nil {
		switch {
		case errors.Is(err, service_cancel_event.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_cancel_event.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		case errors.Is(err, aggregate_event.ErrInvalidState):
			return nil, rpc_error.NewFailedPreconditionError(nil)
		case errors.Is(err, aggregate_event.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "event.v1.Event",
				ResourceName: "id:" + eventID.String(),
				Description:  "event not found",
			})
		case errors.Is(err, aggregate_coordinator.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "coordinator.v1.Coordinator",
				ResourceName: "id:" + session.GetUserID().String(),
				Description:  "coordinator not found",
			})
		case errors.Is(err, aggregate_participation.ErrInvalidState):
			return nil, rpc_error.NewFailedPreconditionError(nil)
		default:
			return nil, rpc_error.NewInternalError()
		}
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func NewRPC(service service_cancel_event.Service) RPC {
	return RPC{service: service}
}
