import type { DescService } from '@bufbuild/protobuf';
import { createClient, type CallOptions, type Client, type Transport } from '@connectrpc/connect';

export abstract class BaseRpcApi<S extends DescService> {
  protected readonly client: Client<S>;

  constructor(service: S, transport: Transport) {
    this.client = createClient(service, transport);
  }

  protected options(signal?: AbortSignal): CallOptions {
    return { signal };
  }
}
