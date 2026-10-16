import { create } from '@bufbuild/protobuf';
import { timestampFromDate } from '@bufbuild/protobuf/wkt';
import { Code, type ConnectRouter } from '@connectrpc/connect';
import { SessionSchema, SessionService } from '../../genproto/session/v1/session_service_pb';
import { clearSessionCookie, createSession, findSession, requireSession, setSessionCookie } from '../auth';
import { chatHub } from '../chat-hub';
import { db, fail, newId, type DateParts, type SessionRec, type VolunteerRec } from '../db';
import { provisionCoordinator } from '../seed';

interface IdTokenClaims {
  iss?: string;
  sub?: string;
  exp?: number;
  given_name?: string;
  family_name?: string;
  middle_name?: string;
  email?: string;
  phone_number?: string;
  birthdate?: string;
  date_of_birth?: string;
}

const COORDINATOR_EMAILS = (process.env.MOCK_COORDINATOR_EMAILS || '')
  .split(',')
  .map((item) => item.trim().toLowerCase())
  .filter(Boolean);

const toSession = (session: SessionRec) =>
  create(SessionSchema, {
    userId: session.userId,
    hasCoordinatorProfile: !!session.userId && db.coordinators.has(session.userId),
    expiresAt: timestampFromDate(session.expiresAt),
  });

const decodeIdToken = (token: string): IdTokenClaims => {
  const payload = token.split('.')[1];
  if (!payload) return fail(Code.InvalidArgument, 'Некорректный ID token');
  try {
    return JSON.parse(Buffer.from(payload, 'base64url').toString('utf8')) as IdTokenClaims;
  } catch {
    return fail(Code.InvalidArgument, 'Некорректный ID token');
  }
};

const parseDate = (value: string | undefined): DateParts | undefined => {
  const match = value?.match(/^(\d{4})-(\d{2})-(\d{2})$/);
  return match ? { year: Number(match[1]), month: Number(match[2]), day: Number(match[3]) } : undefined;
};

const register = (claims: IdTokenClaims & { iss: string; sub: string }): VolunteerRec => {
  const dateOfBirth = parseDate(claims.birthdate ?? claims.date_of_birth);
  const missing = [
    !claims.given_name && 'given_name',
    !claims.family_name && 'family_name',
    !claims.email && 'email',
    !claims.phone_number && 'phone_number',
    !dateOfBirth && 'birthdate',
  ].filter(Boolean);
  if (missing.length) console.warn(`[mock] В ID token нет ${missing.join(', ')}; подставлены учебные значения`);
  const volunteer: VolunteerRec = {
    id: newId(),
    iss: claims.iss,
    sub: claims.sub,
    firstName: claims.given_name || 'Волонтёр',
    lastName: claims.family_name || 'Обнинский',
    middleName: claims.middle_name ?? '',
    email: claims.email || `${claims.sub}@example.com`,
    phone: claims.phone_number || '+70000000000',
    dateOfBirth: dateOfBirth ?? { year: 2000, month: 1, day: 1 },
    coinBalance: 0,
    confirmedMinutes: 0,
    operationSequence: 0,
  };
  db.volunteers.set(volunteer.id, volunteer);
  if (COORDINATOR_EMAILS.includes(volunteer.email.toLowerCase())) provisionCoordinator(volunteer.id);
  return volunteer;
};

const resolveVolunteer = (claims: IdTokenClaims): VolunteerRec => {
  if (!claims.iss || !claims.sub) return fail(Code.InvalidArgument, 'В ID token нет iss или sub');
  if (claims.exp && claims.exp * 1000 < Date.now()) return fail(Code.Unauthenticated, 'ID token истёк');
  const existing = [...db.volunteers.values()].find((item) => item.iss === claims.iss && item.sub === claims.sub);
  return existing ?? register({ ...claims, iss: claims.iss, sub: claims.sub });
};

export const registerSessionService = (router: ConnectRouter) =>
  router.service(SessionService, {
    openSession(_, ctx) {
      const session = findSession(ctx) ?? createSession();
      setSessionCookie(ctx, session);
      return { session: toSession(session) };
    },
    authenticateSession(req, ctx) {
      const current = requireSession(ctx);
      if (current.userId) return fail(Code.FailedPrecondition, 'Сессия уже авторизована');
      if (!req.idToken) return fail(Code.InvalidArgument, 'Нет ID token');
      const volunteer = resolveVolunteer(decodeIdToken(req.idToken));
      current.revoked = true;
      const session = createSession(volunteer.id);
      setSessionCookie(ctx, session);
      return { session: toSession(session) };
    },
    getSession(_, ctx) {
      return { session: toSession(requireSession(ctx)) };
    },
    closeSession(_, ctx) {
      const session = findSession(ctx);
      if (session) {
        session.revoked = true;
        chatHub.closeSession(session.id);
      }
      clearSessionCookie(ctx);
      return {};
    },
  });
