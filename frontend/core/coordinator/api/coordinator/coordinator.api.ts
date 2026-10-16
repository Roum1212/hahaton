import type { Transport } from '@connectrpc/connect';
import { BaseRpcApi } from '~/core/app/rpc/base-api';
import { CoordinatorService } from '#genproto/coordinator/v1/coordinator_service_pb';
import type { ApiCoordinator } from './coordinator.api.types';

export class CoordinatorApi extends BaseRpcApi<typeof CoordinatorService> {
  constructor(transport: Transport) {
    super(CoordinatorService, transport);
  }

  async get(id?: string, signal?: AbortSignal): Promise<ApiCoordinator.Response | undefined> {
    const { coordinator } = await this.client.getCoordinator({ id }, this.options(signal));
    return coordinator;
  }

  async find(page: { limit: number; offset: bigint }, signal?: AbortSignal): Promise<ApiCoordinator.ListResponse> {
    return this.client.listCoordinators(page, this.options(signal));
  }

  async update(data: ApiCoordinator.UpdateRequest): Promise<void> {
    await this.client.updateCoordinator(data);
  }

  async getSettings(): Promise<ApiCoordinator.Settings | undefined> {
    const { settings } = await this.client.getCoordinatorSettings({});
    return settings;
  }

  async updateSettings(cooldownSeconds: bigint): Promise<void> {
    await this.client.updateCoordinatorSettings({ cooperationRequestCooldownSeconds: cooldownSeconds });
  }
}
