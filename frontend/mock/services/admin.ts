import { Code, type ConnectRouter, type HandlerContext } from '@connectrpc/connect';
import { AdminService } from '../../genproto/admin/v1/admin_service_pb';
import { isAdminRequest } from '../auth';
import { compareStrings, db, fail, paginate, requireFound } from '../db';
import { toAdminVolunteer } from '../mappers';
import { provisionCoordinator } from '../seed';
import { isValidAvatarUrl } from './coordinator';

const requireAdmin = (ctx: HandlerContext) => {
  if (!isAdminRequest(ctx.requestHeader.get('authorization') ?? undefined)) fail(Code.Unauthenticated, 'Нужна служебная авторизация');
};

const normalizePhone = (value: string) => {
  const digits = value.replace(/\D/g, '');
  return digits ? `+${digits.startsWith('8') && digits.length === 11 ? `7${digits.slice(1)}` : digits}` : '';
};

export const registerAdminService = (router: ConnectRouter) =>
  router.service(AdminService, {
    findVolunteers(req, ctx) {
      requireAdmin(ctx);
      const { criterion } = req;
      if (!criterion.case || !criterion.value.trim()) return fail(Code.InvalidArgument, 'Нужен ровно один непустой критерий');
      const value = criterion.value.trim();
      const items = [...db.volunteers.values()]
        .filter((item) => {
          if (criterion.case === 'emailAddress') return item.email.toLowerCase() === value.toLowerCase();
          if (criterion.case === 'phoneNumber') return normalizePhone(item.phone) === normalizePhone(value);
          return item.id === value;
        })
        .sort((a, b) => compareStrings(a.id, b.id));
      const { page, totalCount } = paginate(items, req.limit, req.offset);
      return { volunteers: page.map(toAdminVolunteer), totalCount };
    },
    createCoordinator(req, ctx) {
      requireAdmin(ctx);
      if (!req.volunteerId) return fail(Code.InvalidArgument, 'Не указан волонтёр');
      requireFound(db.volunteers.get(req.volunteerId), 'Волонтёр не найден');
      if (db.coordinators.has(req.volunteerId)) return fail(Code.AlreadyExists, 'Профиль координатора уже существует');
      if (req.cooperationRequestCooldownSeconds === undefined || req.cooperationRequestCooldownSeconds < 0n) {
        return fail(Code.InvalidArgument, 'Нужен неотрицательный период ожидания');
      }
      if (!isValidAvatarUrl(req.avatarUrl)) return fail(Code.InvalidArgument, 'Аватар должен быть HTTPS-ссылкой');
      provisionCoordinator(req.volunteerId, {
        avatarUrl: req.avatarUrl,
        description: req.description,
        cooldownSeconds: Number(req.cooperationRequestCooldownSeconds),
      });
      return { id: req.volunteerId };
    },
  });
