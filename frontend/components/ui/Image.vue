<script setup lang="ts">
import Image from 'primevue/image';
import type { Image as ImageType } from '~/core/app/types/index';

import { PLACEHOLDER_IMAGE, convertFromServerBase64, isServerBase64, validateImage } from '~/utils/image.utils';
const props = withDefaults(
  defineProps<{
    src?: ImageType;
    containerClass?: string;
    imageClass?: string;
    type?: 'thumbnail' | 'full';
    preview?: boolean;
    alt?: string;
    previewPrefix?: string;
  }>(),
  {
    src: undefined,
    containerClass: '',
    imageClass: 'w-full object-contain',
    type: 'thumbnail',
    preview: false,
    alt: '',
    previewPrefix: 'p-',
  },
);

const full = ref(PLACEHOLDER_IMAGE);
const thumbnail = ref(PLACEHOLDER_IMAGE);
const enabledPreview = ref<boolean>(props.preview || false);
watch(
  () => props.src,
  () => {
    if (!props.src) {
      full.value = PLACEHOLDER_IMAGE;
      thumbnail.value = full.value;
      enabledPreview.value = false;
    } else if (isServerBase64(props.src)) {
      full.value = convertFromServerBase64(props.src);
      thumbnail.value = full.value;
    } else if (props.src !== PLACEHOLDER_IMAGE) {
      full.value = getImageUrl() + '' + props.src;
      thumbnail.value = getImageUrl() + props.previewPrefix + props.src;
      if (!validateImage(full.value)) full.value = PLACEHOLDER_IMAGE;
      if (!validateImage(thumbnail.value)) thumbnail.value = PLACEHOLDER_IMAGE;
    } else {
      full.value = props.src;
      thumbnail.value = full.value;
    }
  },
  { immediate: true },
);
</script>
<template>
  <Image :preview="enabledPreview" class="flex items-center image-field" :class="containerClass">
    <template #image>
      <img :src="thumbnail" :class="'image-field__image ' + props.imageClass" />
    </template>
    <template #original>
      <img :src="full" class="w-full" />
    </template>
  </Image>
</template>
<style scoped lang="scss">
.image-field {
  max-width: 100%;

  display: flex;
  align-items: center;
  justify-content: center;
  &__image {
    width: 100%;
    height: 100%;
    object-fit: contain;
    object-position: center center;
  }
}
</style>
