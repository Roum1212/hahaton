import type { ConnectRouter } from '@connectrpc/connect';
import { UpdateKind, UpdatesService } from '../../genproto/updates/v1/updates_service_pb';
import { requireSession, requireUser } from '../auth';
import { UpdatesSubscriber, updatesHub } from '../updates-hub';

const HEARTBEAT_MS = 10_000;

export const registerUpdatesService = (router: ConnectRouter) =>
  router.service(UpdatesService, {
    async *subscribeUpdates(_, ctx) {
      const session = requireSession(ctx);
      const user = requireUser(ctx);
      const subscriber = new UpdatesSubscriber(user.id, session.id);
      updatesHub.add(subscriber);
      try {
        yield { kind: UpdateKind.READY };
        while (true) {
          const result = await subscriber.next(ctx.signal, HEARTBEAT_MS);
          if (result === 'aborted') return;
          if (result === 'heartbeat') requireUser(ctx);
          yield { kind: result === 'changed' ? UpdateKind.CHANGED : UpdateKind.HEARTBEAT };
        }
      } finally {
        updatesHub.remove(subscriber);
      }
    },
  });
