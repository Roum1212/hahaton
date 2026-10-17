package rpc_service_volunteer_rpc_list_volunteer_operations

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	rpc_error "go.flida.dev/rpc-error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	volunteer_v1 "github.com/Roum1212/hahaton/genproto/volunteer/v1"
	rpc_service_volunteer_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/volunteer/serde"
	service_list_volunteer_operations "github.com/Roum1212/hahaton/internal/app/service/list-volunteer-operations"
	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_list_volunteer_operations.Service
}

func (x RPC) ListVolunteerOperations(
	ctx context.Context,
	r *connect.Request[volunteer_v1.ListVolunteerOperationsRequest],
) (*connect.Response[volunteer_v1.ListVolunteerOperationsResponse], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	var opts []service_list_volunteer_operations.RequestOption

	if r.Msg.Type != nil {
		value, err := rpc_service_volunteer_serde.OperationTypeFromPB(*r.Msg.Type)
		if err != nil {
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "type",
						Description: err.Error(),
					},
				},
			})
		}

		opts = append(opts, service_list_volunteer_operations.WithOperationType(value))
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

	response, err := x.service.ListVolunteerOperations(
		ctx,
		service_list_volunteer_operations.NewRequest(session.GetUserID(), page, opts...),
	)
	if err != nil {
		switch {
		case errors.Is(err, service_list_volunteer_operations.ErrInvalidRequest):
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

	values := make([]*volunteer_v1.VolunteerOperation, 0, len(response.GetOperations()))
	for _, value := range response.GetOperations() {
		values = append(values, rpc_service_volunteer_serde.OperationToPB(value))
	}

	return connect.NewResponse(&volunteer_v1.ListVolunteerOperationsResponse{
		Operations: values,
		TotalCount: response.GetTotalCount(),
	}), nil
}

func NewRPC(service service_list_volunteer_operations.Service) RPC {
	return RPC{service: service}
}
