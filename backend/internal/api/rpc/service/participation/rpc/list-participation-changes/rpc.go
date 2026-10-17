package rpc_service_participation_rpc_list_participation_changes

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	participation_v1 "github.com/Roum1212/hahaton/genproto/participation/v1"
	rpc_service_participation_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/participation/serde"
	service_list_participation_changes "github.com/Roum1212/hahaton/internal/app/service/list-participation-changes"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_list_participation_changes.Service
}

func (x RPC) ListParticipationChanges(
	ctx context.Context,
	r *connect.Request[participation_v1.ListParticipationChangesRequest],
) (*connect.Response[participation_v1.ListParticipationChangesResponse], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	participationID, err := model_participation_id.NewParticipationID(r.Msg.GetParticipationId())
	if err != nil {
		return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "participation_id",
					Description: err.Error(),
				},
			},
		})
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

	response, err := x.service.ListParticipationChanges(ctx, service_list_participation_changes.NewRequest(
		session.GetUserID(),
		participationID,
		page,
	))
	if err != nil {
		switch {
		case errors.Is(err, service_list_participation_changes.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_list_participation_changes.ErrPermissionDenied):
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

	values := make([]*participation_v1.ParticipationChange, 0, len(response.GetChanges()))
	for _, value := range response.GetChanges() {
		values = append(values, rpc_service_participation_serde.ChangeToPB(value))
	}

	return connect.NewResponse(&participation_v1.ListParticipationChangesResponse{
		Changes:    values,
		TotalCount: response.GetTotalCount(),
	}), nil
}

func NewRPC(service service_list_participation_changes.Service) RPC {
	return RPC{service: service}
}
