import {
  ParticipationCancellationReason,
  ParticipationListScope,
  ParticipationStatus,
} from '#genproto/participation/v1/participation_service_pb';

export { ParticipationCancellationReason, ParticipationListScope, ParticipationStatus };

export interface IParticipation {
  id: string;
  eventId: string;
  volunteerId: string;
  status: ParticipationStatus;
  cancellationReason: ParticipationCancellationReason;
  rejectionReason: string | null;
  awardedCoins: number | null;
  awardedMinutes: number | null;
  createdAt: Date | null;
  updatedAt: Date | null;
}

export interface IParticipationFilter {
  scope: ParticipationListScope;
  eventId?: string;
  volunteerId?: string;
  status?: ParticipationStatus;
}

export interface IParticipationChange {
  sequence: number;
  status: ParticipationStatus;
  cancellationReason: ParticipationCancellationReason;
  rejectionReason: string | null;
  actor: 'coordinator' | 'volunteer' | 'system';
  actorId: string | null;
  createdAt: Date | null;
}

export const PARTICIPATION_STATUS_LABELS: Record<ParticipationStatus, string> = {
  [ParticipationStatus.UNSPECIFIED]: 'Неизвестно',
  [ParticipationStatus.PENDING]: 'На рассмотрении',
  [ParticipationStatus.APPROVED]: 'Принята',
  [ParticipationStatus.REJECTED]: 'Отклонена',
  [ParticipationStatus.CONFIRMED]: 'Участие подтверждено',
  [ParticipationStatus.NOT_ATTENDED]: 'Не пришёл',
  [ParticipationStatus.CANCELLED]: 'Отменена',
  [ParticipationStatus.CLOSED]: 'Закрыта без рассмотрения',
};

export const PARTICIPATION_STATUS_SEVERITY: Record<ParticipationStatus, 'success' | 'info' | 'warn' | 'danger' | 'secondary'> = {
  [ParticipationStatus.UNSPECIFIED]: 'secondary',
  [ParticipationStatus.PENDING]: 'warn',
  [ParticipationStatus.APPROVED]: 'info',
  [ParticipationStatus.REJECTED]: 'danger',
  [ParticipationStatus.CONFIRMED]: 'success',
  [ParticipationStatus.NOT_ATTENDED]: 'danger',
  [ParticipationStatus.CANCELLED]: 'secondary',
  [ParticipationStatus.CLOSED]: 'secondary',
};

export const PARTICIPATION_CANCELLATION_LABELS: Record<ParticipationCancellationReason, string> = {
  [ParticipationCancellationReason.UNSPECIFIED]: '',
  [ParticipationCancellationReason.VOLUNTEER_CANCELLED]: 'Отменена волонтёром',
  [ParticipationCancellationReason.COOPERATION_TERMINATED]: 'Сотрудничество прекращено',
  [ParticipationCancellationReason.EVENT_CANCELLED]: 'Мероприятие отменено',
};
