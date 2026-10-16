<script setup lang="ts">
import { COINS_PER_HOUR, LANDING_PHOTOS } from '~/config/landing.config';

const { isLoggedIn } = useLanding();
const { area, inCityTitle, openDialog } = useSearchArea();
const { store, loading, load } = useLandingEvents();

watch(area, (value) => load(value), { immediate: true, deep: true });
</script>

<template>
  <section id="events" class="landing-section landing-section--muted">
    <div class="landing-container">
      <div class="landing-heading">
        <h2 class="landing-title">Мероприятия {{ inCityTitle }}</h2>
        <NuxtLink v-if="store.upcoming.length" :to="ROUTES.volunteer.events" class="landing-heading__link">
          Все мероприятия
          <i class="pi pi-arrow-right" />
        </NuxtLink>
      </div>

      <div v-if="loading" class="landing-events">
        <Skeleton v-for="n in 3" :key="n" height="420px" border-radius="20px" />
      </div>

      <div v-else-if="!isLoggedIn" class="landing-empty">
        <img :src="LANDING_PHOTOS.eventsEmpty.src" :alt="LANDING_PHOTOS.eventsEmpty.alt" class="landing-empty__image" />
        <h3 class="landing-empty__title">Войдите, чтобы увидеть мероприятия координаторов</h3>
        <NuxtLink :to="ROUTES.login">
          <Button label="Войти" icon="pi pi-sign-in" />
        </NuxtLink>
      </div>

      <div v-else-if="!store.upcoming.length" class="landing-empty">
        <img :src="LANDING_PHOTOS.eventsEmpty.src" :alt="LANDING_PHOTOS.eventsEmpty.alt" class="landing-empty__image" />
        <h3 class="landing-empty__title">Ближайшие мероприятия не найдены</h3>
        <Button label="Изменить город" icon="pi pi-map-marker" severity="secondary" outlined @click="openDialog" />
      </div>

      <div v-else class="landing-events">
        <article v-for="event in store.upcoming" :key="event.id" class="landing-event">
          <LandingPhoto :photo="LANDING_PHOTOS.stepCatalog" ratio="3 / 2" class="landing-event__photo" />
          <div class="landing-event__body">
            <span class="landing-tag">Свободно {{ event.freePlaces }} из {{ event.volunteerLimit }}</span>
            <h3 class="landing-event__title">{{ event.title }}</h3>
            <ul class="landing-event__meta">
              <li>
                <i class="pi pi-calendar" />
                {{ formatDate(event.startsAt, 'd MMMM yyyy, HH:mm') }}
              </li>
              <li v-if="event.location">
                <i class="pi pi-map-marker" />
                {{ event.location }}
              </li>
            </ul>
            <div class="landing-event__footer">
              <div class="landing-event__reward">
                <span class="landing-event__coins">+{{ (event.durationMinutes / 60) * COINS_PER_HOUR }} койнов</span>
                <span class="landing-event__note">за {{ event.durationMinutes / 60 }} ч участия</span>
              </div>
              <NuxtLink :to="`${ROUTES.volunteer.events}/${event.id}`">
                <Button label="Помочь" icon="pi pi-heart" />
              </NuxtLink>
            </div>
          </div>
        </article>
      </div>
    </div>
  </section>
</template>
