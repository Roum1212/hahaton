<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts">
defineProps<{
  modal: DialogInstance<any, any>;
}>();
const siteUrl = useRuntimeConfig().public.siteUrl;
const { minimize } = useDialogManager();

const { copy } = useClipboard('Ссылка скопиррована в буфер обмена');
const handleOpenWindow = (url: string | undefined) => {
  if (!url) return;
  window.open(url, '_blank', 'noopener,noreferrer');
};
</script>
<template>
  <div class="flex items-center">
    <Button
      v-if="modal.options?.detailUrl"
      severity="secondary"
      icon="pi pi-external-link"
      label="Перейти"
      class="!px-0"
      text
      size="small"
      @click.stop="handleOpenWindow(modal.options.detailUrl)" />
    <Button
      v-if="modal.options?.detailUrl"
      icon="pi pi-copy"
      class="!px-0"
      size="small"
      text
      severity="secondary"
      @click.stop="copy(`${siteUrl}${modal.options.detailUrl}`)" />
    <Button
      v-if="modal.options?.minimizable"
      icon="pi pi-window-minimize"
      class="ml-auto"
      size="small"
      text
      severity="secondary"
      @click.stop="minimize(modal.id)" />
  </div>
</template>
