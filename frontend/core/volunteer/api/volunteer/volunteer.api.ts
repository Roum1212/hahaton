import type { Transport } from '@connectrpc/connect';
import { BaseRpcApi } from '~/core/app/rpc/base-api';
import { VolunteerService } from '#genproto/volunteer/v1/volunteer_service_pb';
import type { ApiVolunteer } from './volunteer.api.types';

export class VolunteerApi extends BaseRpcApi<typeof VolunteerService> {
  constructor(transport: Transport) {
    super(VolunteerService, transport);
  }

  async get(signal?: AbortSignal): Promise<ApiVolunteer.Response | undefined> {
    const { volunteer } = await this.client.getVolunteer({}, this.options(signal));
    return volunteer;
  }

  async getForCoordinator(volunteerId: string, signal?: AbortSignal): Promise<ApiVolunteer.CoordinatorView | undefined> {
    const { volunteer } = await this.client.getCoordinatorVolunteer({ volunteerId }, this.options(signal));
    return volunteer;
  }

  async update(data: ApiVolunteer.UpdateRequest): Promise<void> {
    await this.client.updateVolunteer(data);
  }

  async listOperations(
    type: ApiVolunteer.OperationType | undefined,
    page: { limit: number; offset: bigint },
    signal?: AbortSignal,
  ): Promise<ApiVolunteer.OperationsResponse> {
    return this.client.listVolunteerOperations({ type, ...page }, this.options(signal));
  }
}
