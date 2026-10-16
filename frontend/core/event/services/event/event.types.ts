import { EventCandidateRequirement, EventListScope, EventStatus } from '#genproto/event/v1/event_service_pb';
import type { IGeoPoint, ISearchArea } from '~/core/geo/services/geo/geo.types';

export { EventCandidateRequirement, EventListScope, EventStatus };

export interface IEvent {
  id: string;
  coordinatorId: string;
  title: string;
  description: string;
  location: string;
  point: IGeoPoint | null;
  startsAt: Date | null;
  endsAt: Date | null;
  durationMinutes: number;
  requirements: EventCandidateRequirement[];
  volunteerLimit: number;
  status: EventStatus;
  approvedCount: number;
  pendingCount: number;
  freePlaces: number;
}

export interface IEventDto {
  title: string;
  description: string;
  location: string;
  point: IGeoPoint | null;
  startsAt: Date | null;
  endsAt: Date | null;
  requirements: EventCandidateRequirement[];
  volunteerLimit: number;
}

export interface IEventFilter {
  scope: EventListScope;
  searchArea?: ISearchArea;
  coordinatorId?: string;
  status?: EventStatus;
  startsFrom?: Date;
  startsTo?: Date;
  endsFrom?: Date;
  endsTo?: Date;
}

export type IEventResult = { participationId: string; absent: true } | { participationId: string; confirmedMinutes: number };

export const EVENT_DURATION_STEP_MINUTES = 30;

export const EVENT_STATUS_LABELS: Record<EventStatus, string> = {
  [EventStatus.UNSPECIFIED]: 'Неизвестно',
  [EventStatus.PUBLISHED]: 'Опубликовано',
  [EventStatus.CANCELLED]: 'Отменено',
  [EventStatus.COMPLETED]: 'Завершено',
};

export const EVENT_STATUS_SEVERITY: Record<EventStatus, 'success' | 'info' | 'warn' | 'danger' | 'secondary'> = {
  [EventStatus.UNSPECIFIED]: 'secondary',
  [EventStatus.PUBLISHED]: 'info',
  [EventStatus.CANCELLED]: 'danger',
  [EventStatus.COMPLETED]: 'success',
};

export const EVENT_REQUIREMENT_LABELS: Record<EventCandidateRequirement, string> = {
  [EventCandidateRequirement.UNSPECIFIED]: 'Не указано',
  [EventCandidateRequirement.AGE_14_PLUS]: '14+',
  [EventCandidateRequirement.AGE_16_PLUS]: '16+',
  [EventCandidateRequirement.AGE_18_PLUS]: '18+',
  [EventCandidateRequirement.PRELIMINARY_TRAINING]: 'Предварительное обучение',
  [EventCandidateRequirement.RUSSIAN_CITIZENSHIP]: 'Гражданство РФ',
  [EventCandidateRequirement.VOLUNTEERING_EXPERIENCE]: 'Опыт волонтёрства',
};

export const EVENT_AGE_REQUIREMENTS = [
  EventCandidateRequirement.AGE_14_PLUS,
  EventCandidateRequirement.AGE_16_PLUS,
  EventCandidateRequirement.AGE_18_PLUS,
];
