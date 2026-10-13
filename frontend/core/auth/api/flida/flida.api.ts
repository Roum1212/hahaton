import type { Flida } from './flida.types';

const FLIDA_SDK_URL = 'https://id.flida.ru/sdk/v1/sdk.es.js';
const FLIDA_LOGOUT_URL = 'https://id.flida.ru/logout';

let modulePromise: Promise<Flida.Module> | null = null;

const loadModule = (): Promise<Flida.Module> => {
  modulePromise ??= import(/* @vite-ignore */ FLIDA_SDK_URL).catch((e) => {
    modulePromise = null;
    throw e;
  }) as Promise<Flida.Module>;
  return modulePromise;
};

export class FlidaApi {
  private constructor(
    private readonly sdk: Flida.SDK,
    private readonly module: Flida.Module,
  ) {}

  static async create(config: Flida.SDKConfig): Promise<FlidaApi> {
    const module = await loadModule();
    return new FlidaApi(new module.FlidaSDK(config), module);
  }

  get autoAuthorization(): Promise<Flida.TokenResponse | undefined> {
    return this.sdk.autoAuthorization;
  }

  async authorize(signal?: AbortSignal): Promise<Flida.TokenResponse | undefined> {
    return this.sdk.authorize({ signal });
  }

  async getUser(accessToken: string, signal?: AbortSignal): Promise<Flida.UserInfoResponse> {
    return this.sdk.getUser(accessToken, signal);
  }

  isFlidaError(error: unknown): error is Flida.Error {
    return error instanceof this.module.FlidaError;
  }

  static logoutUrl(clientId: string, postLogoutRedirectUri: string, state?: string): string {
    const params = new URLSearchParams({ client_id: clientId, post_logout_redirect_uri: postLogoutRedirectUri });
    if (state) params.set('state', state);
    return `${FLIDA_LOGOUT_URL}?${params}`;
  }
}
