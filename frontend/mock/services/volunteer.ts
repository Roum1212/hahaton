import { Code, type ConnectRouter } from '@connectrpc/connect';
import type { Date as ProtoDate } from '../../genproto/google/type/date_pb';
import { VolunteerOperationType, VolunteerService } from '../../genproto/volunteer/v1/volunteer_service_pb';
import { hasActiveCooperation, requireCoordinator, requireUser } from '../auth';
import { cooperationKey, db, fail, paginate, requireFound } from '../db';
import { toCoordinatorVolunteer, toOperation, toVolunteer } from '../mappers';

const isValidBirthDate = (value: ProtoDate) => {
  const date = new Date(value.year, value.month - 1, value.day);
  return (
    date.getFullYear() === value.year && date.getMonth() === value.month - 1 && date.getDate() === value.day && date.getTime() <= Date.now()
  );
};

const isLinkedToCoordinator = (coordinatorId: string, volunteerId: string) => {
  if (db.cooperations.has(cooperationKey(coordinatorId, volunteerId))) return true;
  const ownEvents = new Set([...db.events.values()].filter((item) => item.coordinatorId === coordinatorId).map((item) => item.id));
  if ([...db.participations.values()].some((item) => item.volunteerId === volunteerId && ownEvents.has(item.eventId))) return true;
  return [...db.rewardOrders.values()].some((item) => item.volunteerId === volunteerId && item.coordinatorId === coordinatorId);
};

export const registerVolunteerService = (router: ConnectRouter) =>
  router.service(VolunteerService, {
    getVolunteer(req, ctx) {
      const user = requireUser(ctx);
      if (req.id !== undefined && !req.id) return fail(Code.InvalidArgument, 'Пустой id');
      if (req.id && req.id !== user.id) return fail(Code.PermissionDenied, 'Чужая карточка недоступна');
      return { volunteer: toVolunteer(user) };
    },
    getCoordinatorVolunteer(req, ctx) {
      const coordinator = requireCoordinator(ctx);
      if (!req.volunteerId) return fail(Code.InvalidArgument, 'Не указан волонтёр');
      const volunteer = requireFound(db.volunteers.get(req.volunteerId), 'Волонтёр не найден');
      if (!isLinkedToCoordinator(coordinator.id, volunteer.id)) return fail(Code.PermissionDenied, 'Нет связи с волонтёром');
      return { volunteer: toCoordinatorVolunteer(volunteer, hasActiveCooperation(coordinator.id, volunteer.id)) };
    },
    updateVolunteer(req, ctx) {
      const user = requireUser(ctx);
      if (req.firstName !== undefined && !req.firstName.trim()) return fail(Code.InvalidArgument, 'Имя не может быть пустым');
      if (req.lastName !== undefined && !req.lastName.trim()) return fail(Code.InvalidArgument, 'Фамилия не может быть пустой');
      if (req.dateOfBirth && !isValidBirthDate(req.dateOfBirth)) return fail(Code.InvalidArgument, 'Некорректная дата рождения');
      if (req.firstName !== undefined) user.firstName = req.firstName.trim();
      if (req.lastName !== undefined) user.lastName = req.lastName.trim();
      if (req.middleName !== undefined) user.middleName = req.middleName.trim();
      if (req.dateOfBirth) user.dateOfBirth = { year: req.dateOfBirth.year, month: req.dateOfBirth.month, day: req.dateOfBirth.day };
      return {};
    },
    listVolunteerOperations(req, ctx) {
      const user = requireUser(ctx);
      if (req.type !== undefined && !Object.values(VolunteerOperationType).includes(req.type)) {
        return fail(Code.InvalidArgument, 'Неизвестный тип операции');
      }
      if (req.type === VolunteerOperationType.UNSPECIFIED) return fail(Code.InvalidArgument, 'Тип операции не указан');
      const items = db.operations
        .filter((item) => item.volunteerId === user.id && (req.type === undefined || item.type === req.type))
        .sort((a, b) => b.sequence - a.sequence);
      const { page, totalCount } = paginate(items, req.limit, req.offset);
      return { operations: page.map(toOperation), totalCount };
    },
  });
