<script setup lang="ts">
import { LANDING_PHOTOS, LANDING_RULES, type ILandingPhoto } from '~/config/landing.config';

const { primaryRoute, isLoggedIn } = useLanding();

const photo: ILandingPhoto = LANDING_PHOTOS.banner;
</script>

<template>
  <section class="landing-section landing-section--tight">
    <div class="landing-container">
      <div class="landing-banner">
        <div class="landing-banner__content">
          <h2 class="landing-banner__title">Каждый час — на счету</h2>
          <ul class="landing-banner__rules">
            <li v-for="rule in LANDING_RULES" :key="rule.text">
              <i :class="rule.icon" />
              {{ rule.text }}
            </li>
          </ul>
          <NuxtLink :to="primaryRoute" class="landing-banner__action">
            <Button :label="isLoggedIn ? 'Перейти в кабинет' : 'Стать волонтёром'" icon="pi pi-arrow-right" icon-pos="right" size="large" />
          </NuxtLink>
        </div>
        <div class="landing-banner__media">
          <img :src="photo.src" :alt="photo.alt" loading="lazy" class="landing-banner__img" />
          <span v-if="photo.source" class="landing-banner__caption">{{ photo.source }}</span>
        </div>
      </div>
    </div>
  </section>
</template>
