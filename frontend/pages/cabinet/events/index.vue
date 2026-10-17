<script setup lang="ts">
import type { CatalogPeriod } from '~/composables/event/useEvents';

definePageMeta({
  title: 'Мероприятия',
});

const { area, inCityTitle, openDialog } = useSearchArea();
const { store, loading, fetching, load, loadMore } = useEventCatalog();
const { name: coordinatorName } = useCoordinatorNames();

const period = ref<CatalogPeriod>('upcoming');
const periodOptions: { label: string; value: CatalogPeriod }[] = [
  { label: 'Предстоящие', value: 'upcoming' },
  { label: 'Все', value: 'all' },
];

watch([area, period], () => load(area.value, period.value), { immediate: true, deep: true });
useLiveRefresh(() => load(area.value, period.value, 1, undefined, true));
</script>

<template>
  <AppPage>
    <template #title>Мероприятия {{ inCityTitle }}</template>
    <template #actions>
      <Button label="Изменить район" icon="pi pi-map" severity="secondary" outlined @click="openDialog" />
    </template>

    <AppFilter v-model="period" :options="periodOptions" />

    <div v-if="loading" class="app-cards">
      <Skeleton v-for="n in 3" :key="n" height="360px" border-radius="20px" />
    </div>

    <div v-else-if="!store.list.length" class="app-card">
      <AppEmpty title="Ближайших доступных вам мероприятий не найдено">
        <div class="flex flex-wrap justify-center gap-2">
          <Button label="Расширить район" severity="secondary" outlined @click="openDialog" />
          <NuxtLink :to="ROUTES.volunteer.coordinators">
            <Button label="Найти координаторов" />
          </NuxtLink>
        </div>
      </AppEmpty>
    </div>

    <template v-else>
      <div class="app-cards">
        <EventCard
          v-for="event in store.list"
          :key="event.id"
          :event="event"
          :coordinator-name="coordinatorName(event.coordinatorId)"
          :to="`${ROUTES.volunteer.events}/${event.id}`">
          <template #actions>
            <NuxtLink :to="`${ROUTES.volunteer.events}/${event.id}`">
              <Button label="Подробнее" size="small" />
            </NuxtLink>
          </template>
        </EventCard>
      </div>
      <div v-if="!store.pagination.loadAll" class="app-more">
        <Button label="Показать ещё" severity="secondary" outlined :loading="fetching" @click="loadMore(area, period)" />
      </div>
    </template>
  </AppPage>
</template>
