import { ConnectError } from '@connectrpc/connect';
import type { ApiPagination, IListResult } from '~/core/app/types';
import { REWARD_API } from '../..';
import type { ApiReward } from '../../api/reward/reward.api.types';
import type { IReward, IRewardDto, IRewardFilter } from './reward.types';

export class RewardService {
  private service = REWARD_API().reward;

  async find(filter: IRewardFilter, pagination: ApiPagination, signal?: AbortSignal): Promise<IListResult<IReward>> {
    const { rewards, totalCount } = await this.service.find({ ...filter, ...toApiPage(pagination) }, signal);
    return { list: rewards.map((item) => this.toInterface(item)), total: fromBigInt(totalCount) };
  }

  async get(id: string, signal?: AbortSignal): Promise<IReward> {
    const reward = await this.service.get(id, signal);
    if (!reward) throw new ConnectError('Пустой ответ сервера');
    return this.toInterface(reward);
  }

  async create(dto: IRewardDto): Promise<IReward> {
    return this.get(await this.service.create(this.toApiCreateDto(dto)));
  }

  async update(id: string, dto: Partial<IRewardDto>): Promise<IReward> {
    await this.service.update(this.toApiUpdateDto(id, dto));
    return this.get(id);
  }

  async delete(id: string): Promise<void> {
    await this.service.delete(id);
  }

  toApiCreateDto(dto: IRewardDto): ApiReward.CreateRequest {
    return {
      title: dto.title,
      description: dto.description,
      priceCoins: BigInt(dto.priceCoins),
      quantityLimit: dto.quantityLimit,
    };
  }

  toApiUpdateDto(id: string, dto: Partial<IRewardDto>): ApiReward.UpdateRequest {
    return {
      id,
      title: dto.title,
      description: dto.description,
      priceCoins: dto.priceCoins === undefined ? undefined : BigInt(dto.priceCoins),
      quantityLimit: dto.quantityLimit,
    };
  }

  toInterface(data: ApiReward.Response): IReward {
    return {
      id: data.id,
      coordinatorId: data.coordinatorId,
      title: data.title,
      description: data.description,
      priceCoins: fromBigInt(data.priceCoins),
      availableQuantity: data.availableQuantity,
      quantityLimit: data.quantityLimit,
    };
  }
}
