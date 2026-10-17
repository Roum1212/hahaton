package rpc_service_participation_rpc_restore_participation

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"

	participation_v1 "github.com/Roum1212/hahaton/genproto/participation/v1"
	service_restore_participation "github.com/Roum1212/hahaton/internal/app/service/restore-participation"
	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_restore_participation.Service
}

func (x RPC) RestoreParticipation(
	ctx context.Context,
	r *connect.Request[participation_v1.RestoreParticipationRequest],
) (*connect.Response[emptypb.Empty], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	participationID, err := model_participation_id.NewParticipationID(r.Msg.GetId())
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

	if err := x.service.RestoreParticipation(ctx, service_restore_participation.NewRequest(
		session.GetUserID(),
		participationID,
	)); err != nil {
		switch {
		case errors.Is(err, service_restore_participation.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_restore_participation.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		case errors.Is(err, service_restore_participation.ErrFailedPrecondition):
			return nil, rpc_error.NewFailedPreconditionError(nil)
		case errors.Is(err, aggregate_event.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "event.v1.Event",
				ResourceName: "",
				Description:  "event not found",
			})
		case errors.Is(err, aggregate_participation.ErrInvalidState):
			return nil, rpc_error.NewFailedPreconditionError(nil)
		case errors.Is(err, aggregate_participation.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "participation.v1.Participation",
				ResourceName: "id:" + participationID.String(),
				Description:  "participation not found",
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

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func NewRPC(service service_restore_participation.Service) RPC {
	return RPC{service: service}
}
