import type { Transport } from '@connectrpc/connect';
import { BaseRpcApi } from '~/core/app/rpc/base-api';
import { RewardOrderService } from '#genproto/reward-order/v1/reward_order_service_pb';
import type { ApiRewardOrder } from './reward-order.api.types';

export class RewardOrderApi extends BaseRpcApi<typeof RewardOrderService> {
  constructor(transport: Transport) {
    super(RewardOrderService, transport);
  }

  async create(rewardId: string, expectedPriceCoins: bigint, requestId: string): Promise<ApiRewardOrder.Response | undefined> {
    const { rewardOrder } = await this.client.createRewardOrder({ rewardId, expectedPriceCoins, requestId });
    return rewardOrder;
  }

  async get(id: string, signal?: AbortSignal): Promise<ApiRewardOrder.Response | undefined> {
    const { rewardOrder } = await this.client.getRewardOrder({ id }, this.options(signal));
    return rewardOrder;
  }

  async find(params: ApiRewardOrder.FindParams, signal?: AbortSignal): Promise<ApiRewardOrder.ListResponse> {
    return this.client.listRewardOrders(params, this.options(signal));
  }

  async issue(id: string): Promise<void> {
    await this.client.issueRewardOrder({ id });
  }

  async cancel(id: string): Promise<void> {
    await this.client.cancelRewardOrder({ id });
  }
}
