package rpc_service_volunteer_rpc_get_volunteer

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	volunteer_v1 "github.com/Roum1212/hahaton/genproto/volunteer/v1"
	rpc_service_volunteer_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/volunteer/serde"
	service_get_volunteer "github.com/Roum1212/hahaton/internal/app/service/get-volunteer"
	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
	aggregate_wallet "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_get_volunteer.Service
}

func (x RPC) GetVolunteer(
	ctx context.Context,
	r *connect.Request[volunteer_v1.GetVolunteerRequest],
) (*connect.Response[volunteer_v1.GetVolunteerResponse], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	var opts []service_get_volunteer.RequestOption

	if r.Msg.Id != nil {
		value, err := model_user_id.NewUserID(*r.Msg.Id)
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

		opts = append(opts, service_get_volunteer.WithID(value))
	}

	response, err := x.service.GetVolunteer(ctx, service_get_volunteer.NewRequest(session.GetUserID(), opts...))
	if err != nil {
		switch {
		case errors.Is(err, service_get_volunteer.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_get_volunteer.ErrPermissionDenied):
			return nil, rpc_error.NewPermissionDeniedError(nil)
		case errors.Is(err, aggregate_volunteer.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "volunteer.v1.Volunteer",
				ResourceName: "id:" + session.GetUserID().String(),
				Description:  "volunteer not found",
			})
		case errors.Is(err, aggregate_wallet.ErrNotFound):
			return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
				ResourceType: "wallet.Wallet",
				ResourceName: "user_id:" + session.GetUserID().String(),
				Description:  "wallet not found",
			})
		default:
			return nil, rpc_error.NewInternalError()
		}
	}

	return connect.NewResponse(&volunteer_v1.GetVolunteerResponse{
		Volunteer: rpc_service_volunteer_serde.AggregateToPBVolunteer(response.GetVolunteer(), response.GetWallet()),
	}), nil
}

func NewRPC(service service_get_volunteer.Service) RPC {
	return RPC{service: service}
}
