<script setup lang="ts">
withDefaults(
  defineProps<{
    modelValue: string | number;
    text?: string | number;
    icon?: string;
    buttonClass?: string;
    hoverable?: boolean;
  }>(),
  {
    text: undefined,
    buttonClass: '',
    icon: 'pi pi-copy',
    hoverable: true,
  },
);

const { toast } = useApp();
function handleCopy(data: string | number) {
  navigator.clipboard
    .writeText(data.toString())
    .then(() => {
      toast.info('Cкопировано');
    })
    .catch((err) => {
      console.error('Ошибка копирования: ', err);
    });
}
</script>
<template>
  <div class="flex gap-1 items-center group">
    <slot v-if="$slots.text" name="text" />
    <span v-else>{{ text }}</span>
    <Button
      size="small"
      :icon="icon"
      class="!w-6 !h-6 p-0"
      :class="[...buttonClass, hoverable ? 'opacity-0 group-hover:opacity-100 transition-opacity duration-200' : '']"
      severity="secondary"
      text
      rounded
      @click.stop="handleCopy(modelValue)" />
  </div>
</template>
