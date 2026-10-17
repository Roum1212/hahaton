<script setup lang="ts">
import type { IReward } from '~/core/reward/services/reward/reward.types';

defineProps<{
  reward: IReward;
  balance: number;
  owner?: string;
}>();

defineEmits<{
  purchase: [];
}>();
</script>

<template>
  <article class="app-reward">
    <div class="app-reward__image" :class="{ 'app-reward__image--photo': reward.imageUrl }">
      <img :src="reward.imageUrl ?? pickIllustration(REWARD_ILLUSTRATIONS, reward.id)" alt="" />
    </div>
    <div class="app-reward__body">
      <span v-if="owner" class="app-reward__owner">{{ owner }}</span>
      <span class="app-reward__title">{{ reward.title }}</span>
      <p v-if="reward.description" class="app-reward__text">{{ reward.description }}</p>
      <div class="app-reward__footer">
        <div class="flex flex-col">
          <span class="app-reward__price">{{ formatCoins(reward.priceCoins) }}</span>
          <span class="app-card__hint">Осталось {{ reward.availableQuantity }}</span>
        </div>
        <Button
          :label="reward.priceCoins > balance ? `Ещё ${reward.priceCoins - balance}` : 'Обменять'"
          size="small"
          :severity="reward.priceCoins > balance ? 'secondary' : undefined"
          :disabled="!reward.availableQuantity || reward.priceCoins > balance"
          @click="$emit('purchase')" />
      </div>
    </div>
  </article>
</template>
