import { ConnectError } from '@connectrpc/connect';
import type { ApiPagination, IListResult } from '~/core/app/types';
import { COOPERATION_API } from '../..';
import type { ApiCooperation } from '../../api/cooperation/cooperation.api.types';
import type { ICooperation, ICooperationFilter, ICooperationInvitation } from './cooperation.types';

export class CooperationService {
  private service = COOPERATION_API().cooperation;

  async request(coordinatorId: string): Promise<ICooperation> {
    return this.toInterface(this.required(await this.service.request(coordinatorId)));
  }

  async approve(volunteerId: string): Promise<void> {
    await this.service.approve(volunteerId);
  }

  async reject(volunteerId: string): Promise<void> {
    await this.service.reject(volunteerId);
  }

  async terminateAsVolunteer(coordinatorId: string): Promise<void> {
    await this.service.terminate({ case: 'coordinatorId', value: coordinatorId });
  }

  async terminateAsCoordinator(volunteerId: string): Promise<void> {
    await this.service.terminate({ case: 'volunteerId', value: volunteerId });
  }

  async getWithCoordinator(coordinatorId: string, signal?: AbortSignal): Promise<ICooperation> {
    return this.toInterface(this.required(await this.service.get({ case: 'coordinatorId', value: coordinatorId }, signal)));
  }

  async getWithVolunteer(volunteerId: string, signal?: AbortSignal): Promise<ICooperation> {
    return this.toInterface(this.required(await this.service.get({ case: 'volunteerId', value: volunteerId }, signal)));
  }

  async find(filter: ICooperationFilter, pagination: ApiPagination, signal?: AbortSignal): Promise<IListResult<ICooperation>> {
    const { cooperations, totalCount } = await this.service.find({ ...filter, ...toApiPage(pagination) }, signal);
    return { list: cooperations.map((item) => this.toInterface(item)), total: fromBigInt(totalCount) };
  }

  async createInvitation(expiresAt: Date | null): Promise<ICooperationInvitation> {
    const invitation = await this.service.createInvitation(expiresAt ? toTimestamp(expiresAt) : undefined);
    if (!invitation) throw new ConnectError('Пустой ответ сервера');
    return this.toInvitation(invitation);
  }

  async findInvitations(
    includeExpired: boolean,
    pagination: ApiPagination,
    signal?: AbortSignal,
  ): Promise<IListResult<ICooperationInvitation>> {
    const { invitations, totalCount } = await this.service.findInvitations(includeExpired, toApiPage(pagination), signal);
    return { list: invitations.map((item) => this.toInvitation(item)), total: fromBigInt(totalCount) };
  }

  async revokeInvitation(id: string): Promise<void> {
    await this.service.revokeInvitation(id);
  }

  async acceptInvitation(token: string): Promise<ICooperation> {
    return this.toInterface(this.required(await this.service.acceptInvitation(token)));
  }

  toInterface(data: ApiCooperation.Response): ICooperation {
    return {
      coordinatorId: data.coordinatorId,
      volunteerId: data.volunteerId,
      status: data.status,
      retryAfter: fromTimestamp(data.retryAfter),
      createdAt: fromTimestamp(data.createdAt),
      updatedAt: fromTimestamp(data.updatedAt),
    };
  }

  toInvitation(data: ApiCooperation.Invitation): ICooperationInvitation {
    return {
      id: data.id,
      coordinatorId: data.coordinatorId,
      url: data.invitationUrl,
      createdAt: fromTimestamp(data.createdAt),
      expiresAt: fromTimestamp(data.expiresAt),
    };
  }

  private required(cooperation: ApiCooperation.Response | undefined): ApiCooperation.Response {
    if (!cooperation) throw new ConnectError('Пустой ответ сервера');
    return cooperation;
  }
}
