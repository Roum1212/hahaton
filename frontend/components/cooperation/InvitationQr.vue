<script setup lang="ts">
import QRCode from 'qrcode';

const props = defineProps<{
  url: string;
}>();

const image = ref('');

watch(
  () => props.url,
  async (value) => {
    image.value = value ? await QRCode.toDataURL(value, { width: 280, margin: 1 }) : '';
  },
  { immediate: true },
);
</script>

<template>
  <img v-if="image" :src="image" alt="QR-код приглашения" class="app-qr" />
</template>
