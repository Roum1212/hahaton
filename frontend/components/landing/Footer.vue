<script setup lang="ts">
import { COINS_DISCLAIMER, LANDING_NAV } from '~/config/landing.config';

const { siteName, scrollTo } = useLanding();
const year = new Date().getFullYear();
</script>

<template>
  <footer class="landing-footer">
    <div class="landing-container landing-footer__inner">
      <div class="landing-footer__brand">
        <NuxtLink :to="ROUTES.home" class="landing-header__brand">
          <img src="/images/common/volonterim_logo_transparent.png" :alt="siteName" class="landing-header__logo" />
          <img src="/images/common/wordmark.png" :alt="siteName" class="landing-header__wordmark" />
        </NuxtLink>
        <p class="landing-footer__note">{{ COINS_DISCLAIMER }}</p>
      </div>

      <nav class="landing-footer__nav">
        <a v-for="item in LANDING_NAV" :key="item.anchor" :href="`#${item.anchor}`" @click.prevent="scrollTo(item.anchor)">
          {{ item.label }}
        </a>
      </nav>

      <div class="landing-footer__nav">
        <NuxtLink :to="ROUTES.login">Войти</NuxtLink>
        <NuxtLink :to="ROUTES.cabinet">Личный кабинет</NuxtLink>
      </div>
    </div>

    <div class="landing-container landing-footer__bottom">
      <span>© {{ year }} {{ siteName }}</span>
      <span>Обнинск</span>
    </div>
  </footer>
</template>
