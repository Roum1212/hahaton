<script setup lang="ts">
import { LANDING_NAV } from '~/config/landing.config';

const { siteName, isLoggedIn, primaryRoute, scrollTo } = useLanding();
const { city, label, confirmed, confirm, openDialog } = useSearchArea();
const { y } = useWindowScroll();

const menuOpen = ref(false);
const scrolled = computed(() => y.value > 8);

const onNavigate = (anchor: string) => {
  menuOpen.value = false;
  scrollTo(anchor);
};
</script>

<template>
  <header class="landing-header" :class="{ 'landing-header--scrolled': scrolled || menuOpen }">
    <div class="landing-container landing-header__inner">
      <NuxtLink :to="ROUTES.home" class="landing-header__brand" @click="menuOpen = false">
        <img src="/images/common/volonterim_logo_transparent.png" :alt="siteName" class="landing-header__logo" />
        <img src="/images/common/wordmark.png" :alt="siteName" class="landing-header__wordmark" />
      </NuxtLink>

      <div class="landing-city">
        <button type="button" class="landing-city__button" @click="openDialog">
          <i class="pi pi-map-marker" />
          <span class="landing-city__label">{{ label }}</span>
        </button>
        <Transition name="landing-menu">
          <div v-if="!confirmed" class="landing-city__confirm">
            <span class="landing-city__question">Ваш город — {{ city.name }}?</span>
            <div class="landing-city__actions">
              <Button label="Да" size="small" @click="confirm" />
              <Button label="Изменить" size="small" severity="secondary" outlined @click="openDialog" />
            </div>
          </div>
        </Transition>
      </div>

      <nav class="landing-header__nav">
        <a
          v-for="item in LANDING_NAV"
          :key="item.anchor"
          :href="`#${item.anchor}`"
          class="landing-header__link"
          @click.prevent="onNavigate(item.anchor)">
          {{ item.label }}
        </a>
      </nav>

      <div class="landing-header__actions">
        <NuxtLink v-if="!isLoggedIn" :to="ROUTES.login" class="landing-header__login">
          <Button label="Войти" severity="secondary" outlined size="small" />
        </NuxtLink>
        <NuxtLink :to="primaryRoute">
          <Button :label="isLoggedIn ? 'Личный кабинет' : 'Стать волонтёром'" size="small" />
        </NuxtLink>
        <button type="button" class="landing-header__burger" :aria-label="menuOpen ? 'Закрыть меню' : 'Меню'" @click="menuOpen = !menuOpen">
          <i :class="menuOpen ? 'pi pi-times' : 'pi pi-bars'" />
        </button>
      </div>
    </div>

    <Transition name="landing-menu">
      <nav v-if="menuOpen" class="landing-header__mobile">
        <a
          v-for="item in LANDING_NAV"
          :key="item.anchor"
          :href="`#${item.anchor}`"
          class="landing-header__mobile-link"
          @click.prevent="onNavigate(item.anchor)">
          {{ item.label }}
          <i class="pi pi-angle-right" />
        </a>
        <NuxtLink v-if="!isLoggedIn" :to="ROUTES.login" class="landing-header__mobile-link" @click="menuOpen = false">
          Войти
          <i class="pi pi-sign-in" />
        </NuxtLink>
      </nav>
    </Transition>

    <GeoSearchAreaDialog />
  </header>
</template>
