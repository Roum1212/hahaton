<script setup lang="ts">
definePageMeta({
  title: 'Редактирование',
  permissions: [ROLES.coordinator],
});

const route = useRoute();
const eventsStore = useEventsStore();
const { generateInitialData, getErrors, toDto } = useEventForm();
const { toast } = useApp();

const eventId = computed(() => String(route.params.id));
const loading = ref(true);
const editable = ref(false);
const form = ref(generateInitialData());
const submitting = ref(false);
const formElement = ref<HTMLFormElement | null>(null);
const { errors, validate } = useFormValidation(formElement, () => getErrors(form.value));

onMounted(async () => {
  const event = await handleAsyncOperation(
    () => eventsStore.fetchOne(eventId.value, true),
    () => {
      loading.value = false;
    },
  );
  if (!event) return;
  editable.value = canEditEvent(event);
  form.value = generateInitialData(event);
});

const onSubmit = async () => {
  if (!(await validate())) return;
  submitting.value = true;
  const event = await handleAsyncOperation(
    () => eventsStore.update(eventId.value, toDto(form.value)),
    () => {
      submitting.value = false;
    },
  );
  if (!event) return;
  toast.success('Изменения сохранены');
  await navigateTo(`${ROUTES.coordinator.events}/${event.id}`);
};
</script>

<template>
  <AppPage :loading="loading" :back="`${ROUTES.coordinator.events}/${eventId}`" back-label="К мероприятию">
    <template #title>Редактирование мероприятия</template>

    <Skeleton v-if="loading" height="420px" border-radius="20px" />

    <form v-else ref="formElement" class="app-card max-w-[860px]" @submit.prevent="onSubmit">
      <Message v-if="!editable" severity="warn" size="small">
        Изменить мероприятие можно только до начала и пока на него нет активных заявок
      </Message>
      <EventForm v-model="form" :errors="errors" :disabled="submitting || !editable" />
      <div class="app-form__actions">
        <NuxtLink :to="`${ROUTES.coordinator.events}/${eventId}`">
          <Button label="Отмена" severity="secondary" text />
        </NuxtLink>
        <Button type="submit" label="Сохранить" :loading="submitting" :disabled="!editable" />
      </div>
    </form>
  </AppPage>
</template>
