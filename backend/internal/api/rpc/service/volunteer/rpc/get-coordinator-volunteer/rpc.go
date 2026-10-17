package rpc_service_volunteer_rpc_get_coordinator_volunteer

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	volunteer_v1 "github.com/Roum1212/hahaton/genproto/volunteer/v1"
	rpc_service_volunteer_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/volunteer/serde"
	service_get_coordinator_volunteer "github.com/Roum1212/hahaton/internal/app/service/get-coordinator-volunteer"
	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_get_coordinator_volunteer.Service
}

func (x RPC) GetCoordinatorVolunteer(
	ctx context.Context,
	r *connect.Request[volunteer_v1.GetCoordinatorVolunteerRequest],
) (*connect.Response[volunteer_v1.GetCoordinatorVolunteerResponse], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	volunteerID, err := model_user_id.NewUserID(r.Msg.GetVolunteerId())
	if err != nil {
		return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "volunteer_id",
					Description: err.Error(),
				},
			},
		})
	}

	response, err := x.service.GetCoordinatorVolunteer(ctx, service_get_coordinator_volunteer.NewRequest(
		session.GetUserID(),
		volunteerID,
	))
	if err != nil {
		switch {
		case errors.Is(err, service_get_coordinator_volunteer.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_get_coordinator_volunteer.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		case errors.Is(err, aggregate_coordinator.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "coordinator.v1.Coordinator",
				ResourceName: "id:" + session.GetUserID().String(),
				Description:  "coordinator not found",
			})
		case errors.Is(err, aggregate_volunteer.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "volunteer.v1.Volunteer",
				ResourceName: "id:" + volunteerID.String(),
				Description:  "volunteer not found",
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

	return connect.NewResponse(&volunteer_v1.GetCoordinatorVolunteerResponse{
		Volunteer: rpc_service_volunteer_serde.CoordinatorVolunteerToPB(
			response.GetVolunteer(),
			response.GetShowContacts(),
		),
	}), nil
}

func NewRPC(service service_get_coordinator_volunteer.Service) RPC {
	return RPC{service: service}
}
