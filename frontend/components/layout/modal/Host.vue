<script setup lang="ts">
import Dialog from 'primevue/dialog';

const { dialogs, close } = useDialogManager();

const hiddenDialogs = computed(() => {
  return Object.values(dialogs.value).filter((el) => el && !el?.visible);
});
</script>

<template>
  <Teleport to="body">
    <LayoutModalMinimize :dialogs="hiddenDialogs" />

    <template v-for="modal in Object.values(dialogs)" :key="modal?.id">
      <Dialog
        v-if="(modal && !modal.options?.type) || modal?.options?.type === 'modal'"
        v-model:visible="modal.visible"
        :header="modal.options?.header || ''"
        :class="modal.options?.class"
        :draggable="false"
        :position="modal.options?.position || 'center'"
        :modal="modal.options?.isModal || true"
        @update:visible="(v: boolean) => !v && close(modal.id)">
        <template #header>
          <component :is="modal.headerContent(modal)" v-if="modal.headerContent" />
          <div v-else class="flex justify-between w-full">
            <span class="p-dialog-title">{{ modal.options?.header }}</span>
            <LayoutModalHeaderControls :modal="modal" />
          </div>
        </template>
        <component
          :is="modal.component"
          :ref="(el: any) => (modal.componentRef = el)"
          v-bind="modal.componentProps"
          @close="close(modal.id)" />
        <template #footer> <component :is="modal.footerContent(modal)" v-if="modal.footerContent" /> </template>
      </Dialog>

      <LayoutSidebarPanel
        v-else-if="modal && modal.options?.type === 'sidebar'"
        v-model:visible="modal.visible"
        :header="modal.options?.header"
        :size="modal.options?.size || 'small'"
        :position="modal.options?.position || 'right'"
        :class="modal.options?.class"
        @update:visible="(v: boolean) => !v && close(modal.id)">
        <template #header>
          <component :is="modal.headerContent(modal)" v-if="modal.headerContent" />
          <div v-else class="flex justify-between w-full">
            <span class="p-dialog-title">{{ modal.options?.header }}</span>
            <LayoutModalHeaderControls :modal="modal" />
          </div>
        </template>
        <template #body>
          <component
            :is="modal.component"
            :ref="(el: any) => (modal.componentRef = el)"
            v-bind="modal.componentProps"
            @close="close(modal.id)" />
        </template>

        <template v-if="!!modal.footerContent" #footer>
          <component :is="modal.footerContent(modal)" v-if="modal.footerContent" />
        </template>
      </LayoutSidebarPanel>
    </template>
  </Teleport>
</template>
