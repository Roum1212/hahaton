<script setup lang="ts">
definePageMeta({
  title: 'Поощрения',
});

const volunteerStore = useVolunteerStore();
const { store, loading, load } = useRewardCatalog();
const { name: coordinatorName } = useCoordinatorNames();
const { purchase } = useRewardPurchase();

const balance = computed(() => volunteerStore.profile?.coinBalance ?? 0);

onMounted(() => {
  load();
  handleAsyncOperation(() => volunteerStore.fetchProfile());
});
</script>

<template>
  <AppPage>
    <template #title>Поощрения</template>
    <template #actions>
      <NuxtLink :to="{ path: ROUTES.volunteer.home, query: { tab: 'wallet' } }" class="app-balance">
        <i class="pi pi-wallet" />{{ formatCoins(balance) }}
      </NuxtLink>
    </template>

    <div v-if="loading" class="app-cards">
      <Skeleton v-for="n in 3" :key="n" height="340px" border-radius="20px" />
    </div>

    <div v-else-if="!store.list.length" class="app-card">
      <AppEmpty title="Поощрения появятся от ваших координаторов" image="/images/landing/step-reward.jpg">
        <NuxtLink :to="ROUTES.volunteer.coordinators">
          <Button label="Найти координаторов" />
        </NuxtLink>
      </AppEmpty>
    </div>

    <div v-else class="app-cards">
      <RewardCard
        v-for="reward in store.list"
        :key="reward.id"
        :reward="reward"
        :balance="balance"
        :owner="coordinatorName(reward.coordinatorId)"
        @purchase="purchase(reward)" />
    </div>
  </AppPage>
</template>
