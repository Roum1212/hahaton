package rpc_service_admin_rpc_force_complete_event

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	admin_v1 "github.com/Roum1212/hahaton/genproto/admin/v1"
	rpc_service_event_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/event/serde"
	service_complete_event "github.com/Roum1212/hahaton/internal/app/service/complete-event"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	aggregate_wallet "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
)

type RPC struct {
	service service_complete_event.Service
}

func (x RPC) ForceCompleteEvent(
	ctx context.Context,
	r *connect.Request[admin_v1.ForceCompleteEventRequest],
) (*connect.Response[admin_v1.ForceCompleteEventResponse], error) {
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

	response, err := x.service.CompleteEvent(ctx, service_complete_event.NewForceRequest(eventID))
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

	return connect.NewResponse(&admin_v1.ForceCompleteEventResponse{Event: rpc_service_event_serde.AggregateToPBEvent(
		response.GetEvent(),
		0,
		0,
	)}), nil
}

func NewRPC(service service_complete_event.Service) RPC {
	return RPC{service: service}
}
