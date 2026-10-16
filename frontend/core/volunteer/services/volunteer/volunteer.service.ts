import { ConnectError } from '@connectrpc/connect';
import type { ApiPagination, IListResult } from '~/core/app/types';
import { VOLUNTEER_API } from '../..';
import type { ApiVolunteer } from '../../api/volunteer/volunteer.api.types';
import type {
  ICoordinatorVolunteer,
  IUpdateVolunteerDto,
  IVolunteer,
  IVolunteerOperation,
  VolunteerOperationType,
} from './volunteer.types';

const joinName = (...parts: string[]) => parts.filter(Boolean).join(' ');

export class VolunteerService {
  private service = VOLUNTEER_API().volunteer;

  async getMe(signal?: AbortSignal): Promise<IVolunteer> {
    const volunteer = await this.service.get(signal);
    if (!volunteer) throw new ConnectError('Пустой ответ сервера');
    return this.toInterface(volunteer);
  }

  async getForCoordinator(volunteerId: string, signal?: AbortSignal): Promise<ICoordinatorVolunteer> {
    const volunteer = await this.service.getForCoordinator(volunteerId, signal);
    if (!volunteer) throw new ConnectError('Пустой ответ сервера');
    return this.toCoordinatorView(volunteer);
  }

  async update(dto: IUpdateVolunteerDto): Promise<IVolunteer> {
    await this.service.update(this.toApiUpdateDto(dto));
    return this.getMe();
  }

  async listOperations(
    type: VolunteerOperationType | undefined,
    pagination: ApiPagination,
    signal?: AbortSignal,
  ): Promise<IListResult<IVolunteerOperation>> {
    const { operations, totalCount } = await this.service.listOperations(type, toApiPage(pagination), signal);
    return { list: operations.map((item) => this.toOperation(item)), total: fromBigInt(totalCount) };
  }

  toApiUpdateDto(dto: IUpdateVolunteerDto): ApiVolunteer.UpdateRequest {
    const { firstName, lastName, middleName, dateOfBirth } = dto;
    return { firstName, lastName, middleName, dateOfBirth: dateOfBirth ? toProtoDate(dateOfBirth) : undefined };
  }

  toInterface(data: ApiVolunteer.Response): IVolunteer {
    return {
      id: data.id,
      firstName: data.firstName,
      lastName: data.lastName,
      middleName: data.middleName,
      fullName: joinName(data.lastName, data.firstName, data.middleName) || data.emailAddress,
      dateOfBirth: fromProtoDate(data.dateOfBirth),
      email: data.emailAddress,
      phone: data.phoneNumber,
      coinBalance: fromBigInt(data.coinBalance),
      confirmedMinutes: fromBigInt(data.confirmedMinutes),
    };
  }

  toCoordinatorView(data: ApiVolunteer.CoordinatorView): ICoordinatorVolunteer {
    return {
      id: data.id,
      firstName: data.firstName,
      lastName: data.lastName,
      middleName: data.middleName,
      fullName: joinName(data.lastName, data.firstName, data.middleName),
      email: data.emailAddress ?? null,
      phone: data.phoneNumber ?? null,
    };
  }

  toOperation(data: ApiVolunteer.Operation): IVolunteerOperation {
    return {
      id: data.id,
      coordinatorId: data.coordinatorId,
      volunteerId: data.volunteerId,
      description: data.description,
      type: data.type,
      participationId: data.source.case === 'participationId' ? data.source.value : null,
      rewardOrderId: data.source.case === 'rewardOrderId' ? data.source.value : null,
      coinDelta: fromBigInt(data.coinDelta),
      minuteDelta: fromBigInt(data.minuteDelta),
      coinBalanceAfter: fromBigInt(data.coinBalanceAfter),
      confirmedMinutesAfter: fromBigInt(data.confirmedMinutesAfter),
      sequence: fromBigInt(data.sequence),
      createdAt: fromTimestamp(data.createdAt),
    };
  }
}
