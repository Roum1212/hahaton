package rpc_service_event_rpc_create_event

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	rpc_error "go.flida.dev/rpc-error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	event_v1 "github.com/Roum1212/hahaton/genproto/event/v1"
	rpc_serde "github.com/Roum1212/hahaton/internal/api/rpc/serde"
	rpc_service_event_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/event/serde"
	rpc_service_geo_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/geo/serde"
	service_create_event "github.com/Roum1212/hahaton/internal/app/service/create-event"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	model_event_schedule "github.com/Roum1212/hahaton/internal/domain/model/event-schedule"
	model_event_title "github.com/Roum1212/hahaton/internal/domain/model/event-title"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	model_volunteer_limit "github.com/Roum1212/hahaton/internal/domain/model/volunteer-limit"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_create_event.Service
}

func (x RPC) CreateEvent(
	ctx context.Context,
	r *connect.Request[event_v1.CreateEventRequest],
) (*connect.Response[event_v1.CreateEventResponse], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	title, err := model_event_title.NewEventTitle(r.Msg.GetTitle())
	if err != nil {
		return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "title",
					Description: err.Error(),
				},
			},
		})
	}

	description := r.Msg.GetDescription()

	imagePath, err := model_image_path.NewImagePath(r.Msg.GetImagePath())
	if err != nil {
		return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "image_path",
					Description: err.Error(),
				},
			},
		})
	}

	location := r.Msg.GetLocation()

	point, err := rpc_service_geo_serde.PointFromPB(r.Msg.GetPoint())
	if err != nil {
		return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "point",
					Description: err.Error(),
				},
			},
		})
	}

	startsAt, err := rpc_serde.TimestampFromPB(r.Msg.GetStartsAt())
	if err != nil {
		return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "starts_at",
					Description: err.Error(),
				},
			},
		})
	}

	endsAt, err := rpc_serde.TimestampFromPB(r.Msg.GetEndsAt())
	if err != nil {
		return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "ends_at",
					Description: err.Error(),
				},
			},
		})
	}

	schedule, err := model_event_schedule.NewEventSchedule(startsAt, endsAt)
	if err != nil {
		return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "ends_at",
					Description: err.Error(),
				},
			},
		})
	}

	requirements, err := rpc_service_event_serde.RequirementsFromPB(r.Msg.GetCandidateRequirements())
	if err != nil {
		return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "candidate_requirements",
					Description: err.Error(),
				},
			},
		})
	}

	volunteerLimit, err := model_volunteer_limit.NewVolunteerLimit(r.Msg.GetVolunteerLimit())
	if err != nil {
		return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "volunteer_limit",
					Description: err.Error(),
				},
			},
		})
	}

	response, err := x.service.CreateEvent(ctx, service_create_event.NewRequest(
		session.GetUserID(),
		title,
		description,
		imagePath,
		location,
		point,
		schedule,
		requirements,
		volunteerLimit,
	))
	if err != nil {
		switch {
		case errors.Is(err, service_create_event.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_create_event.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		case errors.Is(err, service_create_event.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "storage.Image",
				ResourceName: "path:" + r.Msg.GetImagePath(),
				Description:  "image not found",
			})
		case errors.Is(err, aggregate_event.ErrAlreadyExists):
			return nil, rpc_error.NewAlreadyExistsError(nil)
		case errors.Is(err, aggregate_event.ErrInvalidSchedule):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, aggregate_event.ErrInvalidState):
			return nil, rpc_error.NewFailedPreconditionError(nil)
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

	return connect.NewResponse(&event_v1.CreateEventResponse{Id: response.GetEventID().String()}), nil
}

func NewRPC(service service_create_event.Service) RPC {
	return RPC{service: service}
}
