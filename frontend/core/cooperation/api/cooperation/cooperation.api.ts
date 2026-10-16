import type { Transport } from '@connectrpc/connect';
import type { Timestamp } from '@bufbuild/protobuf/wkt';
import { BaseRpcApi } from '~/core/app/rpc/base-api';
import { CooperationService } from '#genproto/cooperation/v1/cooperation_service_pb';
import type { ApiCooperation } from './cooperation.api.types';

type Counterpart = { case: 'coordinatorId'; value: string } | { case: 'volunteerId'; value: string };

export class CooperationApi extends BaseRpcApi<typeof CooperationService> {
  constructor(transport: Transport) {
    super(CooperationService, transport);
  }

  async request(coordinatorId: string): Promise<ApiCooperation.Response | undefined> {
    const { cooperation } = await this.client.requestCooperation({ coordinatorId });
    return cooperation;
  }

  async approve(volunteerId: string): Promise<void> {
    await this.client.approveCooperation({ volunteerId });
  }

  async reject(volunteerId: string): Promise<void> {
    await this.client.rejectCooperation({ volunteerId });
  }

  async terminate(counterpart: Counterpart): Promise<void> {
    await this.client.terminateCooperation({ counterpart });
  }

  async get(counterpart: Counterpart, signal?: AbortSignal): Promise<ApiCooperation.Response | undefined> {
    const { cooperation } = await this.client.getCooperation({ counterpart }, this.options(signal));
    return cooperation;
  }

  async find(params: ApiCooperation.FindParams, signal?: AbortSignal): Promise<ApiCooperation.ListResponse> {
    return this.client.listCooperations(params, this.options(signal));
  }

  async createInvitation(expiresAt?: Timestamp): Promise<ApiCooperation.Invitation | undefined> {
    const { invitation } = await this.client.createCooperationInvitation({ expiresAt });
    return invitation;
  }

  async findInvitations(
    includeExpired: boolean,
    page: { limit: number; offset: bigint },
    signal?: AbortSignal,
  ): Promise<ApiCooperation.InvitationsResponse> {
    return this.client.listCooperationInvitations({ includeExpired, ...page }, this.options(signal));
  }

  async revokeInvitation(id: string): Promise<void> {
    await this.client.revokeCooperationInvitation({ id });
  }

  async acceptInvitation(token: string): Promise<ApiCooperation.Response | undefined> {
    const { cooperation } = await this.client.acceptCooperationInvitation({ token });
    return cooperation;
  }
}
