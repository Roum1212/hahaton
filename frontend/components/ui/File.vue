<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    src: string | undefined | null;
    name?: string | null;
    containerClass?: string;
    classList?: string;
    disabled?: boolean;
  }>(),
  {
    disabled: false,
    name: '',
    containerClass: '',
    classList: 'w-full h-60',
  },
);
const file = ref(PLACEHOLDER_IMAGE);

const extension = ref<string | null>(null);
watch(
  () => props.src,
  () => {
    if (!props.src) {
      file.value = PLACEHOLDER_IMAGE;
      extension.value = null;
    } else if (isServerBase64(props.src)) {
      file.value = convertFromServerBase64(props.src);
      extension.value = getExtension(props.src);
    } else if (props.src !== PLACEHOLDER_IMAGE) {
      file.value = getImageUrl() + props.src;
      extension.value = getExtension(props.src);
    } else {
      file.value = props.src;
      extension.value = null;
    }
  },
  { immediate: true },
);
</script>
<template>
  <div v-if="file" class="flex flex-col items-center justify-center relative" :class="containerClass">
    <div class="grow-0">
      <img :src="`/images/files/file-${extension ? extension : 'default'}.svg`" class="w-16 h-16" />
    </div>
    <div v-if="name" class="w-full text-sm text-center font-medium mb-4">{{ truncateText(name, 7, 5) }}</div>
  </div>
</template>
