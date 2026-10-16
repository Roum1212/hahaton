<script setup lang="ts">
import { LANDING_FAQ } from '~/config/landing.config';
</script>

<template>
  <section id="faq" class="landing-section landing-section--muted">
    <div class="landing-container landing-faq">
      <h2 class="landing-title">Часто задаваемые вопросы</h2>
      <Accordion value="0" class="landing-faq__list">
        <AccordionPanel v-for="(item, idx) in LANDING_FAQ" :key="item.question" :value="String(idx)" class="landing-faq__item">
          <AccordionHeader>{{ item.question }}</AccordionHeader>
          <AccordionContent>
            <p class="landing-faq__answer">{{ item.answer }}</p>
          </AccordionContent>
        </AccordionPanel>
      </Accordion>
    </div>
  </section>
</template>
