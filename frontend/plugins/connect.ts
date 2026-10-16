import { createClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { createAuthInterceptor, createErrorInterceptor, createRefreshInterceptor } from '~/core/app/rpc/interceptors';
import { accessToken } from '~/core/app/rpc/token';
import { AuthService } from '#genproto/auth/v1/auth_pb';

export default defineNuxtPlugin(() => {
  const baseUrl = useRuntimeConfig().public.rpcUrl;
  const { toast } = useApp();

  const fetchWithCredentials: typeof globalThis.fetch = (input, init) => globalThis.fetch(input, { ...init, credentials: 'include' });
  const errorInterceptor = createErrorInterceptor((error) => toast.error(getRpcErrorMessage(error)));
  const authInterceptor = createAuthInterceptor(accessToken.get);

  const publicTransport = createConnectTransport({
    baseUrl,
    fetch: fetchWithCredentials,
    interceptors: [errorInterceptor, authInterceptor],
  });

  const refreshClient = createClient(AuthService, createConnectTransport({ baseUrl, fetch: fetchWithCredentials }));
  let refreshing: Promise<string | null> | null = null;
  const refresh = () => {
    refreshing ??= refreshClient
      .refreshToken({})
      .then(({ accessToken: token }) => {
        accessToken.set(token);
        return token;
      })
      .catch(() => {
        accessToken.clear();
        return null;
      })
      .finally(() => {
        refreshing = null;
      });
    return refreshing;
  };

  const transport = createConnectTransport({
    baseUrl,
    fetch: fetchWithCredentials,
    interceptors: [errorInterceptor, createRefreshInterceptor(refresh), authInterceptor],
  });

  return {
    provide: {
      transport,
      publicTransport,
    },
  };
});
