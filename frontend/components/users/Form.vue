<script setup lang="ts">
import type { ICreateUserDto, IUser } from '~/core/users/services/users/users.types';

defineExpose<{ save: () => Promise<boolean> }>({
  save: onSave,
});

const props = withDefaults(
  defineProps<{
    user?: IUser | undefined;
    loading?: boolean;
  }>(),
  {
    user: undefined,
    loading: false,
  },
);

const { toast } = useApp();
const { save } = useUsersActions();
const { submitting, generateInitialData, roleOptions } = useUsersForm();

const dto = ref<ICreateUserDto>(generateInitialData(props.user));
const isEdit = computed(() => !!props.user?.id);

watch(
  () => props.user,
  (user) => {
    dto.value = generateInitialData(user);
  },
);

async function onSave(): Promise<boolean> {
  const { isValid, message } = validateUserDto(dto.value);
  if (!isValid) {
    if (message) toast.error(message);
    return false;
  }
  submitting.value = true;
  try {
    return !!(await save(props.user?.id, dto.value));
  } finally {
    submitting.value = false;
  }
}
</script>

<template>
  <form class="flex flex-col gap-4" @submit.prevent="onSave">
    <template v-if="loading">
      <Skeleton v-for="n in 4" :key="n" height="2.5rem" />
    </template>
    <template v-else>
      <label class="app-field">
        <span class="app-field__label">Email</span>
        <InputText v-model.trim="dto.email" placeholder="user@example.com" :disabled="submitting" fluid />
      </label>
      <div class="app-form-grid">
        <label class="app-field">
          <span class="app-field__label">Имя</span>
          <InputText v-model.trim="dto.firstName" placeholder="Введите имя" :disabled="submitting" fluid />
        </label>
        <label class="app-field">
          <span class="app-field__label">Фамилия</span>
          <InputText v-model.trim="dto.lastName" placeholder="Введите фамилию" :disabled="submitting" fluid />
        </label>
      </div>
      <label class="app-field">
        <span class="app-field__label">Роль</span>
        <Select v-model="dto.role" :options="roleOptions" option-label="label" option-value="value" :disabled="submitting" fluid />
      </label>
      <label v-if="isEdit" class="flex items-center gap-3">
        <ToggleSwitch v-model="dto.active" :disabled="submitting" />
        <span>Активен</span>
      </label>
    </template>
  </form>
</template>
