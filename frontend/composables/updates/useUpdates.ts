import { Code, ConnectError } from '@connectrpc/connect';
import { UpdateKind } from '~/core/updates/services/updates/updates.types';

const RECONNECT_DELAYS_MS = [1000, 2000, 5000, 10000];
const NOTIFY_DELAY_MS = 300;

let controller: AbortController | null = null;
let notifyTimer: ReturnType<typeof setTimeout> | null = null;

const sleep = (ms: number, signal: AbortSignal) =>
  new Promise<void>((resolve) => {
    if (signal.aborted) return resolve();
    const finish = () => {
      clearTimeout(timer);
      signal.removeEventListener('abort', finish);
      resolve();
    };
    const timer = setTimeout(finish, ms);
    signal.addEventListener('abort', finish, { once: true });
  });

export const useUpdatesConnection = () => {
  const store = useUpdatesStore();
  const authStore = useAuthStore();

  const scheduleNotify = () => {
    if (notifyTimer) return;
    notifyTimer = setTimeout(() => {
      notifyTimer = null;
      store.notify();
    }, NOTIFY_DELAY_MS);
  };

  const run = async (current: AbortController) => {
    let attempt = 0;
    while (!current.signal.aborted) {
      try {
        for await (const kind of store.service.subscribe(current.signal)) {
          if (current.signal.aborted) break;
          attempt = 0;
          store.connected = true;
          if (kind === UpdateKind.CHANGED || kind === UpdateKind.READY) scheduleNotify();
        }
      } catch (e) {
        const error = ConnectError.from(e);
        if (current.signal.aborted) break;
        if (error.code === Code.Unauthenticated) {
          stop();
          authStore.expire();
          break;
        }
      }
      if (current.signal.aborted) break;
      store.connected = false;
      await sleep(RECONNECT_DELAYS_MS[Math.min(attempt, RECONNECT_DELAYS_MS.length - 1)]!, current.signal);
      attempt += 1;
    }
    if (controller === current) stop();
  };

  const start = () => {
    if (controller) return;
    const current = new AbortController();
    controller = current;
    void run(current);
  };

  const stop = () => {
    controller?.abort();
    controller = null;
    store.connected = false;
    if (notifyTimer) clearTimeout(notifyTimer);
    notifyTimer = null;
  };

  return { start, stop };
};

export const useLiveRefresh = (refresh: () => unknown) => {
  const store = useUpdatesStore();
  watch(
    () => store.revision,
    () => refresh(),
  );
};
