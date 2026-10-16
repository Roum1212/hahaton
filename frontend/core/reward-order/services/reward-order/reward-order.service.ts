import { ConnectError } from '@connectrpc/connect';
import type { ApiPagination, IListResult } from '~/core/app/types';
import { REWARD_ORDER_API } from '../..';
import type { ApiRewardOrder } from '../../api/reward-order/reward-order.api.types';
import type { IRewardOrder, IRewardOrderFilter } from './reward-order.types';

export class RewardOrderService {
  private service = REWARD_ORDER_API().rewardOrder;

  async create(rewardId: string, expectedPriceCoins: number, requestId: string): Promise<IRewardOrder> {
    return this.toInterface(this.required(await this.service.create(rewardId, BigInt(expectedPriceCoins), requestId)));
  }

  async get(id: string, signal?: AbortSignal): Promise<IRewardOrder> {
    return this.toInterface(this.required(await this.service.get(id, signal)));
  }

  async find(filter: IRewardOrderFilter, pagination: ApiPagination, signal?: AbortSignal): Promise<IListResult<IRewardOrder>> {
    const { rewardOrders, totalCount } = await this.service.find({ ...filter, ...toApiPage(pagination) }, signal);
    return { list: rewardOrders.map((item) => this.toInterface(item)), total: fromBigInt(totalCount) };
  }

  async issue(id: string): Promise<IRewardOrder> {
    await this.service.issue(id);
    return this.get(id);
  }

  async cancel(id: string): Promise<IRewardOrder> {
    await this.service.cancel(id);
    return this.get(id);
  }

  toInterface(data: ApiRewardOrder.Response): IRewardOrder {
    return {
      id: data.id,
      coordinatorId: data.coordinatorId,
      rewardId: data.rewardId,
      volunteerId: data.volunteerId,
      priceCoins: fromBigInt(data.priceCoins),
      rewardTitle: data.rewardTitle,
      status: data.status,
      cancellationReason: data.cancellationReason,
      createdAt: fromTimestamp(data.createdAt),
      updatedAt: fromTimestamp(data.updatedAt),
    };
  }

  private required(order: ApiRewardOrder.Response | undefined): ApiRewardOrder.Response {
    if (!order) throw new ConnectError('Пустой ответ сервера');
    return order;
  }
}
