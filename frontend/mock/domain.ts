import { Code } from '@connectrpc/connect';
import { CooperationStatus } from '../genproto/cooperation/v1/cooperation_service_pb';
import { EventStatus } from '../genproto/event/v1/event_service_pb';
import { ParticipationCancellationReason, ParticipationStatus } from '../genproto/participation/v1/participation_service_pb';
import { RewardOrderCancellationReason, RewardOrderStatus } from '../genproto/reward-order/v1/reward_order_service_pb';
import { VolunteerOperationType } from '../genproto/volunteer/v1/volunteer_service_pb';
import { touch } from './auth';
import { chatHub } from './chat-hub';
import {
  chatAccessKey,
  db,
  fail,
  newId,
  now,
  type Actor,
  type CooperationRec,
  type EventRec,
  type OperationRec,
  type ParticipationRec,
  type RewardOrderRec,
} from './db';

export const SYSTEM: Actor = { case: undefined };

interface OperationInput {
  volunteerId: string;
  coordinatorId: string;
  type: VolunteerOperationType;
  source: OperationRec['source'];
  description: string;
  coinDelta: number;
  minuteDelta: number;
  createdAt?: Date;
}

export const recordOperation = (input: OperationInput): OperationRec => {
  const volunteer = db.volunteers.get(input.volunteerId);
  if (!volunteer) return fail(Code.NotFound, 'Волонтёр не найден');
  const duplicate = db.operations.find(
    (item) => item.type === input.type && item.source.case === input.source.case && item.source.value === input.source.value,
  );
  if (duplicate) return fail(Code.AlreadyExists, 'Операция по этому основанию уже записана');
  const coinBalance = volunteer.coinBalance + input.coinDelta;
  const confirmedMinutes = volunteer.confirmedMinutes + input.minuteDelta;
  if (coinBalance < 0) return fail(Code.FailedPrecondition, 'Недостаточно коинов');
  if (!Number.isSafeInteger(coinBalance) || !Number.isSafeInteger(confirmedMinutes))
    return fail(Code.InvalidArgument, 'Переполнение итогов');
  volunteer.coinBalance = coinBalance;
  volunteer.confirmedMinutes = confirmedMinutes;
  volunteer.operationSequence += 1;
  const operation: OperationRec = {
    id: newId(),
    coordinatorId: input.coordinatorId,
    volunteerId: input.volunteerId,
    description: input.description,
    type: input.type,
    source: input.source,
    coinDelta: input.coinDelta,
    minuteDelta: input.minuteDelta,
    coinBalanceAfter: coinBalance,
    confirmedMinutesAfter: confirmedMinutes,
    sequence: volunteer.operationSequence,
    createdAt: input.createdAt ?? now(),
  };
  db.operations.push(operation);
  return operation;
};

interface ChangeOptions {
  actor: Actor;
  cancellationReason?: ParticipationCancellationReason;
  rejectionReason?: string;
  at?: Date;
}

export const changeParticipation = (participation: ParticipationRec, status: ParticipationStatus, options: ChangeOptions) => {
  const at = options.at ?? now();
  participation.status = status;
  participation.cancellationReason =
    status === ParticipationStatus.CANCELLED
      ? (options.cancellationReason ?? ParticipationCancellationReason.UNSPECIFIED)
      : ParticipationCancellationReason.UNSPECIFIED;
  if (status === ParticipationStatus.REJECTED) participation.rejectionReason = options.rejectionReason;
  participation.updatedAt = at;
  const sequence = db.participationChanges.filter((item) => item.participationId === participation.id).length + 1;
  db.participationChanges.push({
    participationId: participation.id,
    sequence,
    status,
    cancellationReason: participation.cancellationReason,
    rejectionReason: status === ParticipationStatus.REJECTED ? options.rejectionReason : undefined,
    actor: options.actor,
    createdAt: at,
  });
};

export const latestChatSequence = (eventId: string) => db.chatSequences.get(eventId) ?? 0;

export const openChatAccess = (eventId: string, volunteerId: string) => {
  chatHub.revokeVolunteer(eventId, volunteerId);
  db.chatAccess.set(chatAccessKey(eventId, volunteerId), { visibleAfterSequence: latestChatSequence(eventId), open: true });
};

export const closeChatAccess = (eventId: string, volunteerId: string) => {
  const access = db.chatAccess.get(chatAccessKey(eventId, volunteerId));
  if (access) access.open = false;
  chatHub.revokeVolunteer(eventId, volunteerId);
};

export const approvedCount = (eventId: string) =>
  [...db.participations.values()].filter((item) => item.eventId === eventId && item.status === ParticipationStatus.APPROVED).length;

export const pendingCount = (eventId: string) =>
  [...db.participations.values()].filter((item) => item.eventId === eventId && item.status === ParticipationStatus.PENDING).length;

export const hasOpenApplications = (eventId: string) => approvedCount(eventId) + pendingCount(eventId) > 0;

export const cancelRewardOrder = (order: RewardOrderRec, reason: RewardOrderCancellationReason) => {
  if (order.status === RewardOrderStatus.CANCELLED) return;
  if (order.status === RewardOrderStatus.ISSUED) fail(Code.FailedPrecondition, 'Выданный заказ нельзя отменить');
  recordOperation({
    volunteerId: order.volunteerId,
    coordinatorId: order.coordinatorId,
    type: VolunteerOperationType.REWARD_REFUND,
    source: { case: 'rewardOrderId', value: order.id },
    description: `Возврат за «${order.rewardTitle}»`,
    coinDelta: order.priceCoins,
    minuteDelta: 0,
  });
  order.status = RewardOrderStatus.CANCELLED;
  order.cancellationReason = reason;
  touch(order);
};

export const reservedQuantity = (rewardId: string) =>
  [...db.rewardOrders.values()].filter((item) => item.rewardId === rewardId && item.status === RewardOrderStatus.PENDING).length;

export const issuedQuantity = (rewardId: string) =>
  [...db.rewardOrders.values()].filter((item) => item.rewardId === rewardId && item.status === RewardOrderStatus.ISSUED).length;

export const terminateCooperation = (cooperation: CooperationRec) => {
  if (cooperation.status !== CooperationStatus.ACTIVE) fail(Code.FailedPrecondition, 'Сотрудничество не активно');
  cooperation.status = CooperationStatus.TERMINATED;
  cooperation.retryAfter = undefined;
  touch(cooperation);
  const coordinatorEvents = new Set(
    [...db.events.values()].filter((item) => item.coordinatorId === cooperation.coordinatorId).map((item) => item.id),
  );
  db.participations.forEach((participation) => {
    if (participation.volunteerId !== cooperation.volunteerId || !coordinatorEvents.has(participation.eventId)) return;
    if (participation.status === ParticipationStatus.PENDING || participation.status === ParticipationStatus.APPROVED) {
      changeParticipation(participation, ParticipationStatus.CANCELLED, {
        actor: SYSTEM,
        cancellationReason: ParticipationCancellationReason.COOPERATION_TERMINATED,
      });
    }
  });
  coordinatorEvents.forEach((eventId) => closeChatAccess(eventId, cooperation.volunteerId));
  db.rewardOrders.forEach((order) => {
    if (
      order.volunteerId === cooperation.volunteerId &&
      order.coordinatorId === cooperation.coordinatorId &&
      order.status === RewardOrderStatus.PENDING
    ) {
      cancelRewardOrder(order, RewardOrderCancellationReason.COOPERATION_TERMINATED);
    }
  });
};

export const cancelEventCascade = (event: EventRec) => {
  event.status = EventStatus.CANCELLED;
  db.participations.forEach((participation) => {
    if (participation.eventId !== event.id) return;
    if (participation.status === ParticipationStatus.PENDING || participation.status === ParticipationStatus.APPROVED) {
      changeParticipation(participation, ParticipationStatus.CANCELLED, {
        actor: SYSTEM,
        cancellationReason: ParticipationCancellationReason.EVENT_CANCELLED,
      });
      closeChatAccess(event.id, participation.volunteerId);
    }
  });
  chatHub.revokeVolunteersOfEvent(event.id);
};
