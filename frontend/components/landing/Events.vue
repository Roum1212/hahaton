<script setup lang="ts">
import { COINS_PER_HOUR, LANDING_PHOTOS } from '~/config/landing.config';

const { primaryRoute } = useLanding();
const { city, current, openDialog } = useCity();
const { store, loading, load } = useUpcomingEvents();

watch(city, (value) => load(value), { immediate: true });
</script>

<template>
  <section id="events" class="landing-section landing-section--muted">
    <div class="landing-container">
      <div class="landing-heading">
        <h2 class="landing-title">Мероприятия в {{ current.locative }}</h2>
        <NuxtLink v-if="store.upcoming.length" :to="primaryRoute" class="landing-heading__link">
          Все мероприятия
          <i class="pi pi-arrow-right" />
        </NuxtLink>
      </div>

      <div v-if="loading" class="landing-events">
        <Skeleton v-for="n in 3" :key="n" height="420px" border-radius="20px" />
      </div>

      <div v-else-if="!store.upcoming.length" class="landing-empty">
        <img :src="LANDING_PHOTOS.eventsEmpty.src" :alt="LANDING_PHOTOS.eventsEmpty.alt" class="landing-empty__image" />
        <h3 class="landing-empty__title">Ближайшие мероприятия не найдены</h3>
        <Button label="Изменить город" icon="pi pi-map-marker" severity="secondary" outlined @click="openDialog" />
      </div>

      <div v-else class="landing-events">
        <article v-for="event in store.upcoming" :key="event.id" class="landing-event">
          <LandingPhoto :photo="{ src: event.image ?? PLACEHOLDER_IMAGE, alt: event.title }" ratio="3 / 2" class="landing-event__photo" />
          <div class="landing-event__body">
            <span class="landing-tag">{{ event.category }}</span>
            <h3 class="landing-event__title">{{ event.title }}</h3>
            <ul class="landing-event__meta">
              <li>
                <i class="pi pi-calendar" />
                {{ formatDate(event.startsAt, 'd MMMM yyyy, HH:mm') }}
              </li>
              <li>
                <i class="pi pi-map-marker" />
                {{ event.place }}
              </li>
              <li>
                <i class="pi pi-users" />
                {{ event.slots }} мест
              </li>
            </ul>
            <div class="landing-event__footer">
              <div class="landing-event__reward">
                <span class="landing-event__coins">+{{ event.hours * COINS_PER_HOUR }} койнов</span>
                <span class="landing-event__note">за {{ event.hours }} ч участия</span>
              </div>
              <NuxtLink :to="primaryRoute">
                <Button label="Помочь" icon="pi pi-heart" />
              </NuxtLink>
            </div>
          </div>
        </article>
      </div>
    </div>
  </section>
</template>
