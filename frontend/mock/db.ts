import { randomBytes, randomUUID } from 'node:crypto';
import { Code, ConnectError } from '@connectrpc/connect';
import type { CooperationStatus } from '../genproto/cooperation/v1/cooperation_service_pb';
import type { EventCandidateRequirement, EventStatus } from '../genproto/event/v1/event_service_pb';
import type { ParticipationCancellationReason, ParticipationStatus } from '../genproto/participation/v1/participation_service_pb';
import type { RewardOrderCancellationReason, RewardOrderStatus } from '../genproto/reward-order/v1/reward_order_service_pb';
import type { VolunteerOperationType } from '../genproto/volunteer/v1/volunteer_service_pb';

export interface DateParts {
  year: number;
  month: number;
  day: number;
}

export interface Point {
  latitude: number;
  longitude: number;
}

export interface VolunteerRec {
  id: string;
  iss: string;
  sub: string;
  firstName: string;
  lastName: string;
  middleName: string;
  email: string;
  phone: string;
  dateOfBirth: DateParts;
  avatarPath: string;
  coinBalance: number;
  confirmedMinutes: number;
  operationSequence: number;
}

export interface CoordinatorRec {
  id: string;
  avatarPath: string;
  description: string;
  cooldownSeconds: number;
}

export interface SessionRec {
  id: string;
  userId?: string;
  expiresAt: Date;
  revoked: boolean;
}

export interface CooperationRec {
  coordinatorId: string;
  volunteerId: string;
  status: CooperationStatus;
  retryAfter?: Date;
  createdAt: Date;
  updatedAt: Date;
}

export interface InvitationRec {
  id: string;
  coordinatorId: string;
  token: string;
  createdAt: Date;
  expiresAt?: Date;
}

export interface CompletionRec {
  requestId: string;
  payloadKey: string;
}

export interface EventRec {
  id: string;
  coordinatorId: string;
  title: string;
  description: string;
  location: string;
  point: Point;
  startsAt: Date;
  endsAt: Date;
  requirements: EventCandidateRequirement[];
  volunteerLimit: number;
  imagePath: string;
  status: EventStatus;
  completion?: CompletionRec;
}

export type Actor =
  { case: 'coordinatorId'; value: string } | { case: 'volunteerId'; value: string } | { case: undefined; value?: undefined };

export interface ParticipationRec {
  id: string;
  eventId: string;
  volunteerId: string;
  status: ParticipationStatus;
  cancellationReason: ParticipationCancellationReason;
  rejectionReason?: string;
  accrualOperationId?: string;
  awardedCoins?: number;
  awardedMinutes?: number;
  createdAt: Date;
  updatedAt: Date;
}

export interface ParticipationChangeRec {
  participationId: string;
  sequence: number;
  status: ParticipationStatus;
  cancellationReason: ParticipationCancellationReason;
  rejectionReason?: string;
  actor: Actor;
  createdAt: Date;
}

export interface OperationRec {
  id: string;
  coordinatorId: string;
  volunteerId: string;
  description: string;
  type: VolunteerOperationType;
  source: { case: 'participationId'; value: string } | { case: 'rewardOrderId'; value: string };
  coinDelta: number;
  minuteDelta: number;
  coinBalanceAfter: number;
  confirmedMinutesAfter: number;
  sequence: number;
  createdAt: Date;
}

export interface RewardRec {
  id: string;
  coordinatorId: string;
  title: string;
  description: string;
  priceCoins: number;
  quantityLimit: number;
  imagePath: string;
  deleted: boolean;
}

export interface ImageRec {
  path: string;
  ownerId: string;
  contentType: string;
  data: Buffer;
}

export interface RewardOrderRec {
  id: string;
  coordinatorId: string;
  rewardId: string;
  volunteerId: string;
  priceCoins: number;
  rewardTitle: string;
  status: RewardOrderStatus;
  cancellationReason: RewardOrderCancellationReason;
  requestId: string;
  createdAt: Date;
  updatedAt: Date;
}

export interface ChatAccessRec {
  visibleAfterSequence: number;
  open: boolean;
}

export interface ChatMessageRec {
  id: string;
  eventId: string;
  author: { case: 'coordinatorId'; value: string } | { case: 'volunteerId'; value: string };
  authorDisplayName: string;
  text: string;
  sequence: number;
  requestKey: string;
  createdAt: Date;
}

export interface CityRec {
  id: string;
  name: string;
  regionName: string;
  center?: Point;
  boundary?: Point[];
}

export const db = {
  sessions: new Map<string, SessionRec>(),
  volunteers: new Map<string, VolunteerRec>(),
  coordinators: new Map<string, CoordinatorRec>(),
  cooperations: new Map<string, CooperationRec>(),
  invitations: new Map<string, InvitationRec>(),
  events: new Map<string, EventRec>(),
  participations: new Map<string, ParticipationRec>(),
  participationChanges: [] as ParticipationChangeRec[],
  operations: [] as OperationRec[],
  rewards: new Map<string, RewardRec>(),
  rewardOrders: new Map<string, RewardOrderRec>(),
  chatAccess: new Map<string, ChatAccessRec>(),
  chatMessages: [] as ChatMessageRec[],
  chatSequences: new Map<string, number>(),
  cities: new Map<string, CityRec>(),
  completionResults: new Map<string, string>(),
  images: new Map<string, ImageRec>(),
};

export const newId = () => randomUUID();
export const newToken = () => randomBytes(24).toString('base64url');
export const now = () => new Date();

export const cooperationKey = (coordinatorId: string, volunteerId: string) => `${coordinatorId}:${volunteerId}`;
export const chatAccessKey = (eventId: string, volunteerId: string) => `${eventId}:${volunteerId}`;

export const fail = (code: Code, message: string): never => {
  throw new ConnectError(message, code);
};

export const requireFound = <T>(value: T | undefined, message: string): T => {
  if (value === undefined) fail(Code.NotFound, message);
  return value as T;
};

export const paginate = <T>(items: T[], limit: number, offset: bigint | number, defaultLimit = 20, maxLimit = 100) => {
  const size = limit === 0 ? defaultLimit : limit;
  if (size < 1 || size > maxLimit) fail(Code.InvalidArgument, `limit должен быть от 1 до ${maxLimit}`);
  const start = Number(offset);
  if (start < 0) fail(Code.InvalidArgument, 'offset должен быть неотрицательным');
  return { page: items.slice(start, start + size), totalCount: BigInt(items.length) };
};

export const compareStrings = (a: string, b: string) => (a < b ? -1 : a > b ? 1 : 0);

export const fullName = (volunteer: VolunteerRec) =>
  [volunteer.lastName, volunteer.firstName, volunteer.middleName].filter(Boolean).join(' ');

export const coinsForMinutes = (minutes: number) => (minutes / 30) * 5;
