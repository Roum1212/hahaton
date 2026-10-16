<script setup lang="ts">
import { LANDING_AUDIENCE } from '~/config/landing.config';

const { primaryRoute } = useLanding();

const activeKey = ref(LANDING_AUDIENCE[0]!.key);
const active = computed(() => LANDING_AUDIENCE.find((item) => item.key === activeKey.value) ?? LANDING_AUDIENCE[0]!);
</script>

<template>
  <section id="audience" class="landing-section landing-section--muted">
    <div class="landing-container">
      <h2 class="landing-title">Добро для всех</h2>

      <div class="landing-audience">
        <div class="landing-audience__tabs" role="tablist">
          <button
            v-for="item in LANDING_AUDIENCE"
            :key="item.key"
            type="button"
            role="tab"
            :aria-selected="item.key === activeKey"
            class="landing-audience__tab"
            :class="{ 'landing-audience__tab--active': item.key === activeKey }"
            @click="activeKey = item.key">
            {{ item.label }}
            <i class="pi pi-angle-right" />
          </button>
        </div>

        <Transition name="landing-fade" mode="out-in">
          <div :key="active.key" class="landing-audience__panel">
            <div class="landing-audience__content">
              <h3 class="landing-audience__title">{{ active.title }}</h3>
              <ul class="landing-audience__points">
                <li v-for="point in active.points" :key="point">
                  <i class="pi pi-check-circle" />
                  {{ point }}
                </li>
              </ul>
              <NuxtLink :to="primaryRoute">
                <Button :label="active.action" icon="pi pi-arrow-right" icon-pos="right" />
              </NuxtLink>
            </div>
            <LandingPhoto :photo="active.photo" ratio="4 / 3" class="landing-audience__photo" />
          </div>
        </Transition>
      </div>
    </div>
  </section>
</template>
