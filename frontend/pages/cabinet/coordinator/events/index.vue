<script setup lang="ts">
import type { OwnedEventsFilter } from '~/composables/event/useEvents';

definePageMeta({
  title: 'Мои мероприятия',
  permissions: [ROLES.coordinator],
});

const route = useRoute();
const { store, loading, fetching, load } = useOwnedEvents();
const { open: openCreate } = useEventCreateDialog();

const filter = ref<OwnedEventsFilter>((route.query.filter as OwnedEventsFilter) || 'upcoming');

watch(filter, (value) => load(value), { immediate: true });
useLiveRefresh(() => load(filter.value, 1, true));
</script>

<template>
  <AppPage>
    <template #title>Мои мероприятия</template>
    <template #actions>
      <Button label="Создать мероприятие" icon="pi pi-plus" @click="openCreate" />
    </template>

    <AppFilter v-model="filter" :options="OWNED_EVENTS_FILTERS" />

    <div v-if="loading" class="app-cards">
      <Skeleton v-for="n in 3" :key="n" height="360px" border-radius="20px" />
    </div>

    <div v-else-if="!store.list.length" class="app-card">
      <AppEmpty title="Здесь пока пусто">
        <Button label="Создать мероприятие" icon="pi pi-plus" @click="openCreate" />
      </AppEmpty>
    </div>

    <template v-else>
      <div class="app-cards">
        <EventCard v-for="event in store.list" :key="event.id" :event="event" :to="`${ROUTES.coordinator.events}/${event.id}`">
          <template #actions>
            <NuxtLink v-if="canCompleteEvent(event)" :to="`${ROUTES.coordinator.events}/${event.id}/complete`">
              <Button label="Подвести итоги" size="small" />
            </NuxtLink>
            <NuxtLink v-else :to="`${ROUTES.coordinator.events}/${event.id}`">
              <Button
                :label="event.pendingCount ? `Заявки: ${event.pendingCount}` : 'Открыть'"
                size="small"
                :severity="event.pendingCount ? undefined : 'secondary'"
                :outlined="!event.pendingCount" />
            </NuxtLink>
          </template>
        </EventCard>
      </div>
      <div v-if="!store.pagination.loadAll" class="app-more">
        <Button label="Показать ещё" severity="secondary" outlined :loading="fetching" @click="load(filter, store.pagination.page)" />
      </div>
    </template>
  </AppPage>
</template>
