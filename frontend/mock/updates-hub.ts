import { Code, ConnectError } from '@connectrpc/connect';

export type UpdatesWait = 'changed' | 'heartbeat' | 'aborted';

export class UpdatesSubscriber {
  private pending = false;
  private closed: ConnectError | null = null;
  private waiter: (() => void) | null = null;

  constructor(
    readonly userId: string,
    readonly sessionId: string,
  ) {}

  private wake() {
    this.waiter?.();
    this.waiter = null;
  }

  notify() {
    this.pending = true;
    this.wake();
  }

  close(error: ConnectError) {
    this.closed = error;
    this.wake();
  }

  async next(signal: AbortSignal, heartbeatMs: number): Promise<UpdatesWait> {
    if (!this.pending && !this.closed && !signal.aborted) {
      await new Promise<void>((resolve) => {
        const timer = setTimeout(resolve, heartbeatMs);
        const done = () => {
          clearTimeout(timer);
          resolve();
        };
        this.waiter = done;
        signal.addEventListener('abort', done, { once: true });
      });
    }
    if (this.closed) throw this.closed;
    if (signal.aborted) return 'aborted';
    if (!this.pending) return 'heartbeat';
    this.pending = false;
    return 'changed';
  }
}

const subscribers = new Set<UpdatesSubscriber>();

export const updatesHub = {
  add(subscriber: UpdatesSubscriber) {
    subscribers.add(subscriber);
  },
  remove(subscriber: UpdatesSubscriber) {
    subscribers.delete(subscriber);
  },
  notifyAll() {
    subscribers.forEach((item) => item.notify());
  },
  closeSession(sessionId: string) {
    subscribers.forEach((item) => item.sessionId === sessionId && item.close(new ConnectError('Сессия завершена', Code.Unauthenticated)));
  },
};
