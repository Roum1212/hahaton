package rpc_service_participation_rpc_request_participation

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	participation_v1 "github.com/Roum1212/hahaton/genproto/participation/v1"
	service_request_participation "github.com/Roum1212/hahaton/internal/app/service/request-participation"
	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_request_participation.Service
}

func (x RPC) RequestParticipation(
	ctx context.Context,
	r *connect.Request[participation_v1.RequestParticipationRequest],
) (*connect.Response[participation_v1.RequestParticipationResponse], error) {
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

	response, err := x.service.RequestParticipation(ctx, service_request_participation.NewRequest(
		session.GetUserID(),
		eventID,
	))
	if err != nil {
		switch {
		case errors.Is(err, service_request_participation.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_request_participation.ErrFailedPrecondition):
			return nil, rpc_error.NewFailedPreconditionError(nil)
		case errors.Is(err, aggregate_event.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "event.v1.Event",
				ResourceName: "id:" + eventID.String(),
				Description:  "event not found",
			})
		case errors.Is(err, aggregate_participation.ErrAlreadyExists):
			return nil, rpc_error.NewAlreadyExistsError(nil)
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

	return connect.NewResponse(&participation_v1.RequestParticipationResponse{
		Id: response.GetParticipationID().String(),
	}), nil
}

func NewRPC(service service_request_participation.Service) RPC {
	return RPC{service: service}
}
