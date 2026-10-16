<script setup lang="ts">
import { LANDING_STATS } from '~/config/landing.config';

const { scrollTo } = useLanding();
</script>

<template>
  <section class="landing-section landing-section--tight">
    <div class="landing-container">
      <div class="landing-stats">
        <div v-for="stat in LANDING_STATS" :key="stat.label" class="landing-stats__item">
          <span class="landing-stats__value">{{ stat.value }}</span>
          <span class="landing-stats__label">
            {{ stat.label }}
            <i
              v-if="stat.hint"
              v-tooltip.bottom="{ value: stat.hint, class: 'landing-tooltip' }"
              class="pi pi-question-circle landing-hint" />
          </span>
          <button type="button" class="landing-stats__link" @click="scrollTo(stat.anchor)">
            {{ stat.link }}
            <i class="pi pi-arrow-right" />
          </button>
        </div>
      </div>
    </div>
  </section>
</template>
