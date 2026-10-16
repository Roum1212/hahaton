<script setup lang="ts">
import type { IUpdateVolunteerDto } from '~/core/volunteer/services/volunteer/volunteer.types';

const volunteerStore = useVolunteerStore();
const { logout } = useAuthActions();
const { toast } = useApp();

const generateInitialData = (): IUpdateVolunteerDto => ({
  firstName: volunteerStore.profile?.firstName ?? '',
  lastName: volunteerStore.profile?.lastName ?? '',
  middleName: volunteerStore.profile?.middleName ?? '',
  dateOfBirth: volunteerStore.profile?.dateOfBirth ?? null,
});

const dto = ref<IUpdateVolunteerDto>(generateInitialData());
const submitting = ref(false);
const today = new Date();

watch(
  () => volunteerStore.profile,
  () => {
    dto.value = generateInitialData();
  },
);

const onSave = async () => {
  if (!dto.value.firstName || !dto.value.lastName) {
    toast.error('Укажите имя и фамилию');
    return;
  }
  submitting.value = true;
  const profile = await handleAsyncOperation(
    () => volunteerStore.update(dto.value),
    () => {
      submitting.value = false;
    },
  );
  if (profile) toast.success('Анкета сохранена');
};
</script>

<template>
  <form v-if="volunteerStore.profile" class="app-card" @submit.prevent="onSave">
    <span class="app-card__title">Личные данные</span>
    <div class="app-form-grid">
      <label class="app-field">
        <span class="app-field__label">Фамилия</span>
        <InputText v-model.trim="dto.lastName" :disabled="submitting" fluid />
      </label>
      <label class="app-field">
        <span class="app-field__label">Имя</span>
        <InputText v-model.trim="dto.firstName" :disabled="submitting" fluid />
      </label>
      <label class="app-field">
        <span class="app-field__label">Отчество</span>
        <InputText v-model.trim="dto.middleName" :disabled="submitting" fluid />
      </label>
      <label class="app-field">
        <span class="app-field__label">Дата рождения</span>
        <DatePicker v-model="dto.dateOfBirth" :max-date="today" date-format="dd.mm.yy" :disabled="submitting" show-icon fluid />
      </label>
    </div>

    <span class="app-card__title mt-2">Контакты</span>
    <div class="app-form-grid">
      <label class="app-field">
        <span class="app-field__label">Email</span>
        <InputText :model-value="volunteerStore.profile.email" disabled fluid />
      </label>
      <label class="app-field">
        <span class="app-field__label">Телефон</span>
        <InputText :model-value="volunteerStore.profile.phone" disabled fluid />
      </label>
    </div>

    <div class="app-form__actions">
      <Button label="Выйти из аккаунта" icon="pi pi-sign-out" severity="secondary" text @click="logout" />
      <Button type="submit" label="Сохранить" :loading="submitting" />
    </div>
  </form>
</template>
