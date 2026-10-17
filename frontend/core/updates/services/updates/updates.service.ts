import { UPDATES_API } from '../..';
import type { UpdateKind } from './updates.types';

export class UpdatesService {
  private service = UPDATES_API().updates;

  subscribe(signal?: AbortSignal): AsyncGenerator<UpdateKind> {
    return this.service.subscribe(signal);
  }
}
