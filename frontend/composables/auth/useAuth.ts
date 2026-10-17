import { ConnectError } from '@connectrpc/connect';
import type { FlidaApi } from '~/core/auth/api/flida/flida.api';
import type { IFlidaTokens } from '~/core/auth/services/auth/auth.types';
import type { IDemoAccount } from '~/config/demo.config';

const FLIDA_ERROR_MESSAGES: Partial<Record<string, string>> = {
  invalid_config: 'Флида ID не настроен: проверьте Client ID и Redirect URI',
  invalid_state: 'Сессия входа устарела, попробуйте ещё раз',
  authorization_in_progress: 'Вход уже выполняется',
  popup_blocked: 'Браузер заблокировал окно входа. Разрешите всплывающие окна для сайта',
  timeout: 'Время ожидания входа истекло',
  crypto_unavailable: 'Браузер не поддерживает безопасный вход',
};

const SILENT_FLIDA_ERRORS = ['popup_closed', 'fedcm_unavailable'];

export const useAuthActions = () => {
  const authStore = useAuthStore();
  const loader = useLoader();

  const redirectAfterLogin = async (fallback?: string) => {
    const redirectCookie = useCookie<string | null>('redirectAfterLogin');
    const saved = takeLoginRedirect() ?? normalizeLoginRedirect(redirectCookie.value);
    const allowed = !!saved && (modeFromPath(saved) !== 'coordinator' || authStore.hasCoordinatorProfile);
    const path = allowed ? saved : (fallback ?? useCabinetMode().homeFor());
    redirectCookie.value = null;
    return navigateTo(path, { replace: true });
  };

  const logout = async () => {
    loader.show();
    try {
      await authStore.logout();
    } finally {
      loader.hide();
    }
    const flidaLogoutUrl = authStore.service.flidaLogoutUrl();
    if (flidaLogoutUrl) {
      window.location.assign(flidaLogoutUrl);
      return;
    }
    await navigateTo(ROUTES.login);
  };

  return { redirectAfterLogin, logout };
};

export const useDemoAuth = () => {
  const authStore = useAuthStore();
  const { toast, loader } = useApp();
  const { redirectAfterLogin } = useAuthActions();
  const enabled = !!useRuntimeConfig().public.demoLogin;

  const signIn = async (account: IDemoAccount) => {
    loader.show();
    try {
      const ok = await handleAsyncOperation(async () => {
        await authStore.loginWithFlida(authStore.service.createDemoIdToken(account.subject));
        return true;
      });
      if (!ok) return;
      toast.success(`Вы вошли как ${account.name}`);
      await redirectAfterLogin(useCabinetMode().homeFor(account.mode));
    } finally {
      loader.hide();
    }
  };

  return { enabled, signIn };
};

export const useFlidaAuth = () => {
  const authStore = useAuthStore();
  const { toast, loader } = useApp();
  const { redirectAfterLogin } = useAuthActions();

  const client = shallowRef<FlidaApi | null>(null);
  const initializing = ref(true);
  const authorizing = ref(false);
  const configured = !!useRuntimeConfig().public.flida.clientId;

  const handleError = (error: unknown) => {
    if (error instanceof ConnectError) return;
    if (error instanceof DOMException && error.name === 'AbortError') return;
    if (client.value?.isFlidaError(error)) {
      if (SILENT_FLIDA_ERRORS.includes(error.code)) return;
      toast.error(FLIDA_ERROR_MESSAGES[error.code] ?? 'Не удалось войти через Флида ID');
      return;
    }
    console.error(error);
    toast.error(client.value ? 'Не удалось войти через Флида ID' : 'Не удалось загрузить Флида ID');
  };

  const complete = async (tokens: IFlidaTokens | undefined): Promise<boolean> => {
    if (!tokens) return false;
    if (!tokens.id_token) {
      toast.error('Флида ID не передал данные для входа');
      return false;
    }
    loader.show();
    try {
      await authStore.loginWithFlida(tokens.id_token);
      toast.success('Вы успешно вошли в систему');
      await redirectAfterLogin();
      return true;
    } finally {
      loader.hide();
    }
  };

  const init = async (options: { onReady?: () => void; silent?: boolean } = {}): Promise<boolean> => {
    initializing.value = true;
    try {
      client.value = await authStore.service.createFlidaClient();
      options.onReady?.();
    } catch (error) {
      handleError(error);
      return false;
    } finally {
      initializing.value = false;
    }
    return client.value.autoAuthorization.then(complete).catch((error) => {
      if (!options.silent) handleError(error);
      return false;
    });
  };

  const signIn = async () => {
    if (!client.value || authorizing.value) return;
    authorizing.value = true;
    try {
      await complete(await client.value.authorize());
    } catch (error) {
      handleError(error);
    } finally {
      authorizing.value = false;
    }
  };

  return { configured, initializing, authorizing, init, signIn };
};
