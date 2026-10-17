import type { Transport } from '@connectrpc/connect';
import { BaseRpcApi } from '~/core/app/rpc/base-api';
import { UpdatesService, type UpdateKind } from '#genproto/updates/v1/updates_service_pb';

export class UpdatesApi extends BaseRpcApi<typeof UpdatesService> {
  constructor(transport: Transport) {
    super(UpdatesService, transport);
  }

  async *subscribe(signal?: AbortSignal): AsyncGenerator<UpdateKind> {
    for await (const item of this.client.subscribeUpdates({}, this.options(signal))) yield item.kind;
  }
}
