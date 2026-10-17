package rpc_service_reward

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	reward_v1 "github.com/Roum1212/hahaton/genproto/reward/v1"
	reward_v1connect "github.com/Roum1212/hahaton/genproto/reward/v1/reward_v1connect"
	rpc_service_reward_rpc_create_reward "github.com/Roum1212/hahaton/internal/api/rpc/service/reward/rpc/create-reward"
	rpc_service_reward_rpc_delete_reward "github.com/Roum1212/hahaton/internal/api/rpc/service/reward/rpc/delete-reward"
	rpc_service_reward_rpc_get_reward "github.com/Roum1212/hahaton/internal/api/rpc/service/reward/rpc/get-reward"
	rpc_service_reward_rpc_list_rewards "github.com/Roum1212/hahaton/internal/api/rpc/service/reward/rpc/list-rewards"
	rpc_service_reward_rpc_update_reward "github.com/Roum1212/hahaton/internal/api/rpc/service/reward/rpc/update-reward"
)

var _ reward_v1connect.RewardServiceHandler = Service{}

type Service struct {
	createRewardRPC rpc_service_reward_rpc_create_reward.RPC
	deleteRewardRPC rpc_service_reward_rpc_delete_reward.RPC
	getRewardRPC    rpc_service_reward_rpc_get_reward.RPC
	listRewardsRPC  rpc_service_reward_rpc_list_rewards.RPC
	updateRewardRPC rpc_service_reward_rpc_update_reward.RPC
}

func (x Service) CreateReward(
	ctx context.Context,
	r *connect.Request[reward_v1.CreateRewardRequest],
) (*connect.Response[reward_v1.CreateRewardResponse], error) {
	return x.createRewardRPC.CreateReward(ctx, r)
}

func (x Service) DeleteReward(
	ctx context.Context,
	r *connect.Request[reward_v1.DeleteRewardRequest],
) (*connect.Response[emptypb.Empty], error) {
	return x.deleteRewardRPC.DeleteReward(ctx, r)
}

func (x Service) GetReward(
	ctx context.Context,
	r *connect.Request[reward_v1.GetRewardRequest],
) (*connect.Response[reward_v1.GetRewardResponse], error) {
	return x.getRewardRPC.GetReward(ctx, r)
}

func (x Service) ListRewards(
	ctx context.Context,
	r *connect.Request[reward_v1.ListRewardsRequest],
) (*connect.Response[reward_v1.ListRewardsResponse], error) {
	return x.listRewardsRPC.ListRewards(ctx, r)
}

func (x Service) UpdateReward(
	ctx context.Context,
	r *connect.Request[reward_v1.UpdateRewardRequest],
) (*connect.Response[emptypb.Empty], error) {
	return x.updateRewardRPC.UpdateReward(ctx, r)
}

func NewService(
	createRewardRPC rpc_service_reward_rpc_create_reward.RPC,
	deleteRewardRPC rpc_service_reward_rpc_delete_reward.RPC,
	getRewardRPC rpc_service_reward_rpc_get_reward.RPC,
	listRewardsRPC rpc_service_reward_rpc_list_rewards.RPC,
	updateRewardRPC rpc_service_reward_rpc_update_reward.RPC,
) Service {
	return Service{
		createRewardRPC: createRewardRPC,
		deleteRewardRPC: deleteRewardRPC,
		getRewardRPC:    getRewardRPC,
		listRewardsRPC:  listRewardsRPC,
		updateRewardRPC: updateRewardRPC,
	}
}
