<script setup lang="ts">
const eventsStore = useEventsStore();
const { visible } = useEventCreateDialog();
const { generateInitialData, validate, toDto } = useEventForm();
const { toast } = useApp();

const form = ref(generateInitialData());
const submitting = ref(false);

watch(visible, (value) => {
  if (value) form.value = generateInitialData();
});

const onSubmit = async () => {
  const error = validate(form.value);
  if (error) {
    toast.error(error);
    return;
  }
  submitting.value = true;
  const event = await handleAsyncOperation(
    () => eventsStore.create(toDto(form.value)),
    () => {
      submitting.value = false;
    },
  );
  if (!event) return;
  visible.value = false;
  toast.success('Мероприятие опубликовано');
  await navigateTo(`${ROUTES.coordinator.events}/${event.id}`);
};
</script>

<template>
  <Dialog
    v-model:visible="visible"
    modal
    header="Новое мероприятие"
    :draggable="false"
    :closable="!submitting"
    :style="{ width: '860px' }"
    :breakpoints="{ '767px': '100vw' }">
    <form class="app-form" @submit.prevent="onSubmit">
      <EventForm v-if="visible" v-model="form" :disabled="submitting" />
      <div class="app-form__actions">
        <Button label="Отмена" severity="secondary" text :disabled="submitting" @click="visible = false" />
        <Button type="submit" label="Опубликовать" :loading="submitting" />
      </div>
    </form>
  </Dialog>
</template>
