<script setup lang="ts">
import { EVENT_STATUS_LABELS, EVENT_STATUS_SEVERITY, EventStatus, type IEvent } from '~/core/event/services/event/event.types';

const props = defineProps<{
  event: IEvent;
  coordinatorName?: string;
  to?: string;
}>();

const coins = computed(() => coinsForMinutes(props.event.durationMinutes));
const image = computed(() => pickIllustration(EVENT_ILLUSTRATIONS, props.event.id));
</script>

<template>
  <article class="app-event-card">
    <component :is="to ? resolveComponent('NuxtLink') : 'div'" :to="to" class="app-event-card__image">
      <img :src="image" alt="" />
      <div class="app-event-card__badges">
        <AppStatus
          v-if="event.status !== EventStatus.PUBLISHED"
          :label="EVENT_STATUS_LABELS[event.status]"
          :severity="EVENT_STATUS_SEVERITY[event.status]" />
        <span v-else-if="event.freePlaces <= 0" class="app-chip">Мест нет</span>
      </div>
    </component>
    <div class="app-event-card__body">
      <span class="app-event-card__date">{{ formatEventPeriod(event.startsAt, event.endsAt) }}</span>
      <NuxtLink v-if="to" :to="to" class="app-event-card__title">{{ event.title }}</NuxtLink>
      <span v-else class="app-event-card__title">{{ event.title }}</span>
      <ul class="app-event-card__meta">
        <li v-if="coordinatorName"><i class="pi pi-user" />{{ coordinatorName }}</li>
        <li v-if="event.location"><i class="pi pi-map-marker" />{{ event.location }}</li>
        <li><i class="pi pi-users" />Свободно {{ event.freePlaces }} из {{ event.volunteerLimit }}</li>
      </ul>
      <div class="app-event-card__footer">
        <span class="app-event-card__coins">+{{ formatCoins(coins) }}</span>
        <slot name="actions" />
      </div>
    </div>
  </article>
</template>
