import { ConnectError } from '@connectrpc/connect';
import { timestampDate } from '@bufbuild/protobuf/wkt';
import type { ApiPagination, MapperService } from '~/core/app/types';
import { USERS_API } from '../..';
import type { ApiUser } from '../../api/users/users.api.types';
import type { ICreateUserDto, IUser } from './users.types';

export class UsersService implements MapperService<
  ICreateUserDto,
  ICreateUserDto & { id: string },
  ApiUser.CreateRequest,
  ApiUser.UpdateRequest,
  ApiUser.Response,
  IUser
> {
  private service = USERS_API().users;

  async find(
    search: string | undefined = undefined,
    pagination: ApiPagination = { limit: 50, offset: 0 },
    signal?: AbortSignal,
  ): Promise<{ list: IUser[]; total: number }> {
    const { users, total } = await this.service.find({ search, page: pagination }, signal);
    return { list: users.map((u) => this.toInterface(u)), total };
  }

  async get(id: string): Promise<IUser> {
    return this.toInterface(this.required(await this.service.get(id)));
  }

  async create(dto: ICreateUserDto): Promise<IUser> {
    return this.toInterface(this.required(await this.service.create(this.toApiCreateDto(dto))));
  }

  async update(id: string, dto: ICreateUserDto): Promise<IUser> {
    return this.toInterface(this.required(await this.service.update(this.toApiUpdateDto({ ...dto, id }))));
  }

  async delete(id: string): Promise<boolean> {
    await this.service.delete(id);
    return true;
  }

  toApiCreateDto(dto: ICreateUserDto): ApiUser.CreateRequest {
    const { email, firstName, lastName, role } = dto;
    return { email, firstName, lastName, role };
  }

  toApiUpdateDto(dto: ICreateUserDto & { id: string }): ApiUser.UpdateRequest {
    const { id, email, firstName, lastName, role, active } = dto;
    return { id, email, firstName, lastName, role, active: !!active };
  }

  toInterface(data: ApiUser.Response): IUser {
    return {
      id: data.id,
      email: data.email,
      firstName: data.firstName,
      lastName: data.lastName,
      fullName: [data.firstName, data.lastName].filter(Boolean).join(' ') || data.email,
      role: data.role,
      active: data.active,
      createdAt: data.createdAt ? timestampDate(data.createdAt) : null,
    };
  }

  private required(user: ApiUser.Response | undefined): ApiUser.Response {
    if (!user) throw new ConnectError('Пустой ответ сервера');
    return user;
  }
}
