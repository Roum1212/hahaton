<script setup lang="ts">
const store = useCoordinatorsStore();
const { toast } = useApp();

const loading = ref(true);
const saving = ref(false);
const description = ref('');
const savingAvatar = ref(false);
const cooldownHours = ref(24);

onMounted(async () => {
  await handleAsyncOperation(
    async () => {
      const [me, settings] = await Promise.all([store.fetchMe(), store.fetchSettings()]);
      description.value = me.description;
      cooldownHours.value = Math.round(settings.cooldownSeconds / 3600);
    },
    () => {
      loading.value = false;
    },
  );
});

const saveAvatar = async (avatarPath: string) => {
  savingAvatar.value = true;
  const result = await handleAsyncOperation(
    () => store.updateMe({ avatarPath }),
    () => {
      savingAvatar.value = false;
    },
  );
  if (result) toast.success(avatarPath ? 'Фото обновлено' : 'Фото удалено');
};

const onSave = async () => {
  saving.value = true;
  const result = await handleAsyncOperation(
    () =>
      Promise.all([
        store.updateMe({ description: description.value.trim() }),
        store.updateSettings({ cooldownSeconds: Math.max(0, cooldownHours.value) * 3600 }),
      ]),
    () => {
      saving.value = false;
    },
  );
  if (result) toast.success('Профиль координатора сохранён');
};
</script>

<template>
  <Skeleton v-if="loading" height="320px" border-radius="20px" />

  <form v-else class="app-card" @submit.prevent="onSave">
    <span class="app-card__title">О вас</span>
    <div class="app-field">
      <span class="app-field__label flex items-center gap-2">
        Фото координатора
        <i
          v-tooltip.bottom="'Видят волонтёры. Не связано с фото вашего профиля волонтёра'"
          class="pi pi-question-circle text-surface-400" />
      </span>
      <AppImageUpload
        :image="store.me?.avatarUrl ?? null"
        :name="store.me?.displayName ?? ''"
        :disabled="savingAvatar"
        @uploaded="saveAvatar($event.path)"
        @remove="saveAvatar('')" />
    </div>
    <label class="app-field">
      <span class="app-field__label">Описание</span>
      <Textarea v-model="description" rows="4" auto-resize fluid placeholder="Чем вы занимаетесь и кого ищете" />
    </label>
    <label class="app-field">
      <span class="app-field__label flex items-center gap-2">
        Повторная заявка после отказа, часов
        <i
          v-tooltip.bottom="
            'Через сколько часов волонтёр сможет снова подать заявку на сотрудничество после отказа. Приглашение по ссылке действует сразу.'
          "
          class="pi pi-question-circle text-surface-400" />
      </span>
      <InputNumber v-model="cooldownHours" :min="0" :max="8760" show-buttons fluid />
    </label>
    <div class="app-form__actions">
      <Button type="submit" label="Сохранить" :loading="saving" />
    </div>
  </form>
</template>
