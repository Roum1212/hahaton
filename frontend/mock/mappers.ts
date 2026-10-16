import { create } from '@bufbuild/protobuf';
import { timestampFromDate } from '@bufbuild/protobuf/wkt';
import { AdminVolunteerSchema } from '../genproto/admin/v1/admin_service_pb';
import { CooperationInvitationSchema, CooperationSchema, CooperationStatus } from '../genproto/cooperation/v1/cooperation_service_pb';
import { CoordinatorSchema } from '../genproto/coordinator/v1/coordinator_service_pb';
import { EventChatMessageSchema } from '../genproto/event-chat/v1/event_chat_service_pb';
import { EventSchema } from '../genproto/event/v1/event_service_pb';
import { CitySchema, GeoPointSchema } from '../genproto/geo/v1/geo_service_pb';
import { DateSchema } from '../genproto/google/type/date_pb';
import { ParticipationChangeSchema, ParticipationSchema, ParticipationStatus } from '../genproto/participation/v1/participation_service_pb';
import { RewardOrderSchema } from '../genproto/reward-order/v1/reward_order_service_pb';
import { RewardSchema, RewardStatus } from '../genproto/reward/v1/reward_service_pb';
import { CoordinatorVolunteerSchema, VolunteerOperationSchema, VolunteerSchema } from '../genproto/volunteer/v1/volunteer_service_pb';
import {
  db,
  fullName,
  type ChatMessageRec,
  type CityRec,
  type CooperationRec,
  type CoordinatorRec,
  type EventRec,
  type InvitationRec,
  type OperationRec,
  type ParticipationChangeRec,
  type ParticipationRec,
  type Point,
  type RewardOrderRec,
  type RewardRec,
  type VolunteerRec,
} from './db';
import { approvedCount, issuedQuantity, pendingCount, reservedQuantity } from './domain';

const PUBLIC_URL = process.env.MOCK_PUBLIC_URL || 'http://localhost:3000';

export const ts = (value: Date) => timestampFromDate(value);

export const toPoint = (point: Point) => create(GeoPointSchema, { latitude: point.latitude, longitude: point.longitude });

export const toVolunteer = (item: VolunteerRec) =>
  create(VolunteerSchema, {
    id: item.id,
    dateOfBirth: create(DateSchema, item.dateOfBirth),
    firstName: item.firstName,
    lastName: item.lastName,
    middleName: item.middleName,
    emailAddress: item.email,
    phoneNumber: item.phone,
    coinBalance: BigInt(item.coinBalance),
    confirmedMinutes: BigInt(item.confirmedMinutes),
  });

export const toCoordinatorVolunteer = (item: VolunteerRec, withContacts: boolean) =>
  create(CoordinatorVolunteerSchema, {
    id: item.id,
    firstName: item.firstName,
    lastName: item.lastName,
    middleName: item.middleName,
    emailAddress: withContacts ? item.email : undefined,
    phoneNumber: withContacts ? item.phone : undefined,
  });

export const toAdminVolunteer = (item: VolunteerRec) =>
  create(AdminVolunteerSchema, {
    id: item.id,
    firstName: item.firstName,
    lastName: item.lastName,
    middleName: item.middleName,
    emailAddress: item.email,
    phoneNumber: item.phone,
    hasCoordinatorProfile: db.coordinators.has(item.id),
  });

export const toCoordinator = (item: CoordinatorRec) => {
  const volunteer = db.volunteers.get(item.id);
  return create(CoordinatorSchema, {
    id: item.id,
    avatarUrl: item.avatarUrl,
    description: item.description,
    displayName: volunteer ? fullName(volunteer) : '',
  });
};

export const toCooperation = (item: CooperationRec) =>
  create(CooperationSchema, {
    coordinatorId: item.coordinatorId,
    volunteerId: item.volunteerId,
    status: item.status,
    retryAfter: item.status === CooperationStatus.REJECTED && item.retryAfter ? ts(item.retryAfter) : undefined,
    createdAt: ts(item.createdAt),
    updatedAt: ts(item.updatedAt),
  });

export const toInvitation = (item: InvitationRec) =>
  create(CooperationInvitationSchema, {
    id: item.id,
    coordinatorId: item.coordinatorId,
    invitationUrl: invitationUrl(item.token),
    createdAt: ts(item.createdAt),
    expiresAt: item.expiresAt ? ts(item.expiresAt) : undefined,
  });

export const invitationUrl = (token: string) => `${PUBLIC_URL}/invite/${token}`;

export const toEvent = (item: EventRec) =>
  create(EventSchema, {
    id: item.id,
    coordinatorId: item.coordinatorId,
    description: item.description,
    title: item.title,
    location: item.location,
    point: toPoint(item.point),
    startsAt: ts(item.startsAt),
    endsAt: ts(item.endsAt),
    candidateRequirements: item.requirements,
    volunteerLimit: item.volunteerLimit,
    status: item.status,
    approvedParticipantCount: approvedCount(item.id),
    pendingApplicationCount: pendingCount(item.id),
  });

export const toParticipation = (item: ParticipationRec) =>
  create(ParticipationSchema, {
    id: item.id,
    eventId: item.eventId,
    volunteerId: item.volunteerId,
    cancellationReason: item.cancellationReason,
    rejectionReason: item.status === ParticipationStatus.REJECTED ? item.rejectionReason : undefined,
    status: item.status,
    accrualOperationId: item.status === ParticipationStatus.CONFIRMED ? item.accrualOperationId : undefined,
    awardedCoins: item.status === ParticipationStatus.CONFIRMED && item.awardedCoins !== undefined ? BigInt(item.awardedCoins) : undefined,
    awardedMinutes:
      item.status === ParticipationStatus.CONFIRMED && item.awardedMinutes !== undefined ? BigInt(item.awardedMinutes) : undefined,
    createdAt: ts(item.createdAt),
    updatedAt: ts(item.updatedAt),
  });

export const toParticipationChange = (item: ParticipationChangeRec) =>
  create(ParticipationChangeSchema, {
    sequence: BigInt(item.sequence),
    cancellationReason: item.cancellationReason,
    rejectionReason: item.rejectionReason,
    status: item.status,
    actor: item.actor.case ? item.actor : { case: undefined },
    createdAt: ts(item.createdAt),
  });

export const toOperation = (item: OperationRec) =>
  create(VolunteerOperationSchema, {
    id: item.id,
    coordinatorId: item.coordinatorId,
    volunteerId: item.volunteerId,
    description: item.description,
    type: item.type,
    source: item.source,
    coinDelta: BigInt(item.coinDelta),
    minuteDelta: BigInt(item.minuteDelta),
    coinBalanceAfter: BigInt(item.coinBalanceAfter),
    confirmedMinutesAfter: BigInt(item.confirmedMinutesAfter),
    sequence: BigInt(item.sequence),
    createdAt: ts(item.createdAt),
  });

export const toReward = (item: RewardRec) =>
  create(RewardSchema, {
    id: item.id,
    coordinatorId: item.coordinatorId,
    description: item.description,
    title: item.title,
    priceCoins: BigInt(item.priceCoins),
    availableQuantity: Math.max(0, item.quantityLimit - issuedQuantity(item.id) - reservedQuantity(item.id)),
    quantityLimit: item.quantityLimit,
    status: RewardStatus.ACTIVE,
  });

export const toRewardOrder = (item: RewardOrderRec) =>
  create(RewardOrderSchema, {
    id: item.id,
    coordinatorId: item.coordinatorId,
    rewardId: item.rewardId,
    volunteerId: item.volunteerId,
    priceCoins: BigInt(item.priceCoins),
    rewardTitle: item.rewardTitle,
    cancellationReason: item.cancellationReason,
    status: item.status,
    createdAt: ts(item.createdAt),
    updatedAt: ts(item.updatedAt),
  });

export const toChatMessage = (item: ChatMessageRec) =>
  create(EventChatMessageSchema, {
    id: item.id,
    eventId: item.eventId,
    author: item.author,
    authorDisplayName: item.authorDisplayName,
    text: item.text,
    sequence: BigInt(item.sequence),
    createdAt: ts(item.createdAt),
  });

export const toCity = (item: CityRec) =>
  create(CitySchema, {
    id: item.id,
    name: item.name,
    regionName: item.regionName,
    center: item.center ? toPoint(item.center) : undefined,
    hasBoundary: !!item.boundary?.length,
  });
