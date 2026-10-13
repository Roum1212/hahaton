<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts">
defineProps<{ dialogs: DialogInstance<any, any>[] }>();
const { close, revert } = useDialogManager();
const collapsedPanelRef = ref();

const onToggle = (e: Event) => {
  collapsedPanelRef.value.toggle(e);
};
</script>
<template>
  <Teleport to="body">
    <div class="fixed bottom-2 right-3 z-300 flex flex-row-reverse justify-start w-full gap-1">
      <Button v-if="dialogs.length" severity="secondary" label="Свернутые окна" icon="pi pi-clone" @click.stop="onToggle" />
      <Popover ref="collapsedPanelRef">
        <div class="flex flex-col gap-2">
          <template v-for="h of dialogs" :key="h?.id">
            <div
              v-tooltip.top="h.options?.header"
              class="max-w-[16rem] bg-surface-200 flex gap-2 items-center rounded-lg px-2 py-1 cursor-pointer relative hover:bg-surface-300 justify-between text-white dark:text-surface-900"
              @click="revert(h.id)">
              <i class="pi pi-external-link !text-xs" />
              <span class="overflow-hidden whitespace-nowrap">{{ h.options?.header }}</span>
              <Button icon="pi pi-times" severity="secondary" size="small" text @click.stop="close(h.id)" />
            </div>
          </template>
        </div>
      </Popover>
    </div>
  </Teleport>
</template>
