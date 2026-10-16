<script setup lang="ts">
const { primaryRoute, isLoggedIn } = useLanding();
</script>

<template>
  <section class="landing-section landing-section--tight">
    <div class="landing-container">
      <div class="landing-cta">
        <div class="landing-cta__content">
          <h2 class="landing-cta__title">Начните делать добро — уже сегодня</h2>
          <NuxtLink :to="primaryRoute">
            <Button :label="isLoggedIn ? 'Перейти в кабинет' : 'Стать волонтёром'" icon="pi pi-arrow-right" icon-pos="right" size="large" />
          </NuxtLink>
        </div>
        <div class="landing-cta__art">
          <img src="/images/common/volonterim_logo_transparent.png" alt="" class="landing-cta__logo" />
        </div>
      </div>
    </div>
  </section>
</template>
