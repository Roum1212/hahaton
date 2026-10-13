import { ConnectError } from '@connectrpc/connect';
import { USERS_SERVICES } from '~/core/users';
import { AUTH_API } from '../..';
import { FlidaApi } from '../../api/flida/flida.api';
import type { Flida } from '../../api/flida/flida.types';
import type { IAuthUser } from './auth.types';

const PASSIVE_DISABLED_KEY = 'flida_passive_disabled';

export class AuthSessionService {
  private service = AUTH_API().auth;
  private users = USERS_SERVICES().UsersService;
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

  async loginWithFlida(flidaAccessToken: string): Promise<string> {
    const token = await this.service.loginWithFlida({ flidaAccessToken });
    sessionStorage.removeItem(PASSIVE_DISABLED_KEY);
    return token;
  }

  async logout(): Promise<void> {
    sessionStorage.setItem(PASSIVE_DISABLED_KEY, '1');
    await this.service.logout();
  }

  flidaLogoutUrl(): string | null {
    const { clientId, postLogoutRedirectUri } = this.config;
    return clientId && postLogoutRedirectUri ? FlidaApi.logoutUrl(clientId, postLogoutRedirectUri) : null;
  }

  async me(): Promise<IAuthUser> {
    const user = await this.service.me();
    if (!user) throw new ConnectError('Пустой ответ сервера');
    return this.users.toInterface(user);
  }
}
