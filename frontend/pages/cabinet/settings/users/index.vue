<script setup lang="ts">
definePageMeta({
  title: 'Пользователи',
  permissions: ['admin'],
});

const { openUserFormDialog } = useUsersDialogs();
</script>

<template>
  <AppPage>
    <template #title>Пользователи</template>
    <template #subtitle>Учётные записи, роли и доступ</template>
    <template #actions>
      <UiButton label="Добавить" icon="pi pi-user-plus" @click="openUserFormDialog()" />
    </template>

    <UsersList />
  </AppPage>
</template>
