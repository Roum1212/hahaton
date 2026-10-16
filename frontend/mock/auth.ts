import { randomBytes } from 'node:crypto';
import { Code, type HandlerContext } from '@connectrpc/connect';
import { CooperationStatus } from '../genproto/cooperation/v1/cooperation_service_pb';
import { cooperationKey, db, fail, now, type CoordinatorRec, type SessionRec, type VolunteerRec } from './db';

export const SESSION_COOKIE = 'session';
const SESSION_TTL_MS = Number(process.env.MOCK_SESSION_TTL_MS || 7 * 24 * 60 * 60 * 1000);

const readCookie = (ctx: HandlerContext, name: string): string | undefined =>
  ctx.requestHeader
    .get('cookie')
    ?.split(';')
    .map((part) => part.trim().split('='))
    .find(([key]) => key === name)?.[1];

export const setSessionCookie = (ctx: HandlerContext, session: SessionRec) => {
  const maxAge = Math.max(0, Math.floor((session.expiresAt.getTime() - Date.now()) / 1000));
  ctx.responseHeader.append('Set-Cookie', `${SESSION_COOKIE}=${session.id}; Path=/; HttpOnly; SameSite=Lax; Max-Age=${maxAge}`);
};

export const clearSessionCookie = (ctx: HandlerContext) => {
  ctx.responseHeader.append('Set-Cookie', `${SESSION_COOKIE}=; Path=/; HttpOnly; SameSite=Lax; Max-Age=0`);
};

export const createSession = (userId?: string): SessionRec => {
  const session: SessionRec = {
    id: randomBytes(24).toString('hex'),
    userId,
    expiresAt: new Date(Date.now() + SESSION_TTL_MS),
    revoked: false,
  };
  db.sessions.set(session.id, session);
  return session;
};

export const isSessionValid = (session: SessionRec | undefined): session is SessionRec =>
  !!session && !session.revoked && session.expiresAt.getTime() > Date.now();

export const findSession = (ctx: HandlerContext): SessionRec | undefined => {
  const id = readCookie(ctx, SESSION_COOKIE);
  const session = id ? db.sessions.get(id) : undefined;
  return isSessionValid(session) ? session : undefined;
};

export const requireSession = (ctx: HandlerContext): SessionRec => {
  const session = findSession(ctx);
  if (!session) fail(Code.Unauthenticated, 'Сессия недействительна');
  return session as SessionRec;
};

export const requireUser = (ctx: HandlerContext): VolunteerRec => {
  const { userId } = requireSession(ctx);
  if (!userId) fail(Code.Unauthenticated, 'Требуется вход');
  const volunteer = db.volunteers.get(userId as string);
  if (!volunteer) fail(Code.Unauthenticated, 'Участник не найден');
  return volunteer as VolunteerRec;
};

export const requireCoordinator = (ctx: HandlerContext): CoordinatorRec => {
  const user = requireUser(ctx);
  const coordinator = db.coordinators.get(user.id);
  if (!coordinator) fail(Code.PermissionDenied, 'Требуется профиль координатора');
  return coordinator as CoordinatorRec;
};

export const hasActiveCooperation = (coordinatorId: string, volunteerId: string) =>
  db.cooperations.get(cooperationKey(coordinatorId, volunteerId))?.status === CooperationStatus.ACTIVE;

export const activePartnerIds = (volunteerId: string): Set<string> =>
  new Set(
    [...db.cooperations.values()]
      .filter((item) => item.volunteerId === volunteerId && item.status === CooperationStatus.ACTIVE)
      .map((item) => item.coordinatorId),
  );

export const isAdminRequest = (authorization: string | undefined) => {
  const token = process.env.MOCK_ADMIN_TOKEN || 'mock-admin-token';
  return authorization === `Bearer ${token}`;
};

export const touch = <T extends { updatedAt: Date }>(item: T): T => {
  item.updatedAt = now();
  return item;
};
