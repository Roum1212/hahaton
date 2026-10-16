<script setup lang="ts">
import type { IEvent } from '~/core/event/services/event/event.types';

defineProps<{
  event: IEvent;
}>();
</script>

<template>
  <div v-if="event.point || event.location" class="app-card">
    <div class="app-card__head">
      <span class="app-card__title">Место проведения</span>
      <a v-if="event.point" :href="yandexMapsLink(event.point)" target="_blank" rel="noopener" class="app-link">
        Открыть в Яндекс Картах
        <i class="pi pi-external-link text-xs" />
      </a>
    </div>
    <span class="app-location__address">
      <i class="pi pi-map-marker" />
      {{ event.location || 'Адрес уточнит координатор' }}
    </span>
    <GeoMap v-if="event.point" :center="event.point" :zoom="16" :interactive="false" class="app-location__map" />
  </div>
</template>
