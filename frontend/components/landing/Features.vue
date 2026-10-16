<script setup lang="ts">
import { LANDING_FEATURES } from '~/config/landing.config';
</script>

<template>
  <section class="landing-section landing-section--muted">
    <div class="landing-container">
      <h2 class="landing-title">Как ВОЛОНТЕРИМ решает эти задачи</h2>

      <div class="landing-features">
        <article v-for="item in LANDING_FEATURES" :key="item.title" class="landing-card landing-card--feature">
          <LandingPhoto v-if="item.photo" :photo="item.photo" ratio="4 / 3" class="landing-card__photo" />
          <span v-else class="landing-card__icon">
            <component :is="item.icon" :size="22" />
          </span>
          <h3 class="landing-card__title">{{ item.title }}</h3>
          <p class="landing-card__text">{{ item.text }}</p>
          <div class="landing-tags">
            <span v-for="tag in item.tags" :key="tag" class="landing-tag">{{ tag }}</span>
          </div>
        </article>
      </div>
    </div>
  </section>
</template>
