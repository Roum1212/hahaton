package rpc_service_admin_rpc_find_volunteers

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	admin_v1 "github.com/Roum1212/hahaton/genproto/admin/v1"
	rpc_service_admin_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/admin/serde"
	service_find_volunteers "github.com/Roum1212/hahaton/internal/app/service/find-volunteers"
	model_email_address "github.com/Roum1212/hahaton/internal/domain/model/email-address"
	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
	model_phone_number "github.com/Roum1212/hahaton/internal/domain/model/phone-number"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
)

type RPC struct {
	service service_find_volunteers.Service
}

func (x RPC) FindVolunteers(
	ctx context.Context,
	r *connect.Request[admin_v1.FindVolunteersRequest],
) (*connect.Response[admin_v1.FindVolunteersResponse], error) {
	var opts []service_find_volunteers.RequestOption

	switch criterion := r.Msg.GetCriterion().(type) {
	case *admin_v1.FindVolunteersRequest_VolunteerId:
		value, err := model_user_id.NewUserID(criterion.VolunteerId)
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

		opts = append(opts, service_find_volunteers.WithVolunteerID(value))
	case *admin_v1.FindVolunteersRequest_EmailAddress:
		value, err := model_email_address.NewEmailAddress(criterion.EmailAddress)
		if err != nil {
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "email_address",
						Description: err.Error(),
					},
				},
			})
		}

		opts = append(opts, service_find_volunteers.WithEmailAddress(value))
	case *admin_v1.FindVolunteersRequest_PhoneNumber:
		value, err := model_phone_number.NewPhoneNumber(criterion.PhoneNumber)
		if err != nil {
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "phone_number",
						Description: err.Error(),
					},
				},
			})
		}

		opts = append(opts, service_find_volunteers.WithPhoneNumber(value))
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

	response, err := x.service.FindVolunteers(ctx, service_find_volunteers.NewRequest(page, opts...))
	if err != nil {
		switch {
		case errors.Is(err, service_find_volunteers.ErrInvalidRequest):
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

	values := make([]*admin_v1.AdminVolunteer, 0, len(response.GetVolunteers()))
	for _, value := range response.GetVolunteers() {
		values = append(values, rpc_service_admin_serde.VolunteerToPB(
			value,
			response.GetHasCoordinatorProfiles()[value.GetUserID()],
		))
	}

	return connect.NewResponse(&admin_v1.FindVolunteersResponse{
		Volunteers: values,
		TotalCount: response.GetTotalCount(),
	}), nil
}

func NewRPC(service service_find_volunteers.Service) RPC {
	return RPC{service: service}
}
