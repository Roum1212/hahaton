package rpc_service_coordinator_rpc_list_coordinators

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	rpc_error "go.flida.dev/rpc-error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	coordinator_v1 "github.com/Roum1212/hahaton/genproto/coordinator/v1"
	rpc_service_coordinator_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/coordinator/serde"
	service_list_coordinators "github.com/Roum1212/hahaton/internal/app/service/list-coordinators"
	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_list_coordinators.Service
}

func (x RPC) ListCoordinators(
	ctx context.Context,
	r *connect.Request[coordinator_v1.ListCoordinatorsRequest],
) (*connect.Response[coordinator_v1.ListCoordinatorsResponse], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
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

	response, err := x.service.ListCoordinators(ctx, service_list_coordinators.NewRequest(session.GetUserID(), page))
	if err != nil {
		switch {
		case errors.Is(err, service_list_coordinators.ErrInvalidRequest):
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

	values := make([]*coordinator_v1.Coordinator, 0, len(response.GetCoordinators()))
	for _, value := range response.GetCoordinators() {
		values = append(values, rpc_service_coordinator_serde.AggregateToPBCoordinator(
			value,
			response.GetVolunteers()[value.GetUserID()],
		))
	}

	return connect.NewResponse(&coordinator_v1.ListCoordinatorsResponse{
		Coordinators: values,
		TotalCount:   response.GetTotalCount(),
	}), nil
}

func NewRPC(service service_list_coordinators.Service) RPC {
	return RPC{service: service}
}
