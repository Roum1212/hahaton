package rpc_service_event_rpc_complete_event

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	event_v1 "github.com/Roum1212/hahaton/genproto/event/v1"
	rpc_service_event_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/event/serde"
	service_complete_event "github.com/Roum1212/hahaton/internal/app/service/complete-event"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	aggregate_wallet "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_event_result "github.com/Roum1212/hahaton/internal/domain/model/event-result"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_complete_event.Service
}

func (x RPC) CompleteEvent(
	ctx context.Context,
	r *connect.Request[event_v1.CompleteEventRequest],
) (*connect.Response[event_v1.CompleteEventResponse], error) {
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

	results := make([]model_event_result.EventResult, 0, len(r.Msg.GetResults()))
	for index, value := range r.Msg.GetResults() {
		result, err := rpc_service_event_serde.ResultFromPB(value)
		if err != nil {
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       fmt.Sprintf("results[%d]", index),
						Description: err.Error(),
					},
				},
			})
		}

		results = append(results, result)
	}

	requestID := r.Msg.GetRequestId()

	response, err := x.service.CompleteEvent(ctx, service_complete_event.NewRequest(
		session.GetUserID(),
		eventID,
		results,
		requestID,
	))
	if err != nil {
		switch {
		case errors.Is(err, service_complete_event.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_complete_event.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		case errors.Is(err, service_complete_event.ErrAlreadyExists):
			return nil, rpc_error.NewAlreadyExistsError(nil)
		case errors.Is(err, aggregate_event.ErrAlreadyExists):
			return nil, rpc_error.NewAlreadyExistsError(nil)
		case errors.Is(err, aggregate_event.ErrInvalidState):
			return nil, rpc_error.NewFailedPreconditionError(nil)
		case errors.Is(err, aggregate_event.ErrNotFinished):
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
		case errors.Is(err, aggregate_wallet.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "wallet.Wallet",
				ResourceName: "",
				Description:  "wallet not found",
			})
		case errors.Is(err, aggregate_wallet.ErrOverflow):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		default:
			return nil, rpc_error.NewInternalError()
		}
	}

	return connect.NewResponse(&event_v1.CompleteEventResponse{Event: rpc_service_event_serde.AggregateToPBEvent(
		response.GetEvent(),
		0,
		0,
	)}), nil
}

func NewRPC(service service_complete_event.Service) RPC {
	return RPC{service: service}
}
