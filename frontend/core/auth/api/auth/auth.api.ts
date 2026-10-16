import { createClient, type Client, type Transport } from '@connectrpc/connect';
import { BaseRpcApi } from '~/core/app/rpc/base-api';
import { AuthService } from '#genproto/auth/v1/auth_pb';
import type { ApiAuth } from './auth.api.types';

export class AuthApi extends BaseRpcApi<typeof AuthService> {
  private readonly publicClient: Client<typeof AuthService>;

  constructor(transport: Transport, publicTransport: Transport) {
    super(AuthService, transport);
    this.publicClient = createClient(AuthService, publicTransport);
  }

  async loginWithFlida(data: ApiAuth.LoginWithFlidaRequest): Promise<string> {
    const { accessToken } = await this.publicClient.loginWithFlida(data);
    return accessToken;
  }

  async logout(): Promise<void> {
    await this.publicClient.logout({});
  }

  async me(): Promise<ApiAuth.UserResponse | undefined> {
    const { user } = await this.client.me({});
    return user;
  }
}
