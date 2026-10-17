package rpc_service_participation_rpc_list_participations

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	participation_v1 "github.com/Roum1212/hahaton/genproto/participation/v1"
	rpc_service_participation_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/participation/serde"
	service_list_participations "github.com/Roum1212/hahaton/internal/app/service/list-participations"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_list_participations.Service
}

func (x RPC) ListParticipations(
	ctx context.Context,
	r *connect.Request[participation_v1.ListParticipationsRequest],
) (*connect.Response[participation_v1.ListParticipationsResponse], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	var opts []service_list_participations.RequestOption

	scope := int32(r.Msg.GetScope())

	if r.Msg.EventId != nil {
		value, err := model_event_id.NewEventID(*r.Msg.EventId)
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

		opts = append(opts, service_list_participations.WithEventID(value))
	}

	if r.Msg.Status != nil {
		value, err := rpc_service_participation_serde.StatusFromPB(*r.Msg.Status)
		if err != nil {
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "status",
						Description: err.Error(),
					},
				},
			})
		}

		opts = append(opts, service_list_participations.WithStatus(value))
	}

	if r.Msg.VolunteerId != nil {
		value, err := model_user_id.NewUserID(*r.Msg.VolunteerId)
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

		opts = append(opts, service_list_participations.WithVolunteerID(value))
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

	response, err := x.service.ListParticipations(
		ctx,
		service_list_participations.NewRequest(session.GetUserID(), scope, page, opts...),
	)
	if err != nil {
		switch {
		case errors.Is(err, service_list_participations.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_list_participations.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
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

	values := make([]*participation_v1.Participation, 0, len(response.GetParticipations()))
	for _, value := range response.GetParticipations() {
		values = append(values, rpc_service_participation_serde.AggregateToPBParticipation(value))
	}

	return connect.NewResponse(&participation_v1.ListParticipationsResponse{
		Participations: values,
		TotalCount:     response.GetTotalCount(),
	}), nil
}

func NewRPC(service service_list_participations.Service) RPC {
	return RPC{service: service}
}
