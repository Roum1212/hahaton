package rpc_service_reward_order

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	reward_order_v1 "github.com/Roum1212/hahaton/genproto/reward-order/v1"
	reward_order_v1connect "github.com/Roum1212/hahaton/genproto/reward-order/v1/reward_order_v1connect"
	rpc_service_reward_order_rpc_cancel_reward_order "github.com/Roum1212/hahaton/internal/api/rpc/service/reward-order/rpc/cancel-reward-order"
	rpc_service_reward_order_rpc_create_reward_order "github.com/Roum1212/hahaton/internal/api/rpc/service/reward-order/rpc/create-reward-order"
	rpc_service_reward_order_rpc_get_reward_order "github.com/Roum1212/hahaton/internal/api/rpc/service/reward-order/rpc/get-reward-order"
	rpc_service_reward_order_rpc_issue_reward_order "github.com/Roum1212/hahaton/internal/api/rpc/service/reward-order/rpc/issue-reward-order"
	rpc_service_reward_order_rpc_list_reward_orders "github.com/Roum1212/hahaton/internal/api/rpc/service/reward-order/rpc/list-reward-orders"
)

var _ reward_order_v1connect.RewardOrderServiceHandler = Service{}

type Service struct {
	cancelRewardOrderRPC rpc_service_reward_order_rpc_cancel_reward_order.RPC
	createRewardOrderRPC rpc_service_reward_order_rpc_create_reward_order.RPC
	getRewardOrderRPC    rpc_service_reward_order_rpc_get_reward_order.RPC
	issueRewardOrderRPC  rpc_service_reward_order_rpc_issue_reward_order.RPC
	listRewardOrdersRPC  rpc_service_reward_order_rpc_list_reward_orders.RPC
}

func (x Service) CancelRewardOrder(
	ctx context.Context,
	r *connect.Request[reward_order_v1.CancelRewardOrderRequest],
) (*connect.Response[emptypb.Empty], error) {
	return x.cancelRewardOrderRPC.CancelRewardOrder(ctx, r)
}

func (x Service) CreateRewardOrder(
	ctx context.Context,
	r *connect.Request[reward_order_v1.CreateRewardOrderRequest],
) (*connect.Response[reward_order_v1.CreateRewardOrderResponse], error) {
	return x.createRewardOrderRPC.CreateRewardOrder(ctx, r)
}

func (x Service) GetRewardOrder(
	ctx context.Context,
	r *connect.Request[reward_order_v1.GetRewardOrderRequest],
) (*connect.Response[reward_order_v1.GetRewardOrderResponse], error) {
	return x.getRewardOrderRPC.GetRewardOrder(ctx, r)
}

func (x Service) IssueRewardOrder(
	ctx context.Context,
	r *connect.Request[reward_order_v1.IssueRewardOrderRequest],
) (*connect.Response[emptypb.Empty], error) {
	return x.issueRewardOrderRPC.IssueRewardOrder(ctx, r)
}

func (x Service) ListRewardOrders(
	ctx context.Context,
	r *connect.Request[reward_order_v1.ListRewardOrdersRequest],
) (*connect.Response[reward_order_v1.ListRewardOrdersResponse], error) {
	return x.listRewardOrdersRPC.ListRewardOrders(ctx, r)
}

func NewService(
	cancelRewardOrderRPC rpc_service_reward_order_rpc_cancel_reward_order.RPC,
	createRewardOrderRPC rpc_service_reward_order_rpc_create_reward_order.RPC,
	getRewardOrderRPC rpc_service_reward_order_rpc_get_reward_order.RPC,
	issueRewardOrderRPC rpc_service_reward_order_rpc_issue_reward_order.RPC,
	listRewardOrdersRPC rpc_service_reward_order_rpc_list_reward_orders.RPC,
) Service {
	return Service{
		cancelRewardOrderRPC: cancelRewardOrderRPC,
		createRewardOrderRPC: createRewardOrderRPC,
		getRewardOrderRPC:    getRewardOrderRPC,
		issueRewardOrderRPC:  issueRewardOrderRPC,
		listRewardOrdersRPC:  listRewardOrdersRPC,
	}
}
