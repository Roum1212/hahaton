<script setup lang="ts">
import type { ILandingPhoto } from '~/config/landing.config';

withDefaults(
  defineProps<{
    photo: ILandingPhoto;
    ratio?: string;
    eager?: boolean;
  }>(),
  {
    ratio: '3 / 2',
    eager: false,
  },
);
</script>

<template>
  <figure class="landing-photo">
    <div class="landing-photo__frame" :class="{ 'landing-photo__frame--contain': photo.fit === 'contain' }" :style="{ aspectRatio: ratio }">
      <img :src="photo.src" :alt="photo.alt" :loading="eager ? 'eager' : 'lazy'" class="landing-photo__img" />
      <slot />
    </div>
    <figcaption v-if="photo.source" class="landing-photo__caption">
      <a v-if="photo.href" :href="photo.href" target="_blank" rel="noopener noreferrer">{{ photo.source }}</a>
      <template v-else>{{ photo.source }}</template>
    </figcaption>
  </figure>
</template>
