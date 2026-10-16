<script setup lang="ts">
import { LANDING_PHOTOS, type ILandingPhoto } from '~/config/landing.config';

const { primaryRoute, isLoggedIn, scrollTo } = useLanding();

const photo: ILandingPhoto = LANDING_PHOTOS.team;
</script>

<template>
  <section class="landing-hero">
    <img :src="photo.src" :alt="photo.alt" class="landing-hero__bg" />
    <div class="landing-hero__overlay" />

    <div class="landing-container landing-hero__inner">
      <h1 class="landing-hero__title">
        Цифровой волонтёр
        <span class="landing-hero__accent">Обнинска</span>
      </h1>
      <p class="landing-hero__lead">Учёт добрых дел, волонтёрских часов и понятная система поощрений</p>
      <div class="landing-hero__actions">
        <NuxtLink :to="primaryRoute">
          <Button :label="isLoggedIn ? 'Перейти в кабинет' : 'Стать волонтёром'" icon="pi pi-arrow-right" icon-pos="right" size="large" />
        </NuxtLink>
        <Button label="Смотреть мероприятия" severity="secondary" outlined size="large" @click="scrollTo('events')" />
      </div>
    </div>

    <span v-if="photo.source" class="landing-hero__caption">
      <a v-if="photo.href" :href="photo.href" target="_blank" rel="noopener noreferrer">{{ photo.source }}</a>
      <template v-else>{{ photo.source }}</template>
    </span>
  </section>
</template>
