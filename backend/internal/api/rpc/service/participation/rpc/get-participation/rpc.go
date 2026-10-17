package rpc_service_participation_rpc_get_participation

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	rpc_error "go.flida.dev/rpc-error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	participation_v1 "github.com/Roum1212/hahaton/genproto/participation/v1"
	rpc_service_participation_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/participation/serde"
	service_get_participation "github.com/Roum1212/hahaton/internal/app/service/get-participation"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_get_participation.Service
}

func (x RPC) GetParticipation(
	ctx context.Context,
	r *connect.Request[participation_v1.GetParticipationRequest],
) (*connect.Response[participation_v1.GetParticipationResponse], error) {
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

	response, err := x.service.GetParticipation(ctx, service_get_participation.NewRequest(
		session.GetUserID(),
		participationID,
	))
	if err != nil {
		switch {
		case errors.Is(err, service_get_participation.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_get_participation.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		case errors.Is(err, aggregate_participation.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "participation.v1.Participation",
				ResourceName: "id:" + participationID.String(),
				Description:  "participation not found",
			})
		case errors.Is(err, aggregate_coordinator.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "coordinator.v1.Coordinator",
				ResourceName: "id:" + session.GetUserID().String(),
				Description:  "coordinator not found",
			})
		default:
			return nil, rpc_error.NewInternalError()
		}
	}

	return connect.NewResponse(&participation_v1.GetParticipationResponse{
		Participation: rpc_service_participation_serde.AggregateToPBParticipation(response.GetParticipation()),
	}), nil
}

func NewRPC(service service_get_participation.Service) RPC {
	return RPC{service: service}
}
