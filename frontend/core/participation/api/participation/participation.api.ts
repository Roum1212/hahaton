import type { Transport } from '@connectrpc/connect';
import { BaseRpcApi } from '~/core/app/rpc/base-api';
import { ParticipationService } from '#genproto/participation/v1/participation_service_pb';
import type { ApiParticipation } from './participation.api.types';

export class ParticipationApi extends BaseRpcApi<typeof ParticipationService> {
  constructor(transport: Transport) {
    super(ParticipationService, transport);
  }

  async request(eventId: string): Promise<string> {
    const { id } = await this.client.requestParticipation({ eventId });
    return id;
  }

  async restore(id: string): Promise<void> {
    await this.client.restoreParticipation({ id });
  }

  async approve(id: string): Promise<void> {
    await this.client.approveParticipation({ id });
  }

  async reject(id: string, reason?: string): Promise<void> {
    await this.client.rejectParticipation({ id, reason });
  }

  async cancel(id: string): Promise<void> {
    await this.client.cancelParticipation({ id });
  }

  async get(id: string, signal?: AbortSignal): Promise<ApiParticipation.Response | undefined> {
    const { participation } = await this.client.getParticipation({ id }, this.options(signal));
    return participation;
  }

  async find(params: ApiParticipation.FindParams, signal?: AbortSignal): Promise<ApiParticipation.ListResponse> {
    return this.client.listParticipations(params, this.options(signal));
  }

  async findChanges(
    participationId: string,
    page: { limit: number; offset: bigint },
    signal?: AbortSignal,
  ): Promise<ApiParticipation.ChangesResponse> {
    return this.client.listParticipationChanges({ participationId, ...page }, this.options(signal));
  }
}
