import { Code, type ConnectRouter } from '@connectrpc/connect';
import { EventStatus } from '../../genproto/event/v1/event_service_pb';
import {
  ParticipationCancellationReason,
  ParticipationListScope,
  ParticipationService,
  ParticipationStatus,
} from '../../genproto/participation/v1/participation_service_pb';
import { hasActiveCooperation, requireCoordinator, requireUser } from '../auth';
import {
  compareStrings,
  db,
  fail,
  newId,
  now,
  paginate,
  requireFound,
  type EventRec,
  type ParticipationRec,
  type VolunteerRec,
} from '../db';
import { approvedCount, changeParticipation, closeChatAccess, openChatAccess } from '../domain';
import { toParticipation, toParticipationChange } from '../mappers';
import { requireOwnEvent } from './event';

const requireApplicableEvent = (event: EventRec, volunteer: VolunteerRec) => {
  if (event.coordinatorId === volunteer.id) return fail(Code.FailedPrecondition, 'Нельзя записаться на своё мероприятие');
  if (!hasActiveCooperation(event.coordinatorId, volunteer.id)) return fail(Code.FailedPrecondition, 'Нужно активное сотрудничество');
  if (event.status !== EventStatus.PUBLISHED) return fail(Code.FailedPrecondition, 'Мероприятие не опубликовано');
  if (event.startsAt.getTime() <= Date.now()) return fail(Code.FailedPrecondition, 'Мероприятие уже началось');
};

const requireParticipation = (id: string) => requireFound(db.participations.get(id), 'Заявка не найдена');

const requireReadable = (participation: ParticipationRec, user: VolunteerRec) => {
  const event = db.events.get(participation.eventId);
  if (participation.volunteerId !== user.id && event?.coordinatorId !== user.id) return fail(Code.PermissionDenied, 'Нет доступа к заявке');
};

const volunteerActor = (id: string) => ({ case: 'volunteerId' as const, value: id });
const coordinatorActor = (id: string) => ({ case: 'coordinatorId' as const, value: id });

export const registerParticipationService = (router: ConnectRouter) =>
  router.service(ParticipationService, {
    requestParticipation(req, ctx) {
      const user = requireUser(ctx);
      const event = requireFound(db.events.get(req.eventId), 'Мероприятие не найдено');
      requireApplicableEvent(event, user);
      const exists = [...db.participations.values()].some((item) => item.eventId === event.id && item.volunteerId === user.id);
      if (exists) return fail(Code.AlreadyExists, 'Заявка на это мероприятие уже есть');
      const at = now();
      const participation: ParticipationRec = {
        id: newId(),
        eventId: event.id,
        volunteerId: user.id,
        status: ParticipationStatus.PENDING,
        cancellationReason: ParticipationCancellationReason.UNSPECIFIED,
        createdAt: at,
        updatedAt: at,
      };
      db.participations.set(participation.id, participation);
      changeParticipation(participation, ParticipationStatus.PENDING, { actor: volunteerActor(user.id), at });
      return { id: participation.id };
    },
    restoreParticipation(req, ctx) {
      const user = requireUser(ctx);
      const participation = requireParticipation(req.id);
      if (participation.volunteerId !== user.id) return fail(Code.PermissionDenied, 'Чужая заявка');
      if (participation.status !== ParticipationStatus.CANCELLED)
        return fail(Code.FailedPrecondition, 'Восстановить можно только отменённую заявку');
      requireApplicableEvent(db.events.get(participation.eventId)!, user);
      changeParticipation(participation, ParticipationStatus.PENDING, { actor: volunteerActor(user.id) });
      return {};
    },
    approveParticipation(req, ctx) {
      const coordinator = requireCoordinator(ctx);
      const participation = requireParticipation(req.id);
      const event = requireOwnEvent(participation.eventId, coordinator.id);
      if (participation.status !== ParticipationStatus.PENDING) return fail(Code.FailedPrecondition, 'Заявка уже рассмотрена');
      if (event.status !== EventStatus.PUBLISHED || event.startsAt.getTime() <= Date.now()) {
        return fail(Code.FailedPrecondition, 'Одобрить можно только до начала мероприятия');
      }
      if (!hasActiveCooperation(coordinator.id, participation.volunteerId))
        return fail(Code.FailedPrecondition, 'Нет активного сотрудничества');
      if (approvedCount(event.id) >= event.volunteerLimit) return fail(Code.FailedPrecondition, 'Свободных мест нет');
      changeParticipation(participation, ParticipationStatus.APPROVED, { actor: coordinatorActor(coordinator.id) });
      openChatAccess(event.id, participation.volunteerId);
      return {};
    },
    rejectParticipation(req, ctx) {
      const coordinator = requireCoordinator(ctx);
      const participation = requireParticipation(req.id);
      const event = requireOwnEvent(participation.eventId, coordinator.id);
      if (participation.status !== ParticipationStatus.PENDING) return fail(Code.FailedPrecondition, 'Заявка уже рассмотрена');
      if (event.status !== EventStatus.PUBLISHED) return fail(Code.FailedPrecondition, 'Мероприятие не опубликовано');
      const reason = req.reason?.trim();
      if (req.reason !== undefined && (!reason || [...reason].length > 2000)) {
        return fail(Code.InvalidArgument, 'Причина должна быть от 1 до 2000 символов');
      }
      changeParticipation(participation, ParticipationStatus.REJECTED, {
        actor: coordinatorActor(coordinator.id),
        rejectionReason: reason,
      });
      return {};
    },
    cancelParticipation(req, ctx) {
      const user = requireUser(ctx);
      const participation = requireParticipation(req.id);
      if (participation.volunteerId !== user.id) return fail(Code.PermissionDenied, 'Чужая заявка');
      if (participation.status !== ParticipationStatus.PENDING && participation.status !== ParticipationStatus.APPROVED) {
        return fail(Code.FailedPrecondition, 'Заявку нельзя отменить');
      }
      const event = db.events.get(participation.eventId)!;
      if (event.startsAt.getTime() <= Date.now()) return fail(Code.FailedPrecondition, 'Мероприятие уже началось');
      changeParticipation(participation, ParticipationStatus.CANCELLED, {
        actor: volunteerActor(user.id),
        cancellationReason: ParticipationCancellationReason.VOLUNTEER_CANCELLED,
      });
      closeChatAccess(event.id, user.id);
      return {};
    },
    getParticipation(req, ctx) {
      const user = requireUser(ctx);
      const participation = requireParticipation(req.id);
      requireReadable(participation, user);
      return { participation: toParticipation(participation) };
    },
    listParticipations(req, ctx) {
      const user = requireUser(ctx);
      let filter: (item: ParticipationRec) => boolean;
      if (req.scope === ParticipationListScope.MY_APPLICATIONS) {
        filter = (item) => item.volunteerId === user.id;
      } else if (req.scope === ParticipationListScope.MY_EVENTS) {
        const coordinator = requireCoordinator(ctx);
        filter = (item) => db.events.get(item.eventId)?.coordinatorId === coordinator.id;
      } else {
        return fail(Code.InvalidArgument, 'Не указана область выборки');
      }
      if (req.eventId !== undefined && !req.eventId) return fail(Code.InvalidArgument, 'Пустой фильтр мероприятия');
      if (req.volunteerId !== undefined && !req.volunteerId) return fail(Code.InvalidArgument, 'Пустой фильтр волонтёра');
      if (req.status === ParticipationStatus.UNSPECIFIED) return fail(Code.InvalidArgument, 'Некорректный статус');
      const items = [...db.participations.values()]
        .filter(filter)
        .filter((item) => (!req.eventId || item.eventId === req.eventId) && (!req.volunteerId || item.volunteerId === req.volunteerId))
        .filter((item) => req.status === undefined || item.status === req.status)
        .sort((a, b) => a.createdAt.getTime() - b.createdAt.getTime() || compareStrings(a.id, b.id));
      const { page, totalCount } = paginate(items, req.limit, req.offset);
      return { participations: page.map(toParticipation), totalCount };
    },
    listParticipationChanges(req, ctx) {
      const user = requireUser(ctx);
      const participation = requireParticipation(req.participationId);
      requireReadable(participation, user);
      const items = db.participationChanges
        .filter((item) => item.participationId === participation.id)
        .sort((a, b) => a.sequence - b.sequence);
      const { page, totalCount } = paginate(items, req.limit, req.offset);
      return { changes: page.map(toParticipationChange), totalCount };
    },
  });
