package rpc_service_reward_order_rpc_list_reward_orders

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	reward_order_v1 "github.com/Roum1212/hahaton/genproto/reward-order/v1"
	rpc_service_reward_order_serde "github.com/Roum1212/hahaton/internal/api/rpc/service/reward-order/serde"
	service_list_reward_orders "github.com/Roum1212/hahaton/internal/app/service/list-reward-orders"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
	model_reward_id "github.com/Roum1212/hahaton/internal/domain/model/reward-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
)

type RPC struct {
	service service_list_reward_orders.Service
}

func (x RPC) ListRewardOrders(
	ctx context.Context,
	r *connect.Request[reward_order_v1.ListRewardOrdersRequest],
) (*connect.Response[reward_order_v1.ListRewardOrdersResponse], error) {
	session, ok := session_context.GetSession(ctx)
	if !ok {
		return nil, rpc_error.NewUnauthenticatedError(nil)
	}

	var opts []service_list_reward_orders.RequestOption

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

		opts = append(opts, service_list_reward_orders.WithCoordinatorID(value))
	}

	if r.Msg.RewardId != nil {
		value, err := model_reward_id.NewRewardID(*r.Msg.RewardId)
		if err != nil {
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "reward_id",
						Description: err.Error(),
					},
				},
			})
		}

		opts = append(opts, service_list_reward_orders.WithRewardID(value))
	}

	if r.Msg.Status != nil {
		value, err := rpc_service_reward_order_serde.StatusFromPB(*r.Msg.Status)
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

		opts = append(opts, service_list_reward_orders.WithStatus(value))
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

		opts = append(opts, service_list_reward_orders.WithVolunteerID(value))
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

	response, err := x.service.ListRewardOrders(
		ctx,
		service_list_reward_orders.NewRequest(session.GetUserID(), scope, page, opts...),
	)
	if err != nil {
		switch {
		case errors.Is(err, service_list_reward_orders.ErrInvalidRequest):
			return nil, rpc_error.NewInvalidArgumentError(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "request",
						Description: err.Error(),
					},
				},
			})
		case errors.Is(err, service_list_reward_orders.ErrPermissionDenied):
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

	values := make([]*reward_order_v1.RewardOrder, 0, len(response.GetOrders()))
	for _, value := range response.GetOrders() {
		values = append(values, rpc_service_reward_order_serde.AggregateToPBRewardOrder(value))
	}

	return connect.NewResponse(&reward_order_v1.ListRewardOrdersResponse{
		RewardOrders: values,
		TotalCount:   response.GetTotalCount(),
	}), nil
}

func NewRPC(service service_list_reward_orders.Service) RPC {
	return RPC{service: service}
}
