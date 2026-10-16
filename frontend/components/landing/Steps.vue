<script setup lang="ts">
import { LANDING_STEPS } from '~/config/landing.config';
</script>

<template>
  <section id="steps" class="landing-section">
    <div class="landing-container">
      <h2 class="landing-title">Как это работает</h2>

      <ol class="landing-steps">
        <li v-for="(step, idx) in LANDING_STEPS" :key="step.title" class="landing-steps__item">
          <LandingPhoto v-if="step.photo" :photo="step.photo" ratio="4 / 3" class="landing-steps__photo">
            <span class="landing-steps__index landing-steps__index--badge">Шаг {{ idx + 1 }}</span>
          </LandingPhoto>
          <div v-else class="landing-steps__head">
            <span class="landing-steps__index">Шаг {{ idx + 1 }}</span>
            <span class="landing-steps__icon">
              <component :is="step.icon" :size="20" />
            </span>
          </div>
          <h3 class="landing-card__title">{{ step.title }}</h3>
          <p class="landing-card__text">{{ step.text }}</p>
          <span class="landing-steps__chip">{{ step.chip }}</span>
        </li>
      </ol>
    </div>
  </section>
</template>
