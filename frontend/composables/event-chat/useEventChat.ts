import { Code, ConnectError } from '@connectrpc/connect';
import type { IEventChat, IEventChatMessage } from '~/core/event-chat/services/event-chat/event-chat.types';

const HISTORY_PAGE = 50;

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

  let stopped = true;

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

  const loseAccess = (e: unknown) => {
    const { code } = ConnectError.from(e);
    if (code !== Code.PermissionDenied && code !== Code.NotFound) return;
    accessLost.value = true;
    messages.value = [];
  };

  const loadNewer = async () => {
    const after = lastSequence.value;
    let snapshotSequence: number | undefined;
    let offset = 0;
    while (true) {
      const history = await service.findMessages(eventId.value, snapshotSequence, { limit: HISTORY_PAGE, offset });
      snapshotSequence = history.snapshotSequence;
      snapshot.value = history.snapshotSequence;
      total.value = history.total;
      merge(history.list);
      offset += history.list.length;
      const oldest = Math.min(...history.list.map((item) => item.sequence));
      if (history.list.length < HISTORY_PAGE || oldest <= after || offset >= history.total) return;
    }
  };

  const sync = async () => {
    if (stopped || !chat.value || accessLost.value) return;
    await handleAsyncOperation(
      async () => {
        chat.value = await service.get(eventId.value);
        await loadNewer();
      },
      undefined,
      loseAccess,
    );
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
      loseAccess,
    );
    if (!ok) stopped = true;
  };

  const close = () => {
    stopped = true;
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

  useLiveRefresh(sync);

  onBeforeUnmount(close);

  return { chat, messages, loading, loadingOlder, sending, accessLost, hasOlder, isOwn, open, close, loadOlder, send };
};
