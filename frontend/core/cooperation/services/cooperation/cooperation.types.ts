import { CooperationListScope, CooperationStatus } from '#genproto/cooperation/v1/cooperation_service_pb';

export { CooperationListScope, CooperationStatus };

export interface ICooperation {
  coordinatorId: string;
  volunteerId: string;
  status: CooperationStatus;
  retryAfter: Date | null;
  createdAt: Date | null;
  updatedAt: Date | null;
}

export interface ICooperationFilter {
  scope: CooperationListScope;
  status?: CooperationStatus;
}

export interface ICooperationInvitation {
  id: string;
  coordinatorId: string;
  url: string;
  createdAt: Date | null;
  expiresAt: Date | null;
}

export const COOPERATION_STATUS_LABELS: Record<CooperationStatus, string> = {
  [CooperationStatus.UNSPECIFIED]: 'Неизвестно',
  [CooperationStatus.PENDING]: 'Заявка на рассмотрении',
  [CooperationStatus.ACTIVE]: 'Сотрудничество',
  [CooperationStatus.REJECTED]: 'Заявка отклонена',
  [CooperationStatus.TERMINATED]: 'Сотрудничество прекращено',
};

export const COOPERATION_STATUS_SEVERITY: Record<CooperationStatus, 'success' | 'info' | 'warn' | 'danger' | 'secondary'> = {
  [CooperationStatus.UNSPECIFIED]: 'secondary',
  [CooperationStatus.PENDING]: 'warn',
  [CooperationStatus.ACTIVE]: 'success',
  [CooperationStatus.REJECTED]: 'danger',
  [CooperationStatus.TERMINATED]: 'secondary',
};
