import type { Transport } from '@connectrpc/connect';
import { BaseRpcApi } from '~/core/app/rpc/base-api';
import { UserService } from '~/gen/user/v1/user_pb';
import type { ApiUser } from './users.api.types';

export class UsersApi extends BaseRpcApi<typeof UserService> {
  constructor(transport: Transport) {
    super(UserService, transport);
  }

  async find(params: ApiUser.FindParams, signal?: AbortSignal): Promise<ApiUser.ListResponse> {
    return this.client.listUsers(params, this.options(signal));
  }

  async get(id: string, signal?: AbortSignal): Promise<ApiUser.Response | undefined> {
    const { user } = await this.client.getUser({ id }, this.options(signal));
    return user;
  }

  async create(data: ApiUser.CreateRequest): Promise<ApiUser.Response | undefined> {
    const { user } = await this.client.createUser(data);
    return user;
  }

  async update(data: ApiUser.UpdateRequest): Promise<ApiUser.Response | undefined> {
    const { user } = await this.client.updateUser(data);
    return user;
  }

  async delete(id: string): Promise<void> {
    await this.client.deleteUser({ id });
  }
}
