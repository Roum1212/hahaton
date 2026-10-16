import { timestampDate } from '@bufbuild/protobuf/wkt';
import { Code, type ConnectRouter } from '@connectrpc/connect';
import { CooperationListScope, CooperationService, CooperationStatus } from '../../genproto/cooperation/v1/cooperation_service_pb';
import { requireCoordinator, requireUser, touch } from '../auth';
import { compareStrings, cooperationKey, db, fail, newId, newToken, now, paginate, requireFound, type CooperationRec } from '../db';
import { terminateCooperation } from '../domain';
import { toCooperation, toInvitation } from '../mappers';

const createOrReset = (coordinatorId: string, volunteerId: string, status: CooperationStatus): CooperationRec => {
  const key = cooperationKey(coordinatorId, volunteerId);
  const at = now();
  const existing = db.cooperations.get(key);
  if (existing) {
    existing.status = status;
    existing.retryAfter = undefined;
    existing.updatedAt = at;
    return existing;
  }
  const created: CooperationRec = { coordinatorId, volunteerId, status, createdAt: at, updatedAt: at };
  db.cooperations.set(key, created);
  return created;
};

const requirePending = (coordinatorId: string, volunteerId: string) => {
  const cooperation = requireFound(db.cooperations.get(cooperationKey(coordinatorId, volunteerId)), 'Заявка не найдена');
  if (cooperation.status !== CooperationStatus.PENDING) return fail(Code.FailedPrecondition, 'Заявка уже рассмотрена');
  return cooperation;
};

const isInvitationActive = (expiresAt?: Date) => !expiresAt || Date.now() < expiresAt.getTime();

export const registerCooperationService = (router: ConnectRouter) =>
  router.service(CooperationService, {
    requestCooperation(req, ctx) {
      const user = requireUser(ctx);
      if (!req.coordinatorId) return fail(Code.InvalidArgument, 'Не указан координатор');
      if (req.coordinatorId === user.id) return fail(Code.InvalidArgument, 'Сотрудничество с собой запрещено');
      requireFound(db.coordinators.get(req.coordinatorId), 'Координатор не найден');
      const existing = db.cooperations.get(cooperationKey(req.coordinatorId, user.id));
      if (existing?.status === CooperationStatus.PENDING || existing?.status === CooperationStatus.ACTIVE) {
        return fail(Code.AlreadyExists, 'Заявка уже подана');
      }
      if (existing?.status === CooperationStatus.REJECTED && existing.retryAfter && Date.now() < existing.retryAfter.getTime()) {
        return fail(Code.FailedPrecondition, 'Повторная заявка пока недоступна');
      }
      return { cooperation: toCooperation(createOrReset(req.coordinatorId, user.id, CooperationStatus.PENDING)) };
    },
    approveCooperation(req, ctx) {
      const coordinator = requireCoordinator(ctx);
      const cooperation = requirePending(coordinator.id, req.volunteerId);
      cooperation.status = CooperationStatus.ACTIVE;
      touch(cooperation);
      return {};
    },
    rejectCooperation(req, ctx) {
      const coordinator = requireCoordinator(ctx);
      const cooperation = requirePending(coordinator.id, req.volunteerId);
      cooperation.status = CooperationStatus.REJECTED;
      cooperation.retryAfter = new Date(Date.now() + coordinator.cooldownSeconds * 1000);
      touch(cooperation);
      return {};
    },
    terminateCooperation(req, ctx) {
      const user = requireUser(ctx);
      const { counterpart } = req;
      if (!counterpart.case || !counterpart.value) return fail(Code.InvalidArgument, 'Нужен ровно один участник связи');
      if (counterpart.case === 'volunteerId') requireCoordinator(ctx);
      const key =
        counterpart.case === 'coordinatorId' ? cooperationKey(counterpart.value, user.id) : cooperationKey(user.id, counterpart.value);
      terminateCooperation(requireFound(db.cooperations.get(key), 'Сотрудничество не найдено'));
      return {};
    },
    getCooperation(req, ctx) {
      const user = requireUser(ctx);
      const { counterpart } = req;
      if (!counterpart.case || !counterpart.value) return fail(Code.InvalidArgument, 'Нужен ровно один участник связи');
      if (counterpart.case === 'volunteerId') requireCoordinator(ctx);
      const key =
        counterpart.case === 'coordinatorId' ? cooperationKey(counterpart.value, user.id) : cooperationKey(user.id, counterpart.value);
      return { cooperation: toCooperation(requireFound(db.cooperations.get(key), 'Сотрудничество не найдено')) };
    },
    listCooperations(req, ctx) {
      const user = requireUser(ctx);
      if (req.scope === CooperationListScope.AS_COORDINATOR) requireCoordinator(ctx);
      else if (req.scope !== CooperationListScope.AS_VOLUNTEER) return fail(Code.InvalidArgument, 'Не указана область выборки');
      if (req.status === CooperationStatus.UNSPECIFIED) return fail(Code.InvalidArgument, 'Некорректный статус');
      const items = [...db.cooperations.values()]
        .filter((item) => (req.scope === CooperationListScope.AS_VOLUNTEER ? item.volunteerId : item.coordinatorId) === user.id)
        .filter((item) => req.status === undefined || item.status === req.status)
        .sort((a, b) => compareStrings(a.coordinatorId, b.coordinatorId) || compareStrings(a.volunteerId, b.volunteerId));
      const { page, totalCount } = paginate(items, req.limit, req.offset);
      return { cooperations: page.map(toCooperation), totalCount };
    },
    createCooperationInvitation(req, ctx) {
      const coordinator = requireCoordinator(ctx);
      const expiresAt = req.expiresAt ? timestampDate(req.expiresAt) : undefined;
      if (expiresAt && expiresAt.getTime() <= Date.now()) return fail(Code.InvalidArgument, 'Срок должен быть в будущем');
      const invitation = { id: newId(), coordinatorId: coordinator.id, token: newToken(), createdAt: now(), expiresAt };
      db.invitations.set(invitation.id, invitation);
      return { invitation: toInvitation(invitation) };
    },
    listCooperationInvitations(req, ctx) {
      const coordinator = requireCoordinator(ctx);
      const items = [...db.invitations.values()]
        .filter((item) => item.coordinatorId === coordinator.id && (req.includeExpired || isInvitationActive(item.expiresAt)))
        .sort((a, b) => b.createdAt.getTime() - a.createdAt.getTime() || compareStrings(b.id, a.id));
      const { page, totalCount } = paginate(items, req.limit, req.offset);
      return { invitations: page.map(toInvitation), totalCount };
    },
    revokeCooperationInvitation(req, ctx) {
      const coordinator = requireCoordinator(ctx);
      const invitation = db.invitations.get(req.id);
      if (!invitation) return {};
      if (invitation.coordinatorId !== coordinator.id) return fail(Code.PermissionDenied, 'Чужое приглашение');
      db.invitations.delete(invitation.id);
      return {};
    },
    acceptCooperationInvitation(req, ctx) {
      const user = requireUser(ctx);
      if (!req.token) return fail(Code.InvalidArgument, 'Пустой токен');
      const invitation = requireFound(
        [...db.invitations.values()].find((item) => item.token === req.token),
        'Приглашение не найдено',
      );
      if (!isInvitationActive(invitation.expiresAt)) return fail(Code.FailedPrecondition, 'Срок приглашения истёк');
      if (invitation.coordinatorId === user.id) return fail(Code.InvalidArgument, 'Нельзя принять собственное приглашение');
      const existing = db.cooperations.get(cooperationKey(invitation.coordinatorId, user.id));
      if (existing?.status === CooperationStatus.ACTIVE) return { cooperation: toCooperation(existing) };
      return { cooperation: toCooperation(createOrReset(invitation.coordinatorId, user.id, CooperationStatus.ACTIVE)) };
    },
  });
