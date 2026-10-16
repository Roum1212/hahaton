<script setup lang="ts">
import { LANDING_PHOTOS, LANDING_PROBLEMS } from '~/config/landing.config';
</script>

<template>
  <section id="about" class="landing-section">
    <div class="landing-container">
      <h2 class="landing-title">Почему добрые дела теряются без учёта</h2>

      <LandingPhoto :photo="LANDING_PHOTOS.words" ratio="3 / 1" class="landing-problems__banner" />

      <div class="landing-problems">
        <article v-for="item in LANDING_PROBLEMS" :key="item.title" class="landing-card">
          <LandingPhoto v-if="item.photo" :photo="item.photo" ratio="1 / 1" class="landing-card__thumb" />
          <span v-else class="landing-card__icon landing-card__icon--muted">
            <component :is="item.icon" :size="22" />
          </span>
          <h3 class="landing-card__title">{{ item.title }}</h3>
          <p class="landing-card__text">{{ item.text }}</p>
          <div class="landing-tags">
            <span v-for="tag in item.tags" :key="tag" class="landing-tag landing-tag--danger">{{ tag }}</span>
          </div>
        </article>
      </div>
    </div>
  </section>
</template>
