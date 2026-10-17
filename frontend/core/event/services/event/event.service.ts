import { ConnectError } from '@connectrpc/connect';
import type { ApiPagination, IListResult } from '~/core/app/types';
import { GEO_SERVICES } from '~/core/geo';
import { EVENT_API } from '../..';
import type { ApiEvent } from '../../api/event/event.api.types';
import type { IEvent, IEventDto, IEventFilter, IEventResult } from './event.types';

export class EventService {
  private service = EVENT_API().event;
  private geo = GEO_SERVICES().GeoService;
  private baseUrl = String(useRuntimeConfig().public.rpcUrl);

  async find(filter: IEventFilter, pagination: ApiPagination, signal?: AbortSignal): Promise<IListResult<IEvent>> {
    const { events, totalCount } = await this.service.find(this.toApiFindParams(filter, pagination), signal);
    return { list: events.map((item) => this.toInterface(item)), total: fromBigInt(totalCount) };
  }

  async get(id: string, signal?: AbortSignal): Promise<IEvent> {
    return this.toInterface(this.required(await this.service.get(id, signal)));
  }

  async create(dto: IEventDto): Promise<IEvent> {
    const id = await this.service.create(this.toApiCreateDto(dto));
    return this.get(id);
  }

  async update(id: string, dto: Partial<IEventDto>): Promise<IEvent> {
    await this.service.update(this.toApiUpdateDto(id, dto));
    return this.get(id);
  }

  async cancel(id: string): Promise<IEvent> {
    await this.service.cancel(id);
    return this.get(id);
  }

  async complete(id: string, results: IEventResult[], requestId: string): Promise<IEvent> {
    const apiResults: ApiEvent.Result[] = results.map((item) =>
      'absent' in item
        ? { participationId: item.participationId, outcome: { case: 'absent', value: true } }
        : { participationId: item.participationId, outcome: { case: 'confirmedMinutes', value: BigInt(item.confirmedMinutes) } },
    );
    return this.toInterface(this.required(await this.service.complete(id, apiResults, requestId)));
  }

  toApiFindParams(filter: IEventFilter, pagination: ApiPagination): ApiEvent.FindParams {
    return {
      scope: filter.scope,
      searchArea: filter.searchArea ? this.geo.toApiSearchArea(filter.searchArea) : undefined,
      coordinatorId: filter.coordinatorId,
      status: filter.status,
      startsAtFrom: filter.startsFrom ? toTimestamp(filter.startsFrom) : undefined,
      startsAtTo: filter.startsTo ? toTimestamp(filter.startsTo) : undefined,
      endsAtFrom: filter.endsFrom ? toTimestamp(filter.endsFrom) : undefined,
      endsAtTo: filter.endsTo ? toTimestamp(filter.endsTo) : undefined,
      ...toApiPage(pagination),
    };
  }

  toApiCreateDto(dto: IEventDto): ApiEvent.CreateRequest {
    return {
      title: dto.title,
      description: dto.description,
      location: dto.location,
      imagePath: dto.imagePath ?? '',
      point: dto.point ? this.geo.toApiPoint(dto.point) : undefined,
      startsAt: dto.startsAt ? toTimestamp(dto.startsAt) : undefined,
      endsAt: dto.endsAt ? toTimestamp(dto.endsAt) : undefined,
      candidateRequirements: dto.requirements,
      volunteerLimit: dto.volunteerLimit,
    };
  }

  toApiUpdateDto(id: string, dto: Partial<IEventDto>): ApiEvent.UpdateRequest {
    return {
      id,
      title: dto.title,
      description: dto.description,
      location: dto.location,
      imagePath: dto.imagePath,
      point: dto.point ? this.geo.toApiPoint(dto.point) : undefined,
      startsAt: dto.startsAt ? toTimestamp(dto.startsAt) : undefined,
      endsAt: dto.endsAt ? toTimestamp(dto.endsAt) : undefined,
      candidateRequirements: dto.requirements ? { values: dto.requirements } : undefined,
      volunteerLimit: dto.volunteerLimit,
    };
  }

  toInterface(data: ApiEvent.Response): IEvent {
    const startsAt = fromTimestamp(data.startsAt);
    const endsAt = fromTimestamp(data.endsAt);
    return {
      id: data.id,
      coordinatorId: data.coordinatorId,
      title: data.title,
      description: data.description,
      location: data.location,
      point: this.geo.toPoint(data.point),
      imageUrl: resolveImageUrl(data.imageUrl, this.baseUrl),
      startsAt,
      endsAt,
      durationMinutes: startsAt && endsAt ? Math.round((endsAt.getTime() - startsAt.getTime()) / 60000) : 0,
      requirements: data.candidateRequirements,
      volunteerLimit: data.volunteerLimit,
      status: data.status,
      approvedCount: data.approvedParticipantCount,
      pendingCount: data.pendingApplicationCount,
      freePlaces: Math.max(0, data.volunteerLimit - data.approvedParticipantCount),
    };
  }

  private required(event: ApiEvent.Response | undefined): ApiEvent.Response {
    if (!event) throw new ConnectError('Пустой ответ сервера');
    return event;
  }
}
