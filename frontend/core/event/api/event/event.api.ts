import type { Transport } from '@connectrpc/connect';
import { BaseRpcApi } from '~/core/app/rpc/base-api';
import { EventService } from '#genproto/event/v1/event_service_pb';
import type { ApiEvent } from './event.api.types';

export class EventApi extends BaseRpcApi<typeof EventService> {
  constructor(transport: Transport) {
    super(EventService, transport);
  }

  async find(params: ApiEvent.FindParams, signal?: AbortSignal): Promise<ApiEvent.ListResponse> {
    return this.client.listEvents(params, this.options(signal));
  }

  async get(id: string, signal?: AbortSignal): Promise<ApiEvent.Response | undefined> {
    const { event } = await this.client.getEvent({ id }, this.options(signal));
    return event;
  }

  async create(data: ApiEvent.CreateRequest): Promise<string> {
    const { id } = await this.client.createEvent(data);
    return id;
  }

  async update(data: ApiEvent.UpdateRequest): Promise<void> {
    await this.client.updateEvent(data);
  }

  async cancel(id: string): Promise<void> {
    await this.client.cancelEvent({ id });
  }

  async complete(id: string, results: ApiEvent.Result[], requestId: string): Promise<ApiEvent.Response | undefined> {
    const { event } = await this.client.completeEvent({ id, results, requestId });
    return event;
  }
}
