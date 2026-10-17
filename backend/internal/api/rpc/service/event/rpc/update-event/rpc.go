package rpc_service_event_rpc_update_event

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	rpc_error "go.flida.dev/rpc-error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"

	event_v1 "github.com/Roum1212/hahaton/genproto/event/v1"
	rpc_serde "github.com/Roum1212/hahaton/internal/api/rpc/serde"
	rpc_service_event_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/event/serde"
	rpc_service_geo_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/geo/serde"
	service_update_event "github.com/Roum1212/hahaton/internal/app/service/update-event"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_event_title "github.com/Roum1212/hahaton/internal/domain/model/event-title"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	model_volunteer_limit "github.com/Roum1212/hahaton/internal/domain/model/volunteer-limit"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_update_event.Service
}

func (x RPC) UpdateEvent(
	ctx context.Context,
	r *connect.Request[event_v1.UpdateEventRequest],
) (*connect.Response[emptypb.Empty], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	var opts []service_update_event.RequestOption

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

	if value := r.Msg.Description; value != nil {
		opts = append(opts, service_update_event.WithDescription(*value))
	}

	if value := r.Msg.ImagePath; value != nil {
		path, err := model_image_path.NewImagePath(*value)
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

		opts = append(opts, service_update_event.WithImagePath(path))
	}

	if r.Msg.Title != nil {
		value, err := model_event_title.NewEventTitle(*r.Msg.Title)
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

		opts = append(opts, service_update_event.WithTitle(value))
	}

	if value := r.Msg.Location; value != nil {
		opts = append(opts, service_update_event.WithLocation(*value))
	}

	if r.Msg.Point != nil {
		value, err := rpc_service_geo_serde.PointFromPB(r.Msg.GetPoint())
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

		opts = append(opts, service_update_event.WithPoint(value))
	}

	if r.Msg.StartsAt != nil {
		value, err := rpc_serde.TimestampFromPB(r.Msg.GetStartsAt())
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

		opts = append(opts, service_update_event.WithStartsAt(value))
	}

	if r.Msg.EndsAt != nil {
		value, err := rpc_serde.TimestampFromPB(r.Msg.GetEndsAt())
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

		opts = append(opts, service_update_event.WithEndsAt(value))
	}

	if r.Msg.CandidateRequirements != nil {
		value, err := rpc_service_event_serde.RequirementsFromPB(r.Msg.GetCandidateRequirements().GetValues())
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

		opts = append(opts, service_update_event.WithRequirements(value))
	}

	if r.Msg.VolunteerLimit != nil {
		value, err := model_volunteer_limit.NewVolunteerLimit(*r.Msg.VolunteerLimit)
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

		opts = append(opts, service_update_event.WithVolunteerLimit(value))
	}

	if err := x.service.UpdateEvent(
		ctx,
		service_update_event.NewRequest(session.GetUserID(), eventID, opts...),
	); err != nil {
		switch {
		case errors.Is(err, service_update_event.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_update_event.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		case errors.Is(err, service_update_event.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "storage.Image",
				ResourceName: "path:" + r.Msg.GetImagePath(),
				Description:  "image not found",
			})
		case errors.Is(err, aggregate_event.ErrInvalidSchedule):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, aggregate_event.ErrNotEditable):
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
		default:
			return nil, rpc_error.NewInternalError()
		}
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func NewRPC(service service_update_event.Service) RPC {
	return RPC{service: service}
}
