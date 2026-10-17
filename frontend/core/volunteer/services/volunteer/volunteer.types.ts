import { VolunteerOperationType } from '#genproto/volunteer/v1/volunteer_service_pb';

export { VolunteerOperationType };

export interface IVolunteer {
  id: string;
  firstName: string;
  lastName: string;
  middleName: string;
  fullName: string;
  avatarUrl: string | null;
  dateOfBirth: Date | null;
  email: string;
  phone: string;
  coinBalance: number;
  confirmedMinutes: number;
}

export interface ICoordinatorVolunteer {
  id: string;
  firstName: string;
  lastName: string;
  middleName: string;
  fullName: string;
  avatarUrl: string | null;
  email: string | null;
  phone: string | null;
}

export interface IUpdateVolunteerDto {
  firstName: string;
  lastName: string;
  middleName: string;
  dateOfBirth: Date | null;
}

export interface IVolunteerOperation {
  id: string;
  coordinatorId: string;
  volunteerId: string;
  description: string;
  type: VolunteerOperationType;
  participationId: string | null;
  rewardOrderId: string | null;
  coinDelta: number;
  minuteDelta: number;
  coinBalanceAfter: number;
  confirmedMinutesAfter: number;
  sequence: number;
  createdAt: Date | null;
}

export const VOLUNTEER_OPERATION_TYPE_LABELS: Record<VolunteerOperationType, string> = {
  [VolunteerOperationType.UNSPECIFIED]: 'Операция',
  [VolunteerOperationType.PARTICIPATION_ACCRUAL]: 'Начисление за участие',
  [VolunteerOperationType.REWARD_PURCHASE]: 'Покупка поощрения',
  [VolunteerOperationType.REWARD_REFUND]: 'Возврат за поощрение',
};
