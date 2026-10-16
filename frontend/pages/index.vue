<script setup lang="ts">
definePageMeta({
  layout: false,
});

const { siteName } = useLanding();

useHead({ title: `Цифровой волонтёр Обнинска | ${siteName}` });
</script>

<template>
  <div class="landing">
    <LandingHeader />
    <main>
      <LandingHero />
      <LandingStats />
      <LandingProblems />
      <LandingFeatures />
      <LandingSteps />
      <LandingEvents />
      <LandingAudience />
      <LandingRules />
      <LandingCta />
      <LandingFaq />
    </main>
    <LandingFooter />
  </div>
</template>
