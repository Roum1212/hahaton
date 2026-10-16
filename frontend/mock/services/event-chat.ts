import { create } from '@bufbuild/protobuf';
import { Code, type ConnectRouter, type HandlerContext } from '@connectrpc/connect';
import { EventChatSchema, EventChatService } from '../../genproto/event-chat/v1/event_chat_service_pb';
import { EventStatus } from '../../genproto/event/v1/event_service_pb';
import { ParticipationStatus } from '../../genproto/participation/v1/participation_service_pb';
import { findSession, hasActiveCooperation, requireSession, requireUser } from '../auth';
import { ChatSubscriber, chatHub } from '../chat-hub';
import {
  chatAccessKey,
  compareStrings,
  db,
  fail,
  fullName,
  newId,
  now,
  paginate,
  requireFound,
  type EventRec,
  type VolunteerRec,
} from '../db';
import { latestChatSequence } from '../domain';
import { toChatMessage } from '../mappers';

interface ChatAccess {
  role: 'coordinator' | 'volunteer';
  visibleAfterSequence: number;
  readOnly: boolean;
}

const READABLE_STATUSES = [ParticipationStatus.APPROVED, ParticipationStatus.CONFIRMED, ParticipationStatus.NOT_ATTENDED];

const resolveAccess = (event: EventRec, user: VolunteerRec): ChatAccess | undefined => {
  if (event.coordinatorId === user.id) {
    return { role: 'coordinator', visibleAfterSequence: 0, readOnly: event.status !== EventStatus.PUBLISHED };
  }
  if (event.status === EventStatus.CANCELLED || !hasActiveCooperation(event.coordinatorId, user.id)) return undefined;
  const participation = [...db.participations.values()].find((item) => item.eventId === event.id && item.volunteerId === user.id);
  if (!participation || !READABLE_STATUSES.includes(participation.status)) return undefined;
  const access = db.chatAccess.get(chatAccessKey(event.id, user.id));
  if (!access?.open) return undefined;
  return {
    role: 'volunteer',
    visibleAfterSequence: access.visibleAfterSequence,
    readOnly: event.status !== EventStatus.PUBLISHED || participation.status !== ParticipationStatus.APPROVED,
  };
};

const requireAccess = (eventId: string, user: VolunteerRec) => {
  const event = requireFound(db.events.get(eventId), 'Мероприятие не найдено');
  const access = resolveAccess(event, user);
  if (!access) return fail(Code.PermissionDenied, 'Нет доступа к чату');
  return { event, access };
};

const toChat = (event: EventRec, access: ChatAccess) =>
  create(EventChatSchema, {
    eventId: event.id,
    eventTitle: event.title,
    readOnly: access.readOnly,
    latestSequence: BigInt(latestChatSequence(event.id)),
    visibleAfterSequence: BigInt(access.visibleAfterSequence),
  });

const visibleMessages = (eventId: string, after: number, upTo: number) =>
  db.chatMessages.filter((item) => item.eventId === eventId && item.sequence > after && item.sequence <= upTo);

const ensureStillAllowed = (ctx: HandlerContext, eventId: string, userId: string) => {
  if (!findSession(ctx)) return fail(Code.Unauthenticated, 'Сессия завершена');
  const user = db.volunteers.get(userId)!;
  return requireAccess(eventId, user);
};

export const registerEventChatService = (router: ConnectRouter) =>
  router.service(EventChatService, {
    getEventChat(req, ctx) {
      const user = requireUser(ctx);
      const { event, access } = requireAccess(req.eventId, user);
      return { chat: toChat(event, access) };
    },
    listMyEventChats(req, ctx) {
      const user = requireUser(ctx);
      const items = [...db.events.values()]
        .map((event) => ({ event, access: resolveAccess(event, user) }))
        .filter((item): item is { event: EventRec; access: ChatAccess } => !!item.access)
        .sort((a, b) => b.event.startsAt.getTime() - a.event.startsAt.getTime() || compareStrings(b.event.id, a.event.id));
      const { page, totalCount } = paginate(items, req.limit, req.offset);
      return { chats: page.map((item) => toChat(item.event, item.access)), totalCount };
    },
    sendEventChatMessage(req, ctx) {
      const user = requireUser(ctx);
      const { event, access } = requireAccess(req.eventId, user);
      if (!req.requestId) return fail(Code.InvalidArgument, 'Нужен request_id');
      const text = req.text;
      if (!text.trim() || [...text].length > 2000) return fail(Code.InvalidArgument, 'Сообщение должно быть от 1 до 2000 символов');
      const requestKey = `${access.role}:${user.id}:${event.id}:${req.requestId}`;
      const repeated = db.chatMessages.find((item) => item.requestKey === requestKey);
      if (repeated) {
        if (repeated.text !== text) return fail(Code.AlreadyExists, 'request_id уже использован для другого сообщения');
        if (repeated.sequence <= access.visibleAfterSequence) return fail(Code.FailedPrecondition, 'Сообщение вне доступной истории');
        return { message: toChatMessage(repeated) };
      }
      if (access.readOnly) return fail(Code.FailedPrecondition, 'Чат доступен только для чтения');
      const sequence = latestChatSequence(event.id) + 1;
      db.chatSequences.set(event.id, sequence);
      const message = {
        id: newId(),
        eventId: event.id,
        author:
          access.role === 'coordinator'
            ? { case: 'coordinatorId' as const, value: user.id }
            : { case: 'volunteerId' as const, value: user.id },
        authorDisplayName: fullName(user),
        text,
        sequence,
        requestKey,
        createdAt: now(),
      };
      db.chatMessages.push(message);
      chatHub.publish(message);
      return { message: toChatMessage(message) };
    },
    listEventChatMessages(req, ctx) {
      const user = requireUser(ctx);
      const { event, access } = requireAccess(req.eventId, user);
      const latest = latestChatSequence(event.id);
      const snapshot = req.snapshotSequence === undefined ? latest : Number(req.snapshotSequence);
      if (snapshot < access.visibleAfterSequence || snapshot > latest)
        return fail(Code.InvalidArgument, 'Граница истории вне доступного диапазона');
      const items = visibleMessages(event.id, access.visibleAfterSequence, snapshot).sort((a, b) => b.sequence - a.sequence);
      const { page, totalCount } = paginate(items, req.limit, req.offset);
      return {
        messages: page.map(toChatMessage),
        snapshotSequence: BigInt(snapshot),
        visibleAfterSequence: BigInt(access.visibleAfterSequence),
        totalCount,
      };
    },
    async *subscribeEventChatMessages(req, ctx) {
      const session = requireSession(ctx);
      const user = requireUser(ctx);
      const { event, access } = requireAccess(req.eventId, user);
      const after = Number(req.afterSequence);
      if (after > latestChatSequence(event.id)) fail(Code.InvalidArgument, 'after_sequence больше последнего сообщения');
      const subscriber = new ChatSubscriber(event.id, user.id, access.role, session.id);
      chatHub.add(subscriber);
      try {
        let lastSent = Math.max(after, access.visibleAfterSequence);
        for (const message of visibleMessages(event.id, lastSent, latestChatSequence(event.id))) {
          lastSent = message.sequence;
          yield { message: toChatMessage(message) };
        }
        if (event.status === EventStatus.COMPLETED) return;
        while (true) {
          const item = await subscriber.next(ctx.signal);
          if (item.kind === 'end') return;
          if (item.kind === 'error') throw item.error;
          if (item.message.sequence <= lastSent) continue;
          ensureStillAllowed(ctx, event.id, user.id);
          lastSent = item.message.sequence;
          yield { message: toChatMessage(item.message) };
        }
      } finally {
        chatHub.remove(subscriber);
      }
    },
  });
