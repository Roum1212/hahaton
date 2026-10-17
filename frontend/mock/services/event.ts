import { timestampDate, type Timestamp } from '@bufbuild/protobuf/wkt';
import { Code, type ConnectRouter } from '@connectrpc/connect';
import {
  EventCandidateRequirement,
  EventListScope,
  EventService,
  EventStatus,
  type EventParticipationResult,
} from '../../genproto/event/v1/event_service_pb';
import { ParticipationStatus } from '../../genproto/participation/v1/participation_service_pb';
import { VolunteerOperationType } from '../../genproto/volunteer/v1/volunteer_service_pb';
import { activePartnerIds, hasActiveCooperation, requireCoordinator, requireUser } from '../auth';
import { chatHub } from '../chat-hub';
import { coinsForMinutes, compareStrings, db, fail, newId, now, paginate, requireFound, type EventRec, type VolunteerRec } from '../db';
import { cancelEventCascade, changeParticipation, hasOpenApplications, recordOperation } from '../domain';
import { resolveImagePath } from '../images';
import { toEvent } from '../mappers';
import { buildAreaFilter, validatePoint } from './geo';

const AGE_REQUIREMENTS = [
  EventCandidateRequirement.AGE_14_PLUS,
  EventCandidateRequirement.AGE_16_PLUS,
  EventCandidateRequirement.AGE_18_PLUS,
];
const KNOWN_REQUIREMENTS = Object.values(EventCandidateRequirement).filter(
  (value): value is EventCandidateRequirement => typeof value === 'number' && value !== EventCandidateRequirement.UNSPECIFIED,
);

const validateRequirements = (values: EventCandidateRequirement[]) => {
  if (values.some((value) => !KNOWN_REQUIREMENTS.includes(value))) return fail(Code.InvalidArgument, 'Неизвестное требование');
  if (new Set(values).size !== values.length) return fail(Code.InvalidArgument, 'Требования повторяются');
  if (values.filter((value) => AGE_REQUIREMENTS.includes(value)).length > 1)
    return fail(Code.InvalidArgument, 'Допустим один возрастной порог');
  return values;
};

const validateInterval = (startsAt: Date, endsAt: Date) => {
  if (startsAt.getTime() <= Date.now()) return fail(Code.InvalidArgument, 'Начало должно быть в будущем');
  const duration = endsAt.getTime() - startsAt.getTime();
  if (duration <= 0) return fail(Code.InvalidArgument, 'Окончание должно быть позже начала');
  if (duration % 1_800_000 !== 0) return fail(Code.InvalidArgument, 'Длительность должна быть кратна 30 минутам');
};

const requireDate = (value: Timestamp | undefined, label: string): Date => {
  if (!value) return fail(Code.InvalidArgument, `Не указано ${label}`);
  return timestampDate(value);
};

export const canReadEvent = (event: EventRec, user: VolunteerRec) =>
  event.coordinatorId === user.id ||
  hasActiveCooperation(event.coordinatorId, user.id) ||
  [...db.participations.values()].some((item) => item.eventId === event.id && item.volunteerId === user.id);

export const requireOwnEvent = (eventId: string, coordinatorId: string): EventRec => {
  const event = requireFound(db.events.get(eventId), 'Мероприятие не найдено');
  if (event.coordinatorId !== coordinatorId) return fail(Code.PermissionDenied, 'Чужое мероприятие');
  return event;
};

const intervalFilter = (from: Timestamp | undefined, to: Timestamp | undefined) => {
  const fromDate = from ? timestampDate(from) : undefined;
  const toDate = to ? timestampDate(to) : undefined;
  if (fromDate && toDate && fromDate.getTime() >= toDate.getTime())
    return fail(Code.InvalidArgument, 'Начало интервала должно быть раньше конца');
  return (value: Date) => (!fromDate || value.getTime() >= fromDate.getTime()) && (!toDate || value.getTime() < toDate.getTime());
};

const resultsKey = (results: EventParticipationResult[]) =>
  JSON.stringify(
    results
      .map((item) => [item.participationId, item.outcome.case, String(item.outcome.value)])
      .sort((a, b) => compareStrings(a[0]!, b[0]!)),
  );

export const registerEventService = (router: ConnectRouter) =>
  router.service(EventService, {
    createEvent(req, ctx) {
      const coordinator = requireCoordinator(ctx);
      const title = req.title.trim();
      if (!title) return fail(Code.InvalidArgument, 'Укажите название');
      const point = validatePoint(req.point);
      const startsAt = requireDate(req.startsAt, 'начало');
      const endsAt = requireDate(req.endsAt, 'окончание');
      validateInterval(startsAt, endsAt);
      if (req.volunteerLimit <= 0) return fail(Code.InvalidArgument, 'Лимит участников должен быть больше нуля');
      const event: EventRec = {
        id: newId(),
        coordinatorId: coordinator.id,
        title,
        description: req.description,
        location: req.location,
        point,
        startsAt,
        endsAt,
        requirements: validateRequirements(req.candidateRequirements),
        volunteerLimit: req.volunteerLimit,
        imagePath: resolveImagePath(req.imagePath, coordinator.id),
        status: EventStatus.PUBLISHED,
      };
      db.events.set(event.id, event);
      return { id: event.id };
    },
    getEvent(req, ctx) {
      const user = requireUser(ctx);
      const event = requireFound(db.events.get(req.id), 'Мероприятие не найдено');
      if (!canReadEvent(event, user)) return fail(Code.PermissionDenied, 'Нет доступа к мероприятию');
      return { event: toEvent(event) };
    },
    listEvents(req, ctx) {
      const user = requireUser(ctx);
      let owners: Set<string>;
      if (req.scope === EventListScope.CATALOG) {
        if (!req.searchArea) return fail(Code.InvalidArgument, 'Для каталога нужна область поиска');
        owners = activePartnerIds(user.id);
      } else if (req.scope === EventListScope.OWNED) {
        owners = new Set([requireCoordinator(ctx).id]);
      } else {
        return fail(Code.InvalidArgument, 'Не указана область выборки');
      }
      if (req.coordinatorId !== undefined && !req.coordinatorId) return fail(Code.InvalidArgument, 'Пустой фильтр координатора');
      if (req.status === EventStatus.UNSPECIFIED) return fail(Code.InvalidArgument, 'Некорректный статус');
      const inArea = req.searchArea ? buildAreaFilter(req.searchArea) : () => true;
      const startsIn = intervalFilter(req.startsAtFrom, req.startsAtTo);
      const endsIn = intervalFilter(req.endsAtFrom, req.endsAtTo);
      const items = [...db.events.values()]
        .filter((item) => owners.has(item.coordinatorId))
        .filter((item) => !req.coordinatorId || item.coordinatorId === req.coordinatorId)
        .filter((item) => req.status === undefined || item.status === req.status)
        .filter((item) => startsIn(item.startsAt) && endsIn(item.endsAt) && inArea(item.point))
        .sort((a, b) => a.startsAt.getTime() - b.startsAt.getTime() || compareStrings(a.id, b.id));
      const { page, totalCount } = paginate(items, req.limit, req.offset);
      return { events: page.map(toEvent), totalCount };
    },
    updateEvent(req, ctx) {
      const coordinator = requireCoordinator(ctx);
      const event = requireOwnEvent(req.id, coordinator.id);
      if (event.status !== EventStatus.PUBLISHED || event.startsAt.getTime() <= Date.now()) {
        return fail(Code.FailedPrecondition, 'Мероприятие уже нельзя редактировать');
      }
      if (hasOpenApplications(event.id)) return fail(Code.FailedPrecondition, 'Есть активные заявки, правка запрещена');
      if (req.title !== undefined && !req.title.trim()) return fail(Code.InvalidArgument, 'Название не может быть пустым');
      if (req.volunteerLimit !== undefined && req.volunteerLimit <= 0)
        return fail(Code.InvalidArgument, 'Лимит участников должен быть больше нуля');
      const point = req.point ? validatePoint(req.point) : event.point;
      const startsAt = req.startsAt ? timestampDate(req.startsAt) : event.startsAt;
      const endsAt = req.endsAt ? timestampDate(req.endsAt) : event.endsAt;
      validateInterval(startsAt, endsAt);
      const requirements = req.candidateRequirements ? validateRequirements(req.candidateRequirements.values) : event.requirements;
      const imagePath = req.imagePath === undefined ? event.imagePath : resolveImagePath(req.imagePath, coordinator.id);
      Object.assign(event, {
        title: req.title?.trim() ?? event.title,
        description: req.description ?? event.description,
        location: req.location ?? event.location,
        point,
        startsAt,
        endsAt,
        requirements,
        volunteerLimit: req.volunteerLimit ?? event.volunteerLimit,
        imagePath,
      });
      return {};
    },
    cancelEvent(req, ctx) {
      const coordinator = requireCoordinator(ctx);
      const event = requireOwnEvent(req.id, coordinator.id);
      if (event.status !== EventStatus.PUBLISHED || event.startsAt.getTime() <= Date.now()) {
        return fail(Code.FailedPrecondition, 'Отменить можно только до начала опубликованного мероприятия');
      }
      cancelEventCascade(event);
      return {};
    },
    completeEvent(req, ctx) {
      const coordinator = requireCoordinator(ctx);
      if (!req.requestId) return fail(Code.InvalidArgument, 'Нужен request_id');
      const event = requireOwnEvent(req.id, coordinator.id);
      const payloadKey = `${event.id}|${resultsKey(req.results)}`;
      const requestKey = `${coordinator.id}:${req.requestId}`;
      const previous = db.completionResults.get(requestKey);
      if (previous !== undefined) {
        if (previous !== payloadKey) return fail(Code.AlreadyExists, 'request_id уже использован с другими данными');
        return { event: toEvent(event) };
      }
      if (event.status === EventStatus.COMPLETED) return fail(Code.FailedPrecondition, 'Мероприятие уже завершено');
      if (event.status !== EventStatus.PUBLISHED) return fail(Code.FailedPrecondition, 'Мероприятие не опубликовано');
      if (Date.now() < event.endsAt.getTime()) return fail(Code.FailedPrecondition, 'Завершить можно только после окончания');

      const approved = [...db.participations.values()].filter(
        (item) => item.eventId === event.id && item.status === ParticipationStatus.APPROVED,
      );
      const approvedIds = new Set(approved.map((item) => item.id));
      const seen = new Set<string>();
      for (const result of req.results) {
        if (!result.participationId || seen.has(result.participationId))
          return fail(Code.InvalidArgument, 'Повтор или пустая строка итогов');
        seen.add(result.participationId);
        if (!approvedIds.has(result.participationId)) return fail(Code.InvalidArgument, 'Строка итогов не относится к принятому участнику');
        if (result.outcome.case === 'absent' && !result.outcome.value) return fail(Code.InvalidArgument, 'absent должен быть true');
        if (result.outcome.case === 'confirmedMinutes') {
          const minutes = Number(result.outcome.value);
          if (minutes <= 0 || minutes % 30 !== 0) return fail(Code.InvalidArgument, 'Минуты должны быть положительными и кратными 30');
        }
        if (!result.outcome.case) return fail(Code.InvalidArgument, 'Не указан итог участия');
      }
      if (seen.size !== approvedIds.size) return fail(Code.InvalidArgument, 'Нужны итоги по всем принятым участникам');

      const at = now();
      for (const result of req.results) {
        const participation = db.participations.get(result.participationId)!;
        if (result.outcome.case === 'confirmedMinutes') {
          const minutes = Number(result.outcome.value);
          const coins = coinsForMinutes(minutes);
          const operation = recordOperation({
            volunteerId: participation.volunteerId,
            coordinatorId: coordinator.id,
            type: VolunteerOperationType.PARTICIPATION_ACCRUAL,
            source: { case: 'participationId', value: participation.id },
            description: `Участие в «${event.title}»`,
            coinDelta: coins,
            minuteDelta: minutes,
            createdAt: at,
          });
          participation.accrualOperationId = operation.id;
          participation.awardedCoins = coins;
          participation.awardedMinutes = minutes;
          changeParticipation(participation, ParticipationStatus.CONFIRMED, {
            actor: { case: 'coordinatorId', value: coordinator.id },
            at,
          });
        } else {
          changeParticipation(participation, ParticipationStatus.NOT_ATTENDED, {
            actor: { case: 'coordinatorId', value: coordinator.id },
            at,
          });
        }
      }
      db.participations.forEach((item) => {
        if (item.eventId === event.id && item.status === ParticipationStatus.PENDING) {
          changeParticipation(item, ParticipationStatus.CLOSED, { actor: { case: undefined }, at });
        }
      });
      event.status = EventStatus.COMPLETED;
      event.completion = { requestId: req.requestId, payloadKey };
      db.completionResults.set(requestKey, payloadKey);
      chatHub.completeEvent(event.id);
      return { event: toEvent(event) };
    },
  });
