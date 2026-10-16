import { ConnectError } from '@connectrpc/connect';
import type { ApiPagination, IListResult } from '~/core/app/types';
import { PARTICIPATION_API } from '../..';
import type { ApiParticipation } from '../../api/participation/participation.api.types';
import type { IParticipation, IParticipationChange, IParticipationFilter } from './participation.types';

export class ParticipationService {
  private service = PARTICIPATION_API().participation;

  async request(eventId: string): Promise<IParticipation> {
    return this.get(await this.service.request(eventId));
  }

  async restore(id: string): Promise<IParticipation> {
    await this.service.restore(id);
    return this.get(id);
  }

  async approve(id: string): Promise<IParticipation> {
    await this.service.approve(id);
    return this.get(id);
  }

  async reject(id: string, reason?: string): Promise<IParticipation> {
    await this.service.reject(id, reason?.trim() || undefined);
    return this.get(id);
  }

  async cancel(id: string): Promise<IParticipation> {
    await this.service.cancel(id);
    return this.get(id);
  }

  async get(id: string, signal?: AbortSignal): Promise<IParticipation> {
    const participation = await this.service.get(id, signal);
    if (!participation) throw new ConnectError('Пустой ответ сервера');
    return this.toInterface(participation);
  }

  async find(filter: IParticipationFilter, pagination: ApiPagination, signal?: AbortSignal): Promise<IListResult<IParticipation>> {
    const { participations, totalCount } = await this.service.find({ ...filter, ...toApiPage(pagination) }, signal);
    return { list: participations.map((item) => this.toInterface(item)), total: fromBigInt(totalCount) };
  }

  async findChanges(participationId: string, pagination: ApiPagination, signal?: AbortSignal): Promise<IListResult<IParticipationChange>> {
    const { changes, totalCount } = await this.service.findChanges(participationId, toApiPage(pagination), signal);
    return { list: changes.map((item) => this.toChange(item)), total: fromBigInt(totalCount) };
  }

  toInterface(data: ApiParticipation.Response): IParticipation {
    return {
      id: data.id,
      eventId: data.eventId,
      volunteerId: data.volunteerId,
      status: data.status,
      cancellationReason: data.cancellationReason,
      rejectionReason: data.rejectionReason ?? null,
      awardedCoins: data.awardedCoins === undefined ? null : fromBigInt(data.awardedCoins),
      awardedMinutes: data.awardedMinutes === undefined ? null : fromBigInt(data.awardedMinutes),
      createdAt: fromTimestamp(data.createdAt),
      updatedAt: fromTimestamp(data.updatedAt),
    };
  }

  toChange(data: ApiParticipation.Change): IParticipationChange {
    return {
      sequence: fromBigInt(data.sequence),
      status: data.status,
      cancellationReason: data.cancellationReason,
      rejectionReason: data.rejectionReason ?? null,
      actor: data.actor.case === 'coordinatorId' ? 'coordinator' : data.actor.case === 'volunteerId' ? 'volunteer' : 'system',
      actorId: data.actor.value ?? null,
      createdAt: fromTimestamp(data.createdAt),
    };
  }
}
