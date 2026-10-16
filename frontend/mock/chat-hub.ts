import { Code, ConnectError } from '@connectrpc/connect';
import type { ChatMessageRec } from './db';

type Item = { kind: 'message'; message: ChatMessageRec } | { kind: 'end' } | { kind: 'error'; error: ConnectError };

export class ChatSubscriber {
  private items: Item[] = [];
  private waiter: (() => void) | null = null;

  constructor(
    readonly eventId: string,
    readonly userId: string,
    readonly role: 'coordinator' | 'volunteer',
    readonly sessionId: string,
  ) {}

  private push(item: Item) {
    this.items.push(item);
    this.waiter?.();
    this.waiter = null;
  }

  deliver(message: ChatMessageRec) {
    if (this.items.length > 500) {
      this.fail(new ConnectError('Клиент не успевает получать сообщения', Code.ResourceExhausted));
      return;
    }
    this.push({ kind: 'message', message });
  }

  end() {
    this.push({ kind: 'end' });
  }

  fail(error: ConnectError) {
    this.push({ kind: 'error', error });
  }

  async next(signal: AbortSignal): Promise<Item> {
    while (!this.items.length) {
      if (signal.aborted) return { kind: 'end' };
      await new Promise<void>((resolve) => {
        this.waiter = resolve;
        signal.addEventListener('abort', () => resolve(), { once: true });
      });
    }
    return this.items.shift() as Item;
  }
}

const subscribers = new Set<ChatSubscriber>();

export const chatHub = {
  add(subscriber: ChatSubscriber) {
    subscribers.add(subscriber);
  },
  remove(subscriber: ChatSubscriber) {
    subscribers.delete(subscriber);
  },
  publish(message: ChatMessageRec) {
    subscribers.forEach((item) => item.eventId === message.eventId && item.deliver(message));
  },
  revokeVolunteer(eventId: string, volunteerId: string) {
    subscribers.forEach((item) => {
      if (item.eventId === eventId && item.role === 'volunteer' && item.userId === volunteerId) {
        item.fail(new ConnectError('Доступ к чату закрыт', Code.PermissionDenied));
      }
    });
  },
  completeEvent(eventId: string) {
    subscribers.forEach((item) => item.eventId === eventId && item.end());
  },
  revokeVolunteersOfEvent(eventId: string) {
    subscribers.forEach((item) => {
      if (item.eventId === eventId && item.role === 'volunteer') item.fail(new ConnectError('Доступ к чату закрыт', Code.PermissionDenied));
    });
  },
  closeSession(sessionId: string) {
    subscribers.forEach((item) => item.sessionId === sessionId && item.fail(new ConnectError('Сессия завершена', Code.Unauthenticated)));
  },
};
