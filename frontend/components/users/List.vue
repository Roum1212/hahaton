<script setup lang="ts">
import { USER_ROLE_LABELS, type IUser } from '~/core/users/services/users/users.types';

const usersStore = useUsersStore();
const authStore = useAuthStore();
const { state, load, reload, search, deleteMany } = useUsersActions();
const { openUserFormDialog } = useUsersDialogs();
const { columns } = useUsersTable();

const selected = ref<(string | number)[]>([]);
const isSelf = (user: IUser) => user.id === authStore.user?.id;

const removeSelected = () => deleteMany(selected.value.map(String), () => (selected.value = []));

onBeforeMount(reload);
</script>

<template>
  <UiTable
    v-model:selected="selected"
    :items="usersStore.list"
    :columns="columns"
    :loading="state.loading"
    :fetching="state.fetching"
    :total="usersStore.total"
    :load-all="usersStore.pagination.loadAll"
    empty-title="Пользователей нет"
    empty-text="Измените условия поиска или добавьте пользователя"
    empty-icon="pi pi-users"
    @load-more="load">
    <template #toolbar>
      <IconField class="grow sm:grow-0">
        <InputIcon class="pi pi-search" />
        <InputText
          v-model="state.searchQuery"
          autocomplete="off"
          placeholder="Поиск по имени или email"
          class="w-full sm:w-[280px]"
          @input="search" />
      </IconField>
      <UiButton icon="pi pi-refresh" severity="secondary" outlined :loading="state.loading" @click="reload" />
    </template>

    <template #cell-name="{ row }">
      <div class="app-table__primary">
        {{ row.fullName }}
        <Tag v-if="isSelf(row)" value="Вы" severity="info" rounded />
      </div>
      <div class="app-table__muted">{{ row.email }}</div>
    </template>

    <template #cell-role="{ row }">
      <Tag :value="USER_ROLE_LABELS[row.role]" severity="secondary" rounded />
    </template>

    <template #cell-status="{ row }">
      <Tag :severity="row.active ? 'success' : 'danger'" :value="row.active ? 'Активен' : 'Заблокирован'" rounded />
    </template>

    <template #cell-createdAt="{ row }">
      <span class="app-table__muted">{{ formatDate(row.createdAt, 'dd.MM.yyyy • HH:mm') }}</span>
    </template>

    <template #cell-actions="{ row }">
      <div class="app-table__actions">
        <UiButton
          v-tooltip.top="'Редактировать'"
          icon="pi pi-pencil"
          severity="secondary"
          text
          rounded
          @click="openUserFormDialog(row.id)" />
        <UiButton
          v-tooltip.top="'Удалить'"
          icon="pi pi-trash"
          severity="danger"
          text
          rounded
          :disabled="isSelf(row)"
          @click="deleteMany([row.id])" />
      </div>
    </template>

    <template #bulk>
      <UiButton label="Удалить" icon="pi pi-trash" severity="danger" size="small" @click="removeSelected" />
    </template>
  </UiTable>
</template>
