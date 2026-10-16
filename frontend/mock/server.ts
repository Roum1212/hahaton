import { randomBytes, randomUUID } from 'node:crypto';
import http from 'node:http';
import { create } from '@bufbuild/protobuf';
import { timestampNow } from '@bufbuild/protobuf/wkt';
import { Code, ConnectError, type ConnectRouter, type HandlerContext } from '@connectrpc/connect';
import { connectNodeAdapter } from '@connectrpc/connect-node';
import { AuthService } from '../genproto/auth/v1/auth_pb';
import { UserRole, UserSchema, UserService, type User } from '../genproto/user/v1/user_pb';

const PORT = Number(process.env.MOCK_PORT || 8080);
const FLIDA_USERINFO_URL = process.env.MOCK_FLIDA_USERINFO_URL || 'https://id.flida.ru/api/oidc/userinfo';
const ADMIN_EMAIL = (process.env.MOCK_ADMIN_EMAIL || 'admin@example.com').toLowerCase();
const ACCESS_TOKEN_TTL_MS = Number(process.env.MOCK_ACCESS_TTL_MS || 5 * 60 * 1000);
const REFRESH_COOKIE = 'refresh_token';

interface StoredUser {
  user: User;
  flidaSub?: string;
}

interface FlidaUserInfo {
  sub: string;
  email?: string;
  email_verified?: boolean;
  given_name?: string;
  family_name?: string;
  nickname?: string;
}

const users = new Map<string, StoredUser>();
const accessTokens = new Map<string, { userId: string; expiresAt: number }>();
const refreshTokens = new Map<string, string>();

const addUser = (data: Pick<User, 'email' | 'firstName' | 'lastName' | 'role'>, flidaSub?: string): User => {
  const user = create(UserSchema, { ...data, id: randomUUID(), active: true, createdAt: timestampNow() });
  users.set(user.id, { user, flidaSub });
  return user;
};

addUser({ email: ADMIN_EMAIL, firstName: 'Admin', lastName: 'Template', role: UserRole.ADMIN });
for (let i = 1; i <= 75; i++) {
  addUser({ email: `user${i}@example.com`, firstName: 'User', lastName: `${i}`, role: UserRole.USER });
}

const findByEmail = (email: string) => [...users.values()].find((item) => item.user.email.toLowerCase() === email.toLowerCase());

const fetchFlidaUser = async (token: string): Promise<FlidaUserInfo> => {
  const response = await fetch(FLIDA_USERINFO_URL, { headers: { Authorization: `Bearer ${token}` } }).catch(() => null);
  if (!response) throw new ConnectError('Флида ID недоступен', Code.Unavailable);
  if (response.status === 401) throw new ConnectError('Токен Флида ID недействителен', Code.Unauthenticated);
  if (!response.ok) throw new ConnectError('Флида ID недоступен', Code.Unavailable);
  const info = (await response.json()) as FlidaUserInfo;
  if (!info.sub) throw new ConnectError('Некорректный ответ Флида ID', Code.Internal);
  return info;
};

const resolveFlidaUser = (info: FlidaUserInfo): StoredUser => {
  const bySub = [...users.values()].find((item) => item.flidaSub === info.sub);
  if (bySub) return bySub;
  const byEmail = info.email && info.email_verified ? findByEmail(info.email) : undefined;
  if (byEmail && !byEmail.flidaSub) {
    byEmail.flidaSub = info.sub;
    return byEmail;
  }
  const email = info.email ?? '';
  const user = addUser(
    {
      email,
      firstName: info.given_name ?? info.nickname ?? '',
      lastName: info.family_name ?? '',
      role: email.toLowerCase() === ADMIN_EMAIL && info.email_verified ? UserRole.ADMIN : UserRole.USER,
    },
    info.sub,
  );
  return users.get(user.id)!;
};

const issueTokens = (ctx: HandlerContext, userId: string): string => {
  const access = randomBytes(24).toString('hex');
  const refresh = randomBytes(24).toString('hex');
  accessTokens.set(access, { userId, expiresAt: Date.now() + ACCESS_TOKEN_TTL_MS });
  refreshTokens.set(refresh, userId);
  ctx.responseHeader.append('Set-Cookie', `${REFRESH_COOKIE}=${refresh}; Path=/; HttpOnly; SameSite=Lax; Max-Age=2592000`);
  return access;
};

const readCookie = (ctx: HandlerContext, name: string): string | undefined =>
  ctx.requestHeader
    .get('cookie')
    ?.split(';')
    .map((part) => part.trim().split('='))
    .find(([key]) => key === name)?.[1];

const requireUser = (ctx: HandlerContext): User => {
  const token = ctx.requestHeader.get('authorization')?.replace(/^Bearer\s+/i, '');
  const session = token ? accessTokens.get(token) : undefined;
  if (!session || session.expiresAt < Date.now()) throw new ConnectError('Токен недействителен', Code.Unauthenticated);
  const stored = users.get(session.userId);
  if (!stored || !stored.user.active) throw new ConnectError('Пользователь заблокирован', Code.Unauthenticated);
  return stored.user;
};

const requireAdmin = (ctx: HandlerContext): User => {
  const user = requireUser(ctx);
  if (user.role !== UserRole.ADMIN) throw new ConnectError('Недостаточно прав', Code.PermissionDenied);
  return user;
};

const getStored = (id: string): StoredUser => {
  const stored = users.get(id);
  if (!stored) throw new ConnectError('Пользователь не найден', Code.NotFound);
  return stored;
};

const routes = (router: ConnectRouter) => {
  router.service(AuthService, {
    async loginWithFlida(req, ctx) {
      const stored = resolveFlidaUser(await fetchFlidaUser(req.flidaAccessToken));
      if (!stored.user.active) throw new ConnectError('Пользователь заблокирован', Code.PermissionDenied);
      return { accessToken: issueTokens(ctx, stored.user.id) };
    },
    refreshToken(_, ctx) {
      const refresh = readCookie(ctx, REFRESH_COOKIE);
      const userId = refresh ? refreshTokens.get(refresh) : undefined;
      if (!refresh || !userId) throw new ConnectError('Сессия истекла', Code.Unauthenticated);
      refreshTokens.delete(refresh);
      return { accessToken: issueTokens(ctx, userId) };
    },
    logout(_, ctx) {
      const refresh = readCookie(ctx, REFRESH_COOKIE);
      if (refresh) refreshTokens.delete(refresh);
      ctx.responseHeader.append('Set-Cookie', `${REFRESH_COOKIE}=; Path=/; HttpOnly; SameSite=Lax; Max-Age=0`);
      return {};
    },
    me(_, ctx) {
      return { user: requireUser(ctx) };
    },
  });

  router.service(UserService, {
    listUsers(req, ctx) {
      requireAdmin(ctx);
      const search = req.search.trim().toLowerCase();
      const filtered = [...users.values()]
        .map((item) => item.user)
        .filter((user) => !search || `${user.email} ${user.firstName} ${user.lastName}`.toLowerCase().includes(search));
      const offset = req.page?.offset ?? 0;
      const limit = req.page?.limit || 50;
      return { users: filtered.slice(offset, offset + limit), total: filtered.length };
    },
    getUser(req, ctx) {
      requireUser(ctx);
      return { user: getStored(req.id).user };
    },
    createUser(req, ctx) {
      requireAdmin(ctx);
      if (findByEmail(req.email)) throw new ConnectError('Пользователь с таким email уже существует', Code.AlreadyExists);
      return { user: addUser(req) };
    },
    updateUser(req, ctx) {
      const current = requireUser(ctx);
      if (current.id !== req.id && current.role !== UserRole.ADMIN) throw new ConnectError('Недостаточно прав', Code.PermissionDenied);
      const stored = getStored(req.id);
      const role = current.role === UserRole.ADMIN && req.role !== UserRole.UNSPECIFIED ? req.role : stored.user.role;
      const active = current.id === req.id ? stored.user.active : req.active;
      stored.user = create(UserSchema, {
        ...stored.user,
        email: req.email,
        firstName: req.firstName,
        lastName: req.lastName,
        role,
        active,
      });
      return { user: stored.user };
    },
    deleteUser(req, ctx) {
      const current = requireAdmin(ctx);
      if (current.id === req.id) throw new ConnectError('Нельзя удалить самого себя', Code.FailedPrecondition);
      getStored(req.id);
      users.delete(req.id);
      return {};
    },
  });
};

const handler = connectNodeAdapter({ routes });

http
  .createServer((req, res) => {
    const origin = req.headers.origin;
    if (origin) {
      res.setHeader('Access-Control-Allow-Origin', origin);
      res.setHeader('Access-Control-Allow-Credentials', 'true');
      res.setHeader('Vary', 'Origin');
    }
    if (req.method === 'OPTIONS') {
      res.setHeader('Access-Control-Allow-Methods', 'POST, GET, OPTIONS');
      res.setHeader('Access-Control-Allow-Headers', req.headers['access-control-request-headers'] || '*');
      res.setHeader('Access-Control-Max-Age', '7200');
      res.writeHead(204).end();
      return;
    }
    handler(req, res);
  })
  .listen(PORT, () => {
    console.log(`Mock Connect RPC server: http://localhost:${PORT}`);
    console.log(`Admin email: ${ADMIN_EMAIL}`);
  });
