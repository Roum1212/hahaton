package rpc_service_event_rpc_list_events

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	event_v1 "github.com/Roum1212/hahaton/genproto/event/v1"
	rpc_serde "github.com/Roum1212/hahaton/internal/api/rpc/serde"
	rpc_service_event_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/event/serde"
	rpc_service_geo_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/geo/serde"
	service_list_events "github.com/Roum1212/hahaton/internal/app/service/list-events"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_list_events.Service
}

func (x RPC) ListEvents(
	ctx context.Context,
	r *connect.Request[event_v1.ListEventsRequest],
) (*connect.Response[event_v1.ListEventsResponse], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	var opts []service_list_events.RequestOption

	scope := int32(r.Msg.GetScope())

	if r.Msg.CoordinatorId != nil {
		value, err := model_user_id.NewUserID(*r.Msg.CoordinatorId)
		if err != nil {
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "coordinator_id",
						Description: err.Error(),
					},
				},
			})
		}

		opts = append(opts, service_list_events.WithCoordinatorID(value))
	}

	if r.Msg.SearchArea != nil {
		value, err := rpc_service_geo_serde.SearchAreaFromPB(r.Msg.GetSearchArea())
		if err != nil {
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "search_area",
						Description: err.Error(),
					},
				},
			})
		}

		opts = append(opts, service_list_events.WithSearchArea(value))
	}

	if r.Msg.Status != nil {
		value, err := rpc_service_event_serde.StatusFromPB(*r.Msg.Status)
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

		opts = append(opts, service_list_events.WithStatus(value))
	}

	if r.Msg.StartsAtFrom != nil {
		value, err := rpc_serde.TimestampFromPB(r.Msg.GetStartsAtFrom())
		if err != nil {
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "starts_at_from",
						Description: err.Error(),
					},
				},
			})
		}

		opts = append(opts, service_list_events.WithStartsAtFrom(value))
	}

	if r.Msg.StartsAtTo != nil {
		value, err := rpc_serde.TimestampFromPB(r.Msg.GetStartsAtTo())
		if err != nil {
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "starts_at_to",
						Description: err.Error(),
					},
				},
			})
		}

		opts = append(opts, service_list_events.WithStartsAtTo(value))
	}

	if r.Msg.EndsAtFrom != nil {
		value, err := rpc_serde.TimestampFromPB(r.Msg.GetEndsAtFrom())
		if err != nil {
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "ends_at_from",
						Description: err.Error(),
					},
				},
			})
		}

		opts = append(opts, service_list_events.WithEndsAtFrom(value))
	}

	if r.Msg.EndsAtTo != nil {
		value, err := rpc_serde.TimestampFromPB(r.Msg.GetEndsAtTo())
		if err != nil {
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "ends_at_to",
						Description: err.Error(),
					},
				},
			})
		}

		opts = append(opts, service_list_events.WithEndsAtTo(value))
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

	response, err := x.service.ListEvents(
		ctx,
		service_list_events.NewRequest(session.GetUserID(), scope, page, opts...),
	)
	if err != nil {
		switch {
		case errors.Is(err, service_list_events.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_list_events.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		case errors.Is(err, service_list_events.ErrFailedPrecondition):
			return nil, rpc_error.NewFailedPreconditionError(nil)
		case errors.Is(err, service_list_events.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "geo.v1.City",
				ResourceName: "id:" + r.Msg.GetSearchArea().GetCityId(),
				Description:  "city not found",
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

	values := make([]*event_v1.Event, 0, len(response.GetEvents()))
	for _, value := range response.GetEvents() {
		values = append(values, rpc_service_event_serde.AggregateToPBEvent(
			value,
			response.GetApprovedCounts()[value.GetEventID()],
			response.GetPendingCounts()[value.GetEventID()],
		))
	}

	return connect.NewResponse(&event_v1.ListEventsResponse{
		Events:     values,
		TotalCount: response.GetTotalCount(),
	}), nil
}

func NewRPC(service service_list_events.Service) RPC {
	return RPC{service: service}
}
