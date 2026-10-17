<script setup lang="ts">
const props = defineProps<{
  eventId: string;
}>();

const eventId = computed(() => props.eventId);
const { chat, messages, loading, loadingOlder, sending, accessLost, hasOlder, isOwn, open, loadOlder, send } =
  useEventChat(eventId);

const text = ref('');
const feed = ref<HTMLElement | null>(null);

const scrollToBottom = () => nextTick(() => feed.value?.scrollTo({ top: feed.value.scrollHeight }));

const onSend = async () => {
  if (await send(text.value)) {
    text.value = '';
    scrollToBottom();
  }
};

const onKeydown = (event: KeyboardEvent) => {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault();
    onSend();
  }
};

watch(
  () => messages.value.length,
  (length, previous) => {
    if (!previous || length - previous === 1) scrollToBottom();
  },
);

onMounted(open);
</script>

<template>
  <div v-if="accessLost" class="app-card">
    <AppEmpty title="Доступ к чату закрыт" />
  </div>

  <div v-else class="app-chat">
    <div ref="feed" class="app-chat__feed app-scroll">
      <div v-if="hasOlder" class="app-more">
        <Button
          label="Предыдущие сообщения"
          severity="secondary"
          text
          size="small"
          :loading="loadingOlder"
          @click="loadOlder" />
      </div>
      <div v-if="loading" class="flex flex-col gap-3">
        <Skeleton v-for="n in 4" :key="n" height="48px" width="60%" />
      </div>
      <div v-else-if="!messages.length" class="app-chat__empty">
        {{ chat?.readOnly ? 'Сообщений пока нет' : 'Напишите первое сообщение участникам' }}
      </div>
      <div
        v-for="message in messages"
        :key="message.id"
        class="app-chat__message"
        :class="{ 'app-chat__message--own': isOwn(message) }">
        <div class="app-chat__author">
          {{ message.authorDisplayName }}
          <span v-if="message.author === 'coordinator'" class="app-chip app-chip--violet app-chat__badge"
            >Координатор</span
          >
        </div>
        <div class="app-chat__bubble">{{ message.text }}</div>
        <span class="app-chat__time">{{ formatDateTime(message.createdAt) }}</span>
      </div>
    </div>

    <Message v-if="chat?.readOnly" severity="secondary" size="small">Чат доступен только для чтения</Message>
    <form v-else class="app-chat__composer" @submit.prevent="onSend">
      <Textarea
        v-model="text"
        rows="1"
        auto-resize
        maxlength="2000"
        placeholder="Сообщение"
        :disabled="loading"
        fluid
        @keydown="onKeydown" />
      <Button type="submit" icon="pi pi-send" :loading="sending" :disabled="!text.trim()" aria-label="Отправить" />
    </form>
  </div>
</template>
