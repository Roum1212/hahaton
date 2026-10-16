import { createConnectTransport } from '@connectrpc/connect-web';
import { createErrorInterceptor } from '~/core/app/rpc/interceptors';

export default defineNuxtPlugin(() => {
  const baseUrl = useRuntimeConfig().public.rpcUrl;
  const { toast } = useApp();

  const fetchWithCredentials: typeof globalThis.fetch = (input, init) => globalThis.fetch(input, { ...init, credentials: 'include' });

  const transport = createConnectTransport({
    baseUrl,
    fetch: fetchWithCredentials,
    interceptors: [
      createErrorInterceptor(
        (error, method) => toast.error(getRpcErrorMessage(error, RPC_METHOD_ERROR_MESSAGES[method])),
        () => useAuthStore().expire(),
      ),
    ],
  });

  return {
    provide: {
      transport,
    },
  };
});
