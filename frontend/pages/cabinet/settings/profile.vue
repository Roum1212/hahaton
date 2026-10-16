<script setup lang="ts">
import { USER_ROLE_LABELS, type ICreateUserDto } from '~/core/users/services/users/users.types';

definePageMeta({
  title: 'Профиль',
});

const authStore = useAuthStore();
const { save } = useUsersActions();
const { logout } = useAuthActions();
const { generateInitialData } = useUsersForm();
const { toast } = useApp();

const dto = ref<ICreateUserDto>(generateInitialData(authStore.user ?? undefined));
const submitting = ref(false);

const onSave = async () => {
  if (!authStore.user) return;
  const { isValid, message } = validateUserDto(dto.value);
  if (!isValid) {
    if (message) toast.error(message);
    return;
  }
  submitting.value = true;
  await save(authStore.user.id, dto.value);
  submitting.value = false;
};
</script>

<template>
  <AppPage>
    <template #title>Профиль</template>
    <template #subtitle>Данные вашей учётной записи</template>

    <div v-if="authStore.user" class="app-card max-w-[720px]">
      <div class="flex items-center gap-3">
        <span class="app-user__avatar">{{ authStore.user.fullName.slice(0, 1).toUpperCase() }}</span>
        <div class="flex flex-col">
          <span class="app-card__title">{{ authStore.user.fullName }}</span>
          <span class="app-card__hint">{{ USER_ROLE_LABELS[authStore.user.role] }}</span>
        </div>
      </div>

      <form class="flex flex-col gap-4" @submit.prevent="onSave">
        <div class="app-form-grid">
          <label class="app-field">
            <span class="app-field__label">Имя</span>
            <InputText v-model.trim="dto.firstName" :disabled="submitting" fluid />
          </label>
          <label class="app-field">
            <span class="app-field__label">Фамилия</span>
            <InputText v-model.trim="dto.lastName" :disabled="submitting" fluid />
          </label>
        </div>
        <label class="app-field">
          <span class="app-field__label">Email</span>
          <InputText v-model.trim="dto.email" :disabled="submitting" fluid />
        </label>
        <div class="flex justify-end">
          <Button type="submit" label="Сохранить" icon="pi pi-check" :loading="submitting" />
        </div>
      </form>
    </div>

    <div class="app-logout-mobile max-w-[720px]">
      <Button label="Выйти" icon="pi pi-sign-out" severity="danger" fluid @click="logout" />
    </div>
  </AppPage>
</template>
