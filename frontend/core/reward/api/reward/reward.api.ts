import type { Transport } from '@connectrpc/connect';
import { BaseRpcApi } from '~/core/app/rpc/base-api';
import { RewardService } from '#genproto/reward/v1/reward_service_pb';
import type { ApiReward } from './reward.api.types';

export class RewardApi extends BaseRpcApi<typeof RewardService> {
  constructor(transport: Transport) {
    super(RewardService, transport);
  }

  async find(params: ApiReward.FindParams, signal?: AbortSignal): Promise<ApiReward.ListResponse> {
    return this.client.listRewards(params, this.options(signal));
  }

  async get(id: string, signal?: AbortSignal): Promise<ApiReward.Response | undefined> {
    const { reward } = await this.client.getReward({ id }, this.options(signal));
    return reward;
  }

  async create(data: ApiReward.CreateRequest): Promise<string> {
    const { id } = await this.client.createReward(data);
    return id;
  }

  async update(data: ApiReward.UpdateRequest): Promise<void> {
    await this.client.updateReward(data);
  }

  async delete(id: string): Promise<void> {
    await this.client.deleteReward({ id });
  }
}
