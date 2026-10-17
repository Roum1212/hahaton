import { Code, ConnectError } from '@connectrpc/connect';
import { UpdateKind } from '~/core/updates/services/updates/updates.types';

const RECONNECT_DELAYS_MS = [1000, 2000, 5000, 10000];
const NOTIFY_DELAY_MS = 300;

let running = false;
let controller: AbortController | null = null;
let notifyTimer: ReturnType<typeof setTimeout> | null = null;

const sleep = (ms: number, signal: AbortSignal) =>
  new Promise<void>((resolve) => {
    const timer = setTimeout(resolve, ms);
    signal.addEventListener(
      'abort',
      () => {
        clearTimeout(timer);
        resolve();
      },
      { once: true },
    );
  });

export const useUpdatesConnection = () => {
  const store = useUpdatesStore();
  const authStore = useAuthStore();

  const scheduleNotify = () => {
    if (notifyTimer) clearTimeout(notifyTimer);
    notifyTimer = setTimeout(() => {
      notifyTimer = null;
      store.notify();
    }, NOTIFY_DELAY_MS);
  };

  const run = async () => {
    let attempt = 0;
    let synced = false;
    while (running) {
      const current = new AbortController();
      controller = current;
      try {
        for await (const kind of store.service.subscribe(current.signal)) {
          attempt = 0;
          store.connected = true;
          if (kind === UpdateKind.CHANGED || (kind === UpdateKind.READY && synced)) scheduleNotify();
          if (kind === UpdateKind.READY) synced = true;
        }
      } catch (e) {
        const error = ConnectError.from(e);
        if (!running || error.code === Code.Canceled) break;
        if (error.code === Code.Unauthenticated) {
          authStore.expire();
          break;
        }
      }
      store.connected = false;
      if (!running) break;
      await sleep(RECONNECT_DELAYS_MS[Math.min(attempt, RECONNECT_DELAYS_MS.length - 1)]!, current.signal);
      attempt += 1;
    }
    store.connected = false;
  };

  const start = () => {
    if (running) return;
    running = true;
    run();
  };

  const stop = () => {
    running = false;
    controller?.abort();
    controller = null;
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
