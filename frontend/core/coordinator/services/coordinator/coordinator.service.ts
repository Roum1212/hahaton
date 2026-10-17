import { ConnectError } from '@connectrpc/connect';
import type { ApiPagination, IListResult } from '~/core/app/types';
import { COORDINATOR_API } from '../..';
import type { ApiCoordinator } from '../../api/coordinator/coordinator.api.types';
import type { ICoordinator, ICoordinatorSettings, IUpdateCoordinatorDto } from './coordinator.types';

export class CoordinatorService {
  private service = COORDINATOR_API().coordinator;
  private baseUrl = String(useRuntimeConfig().public.rpcUrl);

  async get(id?: string, signal?: AbortSignal): Promise<ICoordinator> {
    const coordinator = await this.service.get(id, signal);
    if (!coordinator) throw new ConnectError('Пустой ответ сервера');
    return this.toInterface(coordinator);
  }

  async find(pagination: ApiPagination, signal?: AbortSignal): Promise<IListResult<ICoordinator>> {
    const { coordinators, totalCount } = await this.service.find(toApiPage(pagination), signal);
    return { list: coordinators.map((item) => this.toInterface(item)), total: fromBigInt(totalCount) };
  }

  async update(dto: IUpdateCoordinatorDto): Promise<ICoordinator> {
    await this.service.update({ description: dto.description, avatarPath: dto.avatarPath });
    return this.get();
  }

  async getSettings(): Promise<ICoordinatorSettings> {
    const settings = await this.service.getSettings();
    return { cooldownSeconds: fromBigInt(settings?.cooperationRequestCooldownSeconds) };
  }

  async updateSettings(settings: ICoordinatorSettings): Promise<ICoordinatorSettings> {
    await this.service.updateSettings(BigInt(settings.cooldownSeconds));
    return this.getSettings();
  }

  toInterface(data: ApiCoordinator.Response): ICoordinator {
    return {
      id: data.id,
      displayName: data.displayName,
      description: data.description,
      avatarUrl: resolveImageUrl(data.avatarUrl, this.baseUrl),
    };
  }
}
