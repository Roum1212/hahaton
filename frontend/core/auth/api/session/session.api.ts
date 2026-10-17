import type { Transport } from '@connectrpc/connect';
import { BaseRpcApi } from '~/core/app/rpc/base-api';
import { SessionService } from '#genproto/session/v1/session_service_pb';
import type { ApiSession } from './session.api.types';

export class SessionApi extends BaseRpcApi<typeof SessionService> {
  constructor(transport: Transport) {
    super(SessionService, transport);
  }

  async open(): Promise<void> {
    await this.client.openSession({});
  }

  async authenticate(idToken: string): Promise<void> {
    await this.client.authenticateSession({ idToken });
  }

  async get(): Promise<ApiSession.Response | undefined> {
    const { session } = await this.client.getSession({});
    return session;
  }

  async close(): Promise<void> {
    await this.client.closeSession({});
  }
}
