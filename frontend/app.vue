<template>
  <Toast position="top-right">
    <template #message="{ message }">
      <div class="flex w-full flex-col gap-4">
        <div>
          <strong>{{ message.summary }}</strong>
          <div>{{ message.detail }}</div>
        </div>
        <div v-if="message.data?.buttons?.length" class="flex gap-2">
          <Button
            v-for="(button, idx) in message.data.buttons"
            :key="idx"
            :label="button.label"
            outlined
            :severity="button.severity || 'secondary'"
            size="small"
            @click.stop="
              () => {
                button.action?.();
                message.closeCallback?.();
              }
            " />
        </div>
      </div>
    </template>
  </Toast>

  <LayoutLoader />
  <LayoutModalHost />

  <ConfirmDialog group="confirmDelete">
    <template #container="{ message, acceptCallback, rejectCallback }">
      <div class="flex flex-col items-center rounded-border bg-surface-0 p-8 dark:bg-surface-900">
        <div class="-mt-20 inline-flex h-24 w-24 items-center justify-center rounded-full bg-red-400 text-white">
          <i class="pi pi-question text-5xl" />
        </div>
        <span class="mb-2 mt-6 block text-2xl font-bold">{{ message.header }}</span>
        <p class="mb-0">{{ message.message }}</p>
        <div class="mt-6 flex items-center gap-2">
          <Button :label="message.acceptLabel" class="min-w-8rem" severity="danger" @click="acceptCallback" />
          <Button :label="message.rejectLabel" severity="secondary" class="w-32" @click="rejectCallback" />
        </div>
      </div>
    </template>
  </ConfirmDialog>

  <ConfirmDialog group="confirmSave">
    <template #container="{ message, acceptCallback, rejectCallback }">
      <div class="flex flex-col items-center rounded-border bg-surface-0 p-8 dark:bg-surface-900">
        <div class="-mt-20 inline-flex h-24 w-24 items-center justify-center rounded-full bg-green-400 text-white">
          <i class="pi pi-question text-5xl" />
        </div>
        <span class="mb-2 mt-6 block text-2xl font-bold">{{ message.header }}</span>
        <p class="mb-0">{{ message.message }}</p>
        <div class="mt-6 flex items-center gap-2">
          <Button :label="message.acceptLabel" class="min-w-8rem" severity="success" @click="acceptCallback" />
          <Button :label="message.rejectLabel" severity="secondary" class="w-32" @click="rejectCallback" />
        </div>
      </div>
    </template>
  </ConfirmDialog>

  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>
