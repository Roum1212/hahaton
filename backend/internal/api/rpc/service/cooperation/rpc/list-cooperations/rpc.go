package rpc_service_cooperation_rpc_list_cooperations

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	cooperation_v1 "github.com/Roum1212/hahaton/genproto/cooperation/v1"
	rpc_service_cooperation_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/cooperation/serde"
	service_list_cooperations "github.com/Roum1212/hahaton/internal/app/service/list-cooperations"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_list_cooperations.Service
}

func (x RPC) ListCooperations(
	ctx context.Context,
	r *connect.Request[cooperation_v1.ListCooperationsRequest],
) (*connect.Response[cooperation_v1.ListCooperationsResponse], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	var opts []service_list_cooperations.RequestOption

	scope := int32(r.Msg.GetScope())

	if r.Msg.Status != nil {
		value, err := rpc_service_cooperation_serde.StatusFromPB(*r.Msg.Status)
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

		opts = append(opts, service_list_cooperations.WithStatus(value))
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

	response, err := x.service.ListCooperations(
		ctx,
		service_list_cooperations.NewRequest(session.GetUserID(), scope, page, opts...),
	)
	if err != nil {
		switch {
		case errors.Is(err, service_list_cooperations.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_list_cooperations.ErrPermissionDenied):
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

	values := make([]*cooperation_v1.Cooperation, 0, len(response.GetCooperations()))
	for _, value := range response.GetCooperations() {
		values = append(values, rpc_service_cooperation_serde.AggregateToPBCooperation(value))
	}

	return connect.NewResponse(&cooperation_v1.ListCooperationsResponse{
		Cooperations: values,
		TotalCount:   response.GetTotalCount(),
	}), nil
}

func NewRPC(service service_list_cooperations.Service) RPC {
	return RPC{service: service}
}
