import type { IEvent } from './events.types';

export class EventsService {
  async findUpcoming(_city: string, _signal?: AbortSignal): Promise<IEvent[]> {
    return [];
  }
}
