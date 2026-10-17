import { ConnectError } from '@connectrpc/connect';
import { AUTH_API } from '../..';
import type { ApiSession } from '../../api/session/session.api.types';
import { FlidaApi } from '../../api/flida/flida.api';
import type { Flida } from '../../api/flida/flida.types';
import type { ISession } from './auth.types';

const PASSIVE_DISABLED_KEY = 'flida_passive_disabled';

export class AuthSessionService {
  private service = AUTH_API().session;
  private config = useRuntimeConfig().public.flida;

  async createFlidaClient(): Promise<FlidaApi> {
    const passiveDisabled = sessionStorage.getItem(PASSIVE_DISABLED_KEY) === '1';
    const fedcm = (this.config.fedcm as Flida.FedCMMode[]).filter((mode) => !(passiveDisabled && mode === 'passive'));
    return FlidaApi.create({
      clientId: this.config.clientId,
      redirectUri: this.config.redirectUri,
      mode: this.config.mode as Flida.AuthorizationMode,
      fedcm,
    });
  }

  async open(): Promise<ISession> {
    await this.service.open();
    return this.get();
  }

  async authenticate(idToken: string): Promise<ISession> {
    await this.service.authenticate(idToken);
    sessionStorage.removeItem(PASSIVE_DISABLED_KEY);
    return this.get();
  }

  async get(): Promise<ISession> {
    return this.toInterface(this.required(await this.service.get()));
  }

  async close(): Promise<void> {
    sessionStorage.setItem(PASSIVE_DISABLED_KEY, '1');
    await this.service.close();
  }

  createDemoIdToken(subject: string): string {
    const encode = (value: object) => btoa(JSON.stringify(value)).replace(/=+$/, '').replace(/\+/g, '-').replace(/\//g, '_');
    const claims = { iss: DEMO_ISSUER, sub: subject, exp: Math.floor(Date.now() / 1000) + 300 };
    return `${encode({ alg: 'none', typ: 'JWT' })}.${encode(claims)}.`;
  }

  flidaLogoutUrl(): string | null {
    const { clientId, postLogoutRedirectUri } = this.config;
    return clientId && postLogoutRedirectUri ? FlidaApi.logoutUrl(clientId, postLogoutRedirectUri) : null;
  }

  toInterface(data: ApiSession.Response): ISession {
    return {
      userId: data.userId || null,
      hasCoordinatorProfile: data.hasCoordinatorProfile,
      expiresAt: fromTimestamp(data.expiresAt),
    };
  }

  private required(session: ApiSession.Response | undefined): ApiSession.Response {
    if (!session) throw new ConnectError('Пустой ответ сервера');
    return session;
  }
}
