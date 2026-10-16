import { Code, ConnectError } from '@connectrpc/connect';
import type { IEventChat, IEventChatMessage } from '~/core/event-chat/services/event-chat/event-chat.types';

const HISTORY_PAGE = 50;
const RECONNECT_DELAY_MS = 2000;

export const useEventChat = (eventId: Ref<string>) => {
  const service = useEventChatsStore().service;
  const authStore = useAuthStore();

  const chat = ref<IEventChat | null>(null);
  const messages = ref<IEventChatMessage[]>([]);
  const loading = ref(true);
  const loadingOlder = ref(false);
  const sending = ref(false);
  const accessLost = ref(false);
  const total = ref(0);
  const snapshot = ref<number | undefined>(undefined);

  let controller: AbortController | null = null;
  let stopped = false;

  const isOwn = (message: IEventChatMessage) => message.authorId === authStore.session?.userId;
  const hasOlder = computed(() => messages.value.length < total.value);
  const lastSequence = computed(() => messages.value.at(-1)?.sequence ?? snapshot.value ?? 0);

  const merge = (items: IEventChatMessage[]) => {
    const byId = new Map(messages.value.map((item) => [item.id, item]));
    items.forEach((item) => byId.set(item.id, item));
    messages.value = [...byId.values()].sort((a, b) => a.sequence - b.sequence);
  };

  const loadHistory = async () => {
    const history = await service.findMessages(eventId.value, undefined, { limit: HISTORY_PAGE, offset: 0 });
    snapshot.value = history.snapshotSequence;
    total.value = history.total;
    merge(history.list);
  };

  const loadOlder = async () => {
    if (!hasOlder.value || loadingOlder.value) return;
    loadingOlder.value = true;
    const loadedFromSnapshot = messages.value.filter((item) => item.sequence <= (snapshot.value ?? 0)).length;
    const history = await handleAsyncOperation(
      () => service.findMessages(eventId.value, snapshot.value, { limit: HISTORY_PAGE, offset: loadedFromSnapshot }),
      () => {
        loadingOlder.value = false;
      },
    );
    if (history) merge(history.list);
  };

  const subscribe = async () => {
    while (!stopped && chat.value && !chat.value.readOnly) {
      controller = new AbortController();
      try {
        for await (const message of service.subscribe(eventId.value, lastSequence.value, controller.signal)) merge([message]);
        chat.value = await service.get(eventId.value);
        return;
      } catch (e) {
        const error = ConnectError.from(e);
        if (stopped || error.code === Code.Canceled) return;
        if (error.code === Code.PermissionDenied || error.code === Code.NotFound) {
          accessLost.value = true;
          messages.value = [];
          return;
        }
        if (error.code === Code.Unauthenticated) {
          authStore.expire();
          return;
        }
        await new Promise((resolve) => setTimeout(resolve, RECONNECT_DELAY_MS));
      }
    }
  };

  const open = async () => {
    stopped = false;
    loading.value = true;
    accessLost.value = false;
    messages.value = [];
    const ok = await handleAsyncOperation(
      async () => {
        chat.value = await service.get(eventId.value);
        await loadHistory();
        return true;
      },
      () => {
        loading.value = false;
      },
      (e) => {
        if (ConnectError.from(e).code === Code.PermissionDenied) accessLost.value = true;
        return false;
      },
    );
    if (ok) subscribe();
  };

  const close = () => {
    stopped = true;
    controller?.abort();
  };

  const send = async (text: string) => {
    const value = text.trim();
    if (!value || sending.value) return false;
    sending.value = true;
    const requestId = createRequestId();
    const message = await handleAsyncOperation(
      () => service.send(eventId.value, value, requestId),
      () => {
        sending.value = false;
      },
    );
    if (!message) return false;
    merge([message]);
    return true;
  };

  onBeforeUnmount(close);

  return { chat, messages, loading, loadingOlder, sending, accessLost, hasOlder, isOwn, open, close, loadOlder, send };
};
