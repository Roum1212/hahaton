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
const savingAvatar = ref(false);

const saveAvatar = async (avatarPath: string) => {
  savingAvatar.value = true;
  const profile = await handleAsyncOperation(
    () => volunteerStore.updateAvatar(avatarPath),
    () => {
      savingAvatar.value = false;
    },
  );
  if (profile) toast.success(avatarPath ? 'Фото обновлено' : 'Фото удалено');
};
const today = new Date();
const formElement = ref<HTMLFormElement | null>(null);
const { errors, errorId, validate, reset } = useFormValidation(formElement, () => {
  const errors: Partial<Record<'firstName' | 'lastName' | 'dateOfBirth', string>> = {};
  if (!dto.value.firstName.trim()) errors.firstName = 'Укажите имя';
  if (!dto.value.lastName.trim()) errors.lastName = 'Укажите фамилию';
  const date = dto.value.dateOfBirth;
  if (date && (!Number.isFinite(date.getTime()) || date > today))
    errors.dateOfBirth = 'Дата рождения не может быть в будущем';
  return errors;
});

watch(
  () => volunteerStore.profile,
  () => {
    dto.value = generateInitialData();
    reset();
  },
);

const onSave = async () => {
  if (!(await validate())) return;
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
  <form v-if="volunteerStore.profile" ref="formElement" class="app-card" @submit.prevent="onSave">
    <span class="app-card__title">Фото профиля</span>
    <AppImageUpload
      :image="volunteerStore.profile.avatarUrl"
      :name="volunteerStore.profile.fullName"
      :disabled="savingAvatar"
      @uploaded="saveAvatar($event.path)"
      @remove="saveAvatar('')" />

    <span class="app-card__title mt-2">Личные данные</span>
    <div class="app-form-grid">
      <label class="app-field" :data-invalid="!!errors.lastName">
        <span class="app-field__label">Фамилия</span>
        <InputText
          v-model.trim="dto.lastName"
          :invalid="!!errors.lastName"
          :aria-describedby="errors.lastName ? errorId('lastName') : undefined"
          :disabled="submitting"
          fluid />
        <small v-if="errors.lastName" :id="errorId('lastName')" class="app-field__error">{{ errors.lastName }}</small>
      </label>
      <label class="app-field" :data-invalid="!!errors.firstName">
        <span class="app-field__label">Имя</span>
        <InputText
          v-model.trim="dto.firstName"
          :invalid="!!errors.firstName"
          :aria-describedby="errors.firstName ? errorId('firstName') : undefined"
          :disabled="submitting"
          fluid />
        <small v-if="errors.firstName" :id="errorId('firstName')" class="app-field__error">{{
          errors.firstName
        }}</small>
      </label>
      <label class="app-field">
        <span class="app-field__label">Отчество</span>
        <InputText v-model.trim="dto.middleName" :disabled="submitting" fluid />
      </label>
      <label class="app-field" :data-invalid="!!errors.dateOfBirth">
        <span class="app-field__label">Дата рождения</span>
        <DatePicker
          v-model="dto.dateOfBirth"
          :invalid="!!errors.dateOfBirth"
          :aria-describedby="errors.dateOfBirth ? errorId('dateOfBirth') : undefined"
          :max-date="today"
          date-format="dd.mm.yy"
          :disabled="submitting"
          show-icon
          fluid
          :pt="{
            pcInputText: { root: { 'aria-describedby': errors.dateOfBirth ? errorId('dateOfBirth') : undefined } },
          }" />
        <small v-if="errors.dateOfBirth" :id="errorId('dateOfBirth')" class="app-field__error">{{
          errors.dateOfBirth
        }}</small>
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
