<script setup lang="ts">
import {
  ParticipationListScope,
  ParticipationStatus,
  type IParticipation,
} from '~/core/participation/services/participation/participation.types';

definePageMeta({
  title: 'Мой профиль',
});

const authStore = useAuthStore();
const volunteerStore = useVolunteerStore();
const participationsStore = useParticipationsStore();
const { store: cooperationsStore, load: loadCooperations } = useVolunteerCooperations();
const { ensure: ensureCoordinators, name: coordinatorName } = useCoordinatorNames();
const { ensure: ensureEvents, event } = useEventLookup();
const { label: cityLabel } = useSearchArea();
const { active } = useProfileTabs(VOLUNTEER_TABS);

const loading = ref(true);
const upcoming = ref<IParticipation[]>([]);
const doneCount = ref(0);

const profile = computed(() => volunteerStore.profile);
const hours = computed(() => Math.floor((profile.value?.confirmedMinutes ?? 0) / 60));
const coordinatorIds = computed(() => cooperationsStore.activeCoordinatorIds);

onMounted(async () => {
  loading.value = true;
  await handleAsyncOperation(
    async () => {
      const [, approved, confirmed] = await Promise.all([
        volunteerStore.fetchProfile(),
        participationsStore.service.find(
          { scope: ParticipationListScope.MY_APPLICATIONS, status: ParticipationStatus.APPROVED },
          { limit: 20, offset: 0 },
        ),
        participationsStore.service.find(
          { scope: ParticipationListScope.MY_APPLICATIONS, status: ParticipationStatus.CONFIRMED },
          { limit: 1, offset: 0 },
        ),
        loadCooperations(),
      ]);
      doneCount.value = confirmed.total;
      await Promise.all([ensureEvents(approved.list.map((item) => item.eventId)), ensureCoordinators(coordinatorIds.value)]);
      upcoming.value = approved.list.filter((item) => !isPast(event(item.eventId)?.endsAt)).slice(0, 4);
    },
    () => {
      loading.value = false;
    },
  );
});
</script>

<template>
  <div class="app-page">
    <section class="app-hero">
      <div class="app-hero__cover" />
      <div class="app-hero__card">
        <AppAvatar :name="profile?.fullName ?? ''" size="lg" class="app-hero__avatar" />
        <div class="app-hero__body">
          <Skeleton v-if="!profile" width="18rem" height="2rem" />
          <template v-else>
            <h1 class="app-hero__name">{{ profile.fullName }}</h1>
            <div class="app-hero__chips">
              <span class="app-chip app-chip--accent"><i class="pi pi-heart-fill" />Волонтёр</span>
              <span v-if="authStore.hasCoordinatorProfile" class="app-chip app-chip--violet"><i class="pi pi-briefcase" />Координатор</span>
              <span class="app-chip"><i class="pi pi-map-marker" />{{ cityLabel }}</span>
            </div>
            <div class="app-hero__meta">
              <span v-if="profile.email"><i class="pi pi-envelope" />{{ profile.email }}</span>
              <span v-if="profile.phone"><i class="pi pi-phone" />{{ profile.phone }}</span>
            </div>
          </template>
        </div>
        <div class="app-hero__side">
          <div class="app-hero__stats">
            <div class="app-hero__stat">
              <b>{{ profile?.coinBalance ?? 0 }}</b>
              <span>{{ pluralize(profile?.coinBalance ?? 0, ['коин', 'коина', 'коинов']) }}</span>
            </div>
            <div class="app-hero__stat">
              <b>{{ hours }}</b>
              <span>{{ pluralize(hours, ['час', 'часа', 'часов']) }}</span>
            </div>
            <div class="app-hero__stat">
              <b>{{ doneCount }}</b>
              <span>{{ pluralize(doneCount, ['мероприятие', 'мероприятия', 'мероприятий']) }}</span>
            </div>
          </div>
          <div class="app-hero__actions">
            <NuxtLink :to="ROUTES.volunteer.events" class="flex-1">
              <Button label="Найти мероприятие" fluid />
            </NuxtLink>
            <NuxtLink :to="ROUTES.volunteer.rewards">
              <Button v-tooltip.bottom="'Поощрения'" icon="pi pi-gift" severity="secondary" outlined aria-label="Поощрения" />
            </NuxtLink>
          </div>
        </div>
      </div>
    </section>

    <div class="app-columns">
      <div class="app-columns__main">
        <AppTabs v-model="active" :tabs="VOLUNTEER_TABS" />
        <VolunteerApplicationsTab v-if="active === 'applications'" />
        <VolunteerPurchasesTab v-else-if="active === 'purchases'" />
        <VolunteerWalletTab v-else-if="active === 'wallet'" />
        <VolunteerAboutTab v-else />
      </div>

      <aside class="app-columns__aside">
        <div class="app-card">
          <div class="app-card__head">
            <span class="app-card__title">Ближайшие дела</span>
            <NuxtLink :to="ROUTES.volunteer.events" class="app-link">Каталог</NuxtLink>
          </div>
          <div v-if="loading" class="flex flex-col gap-3">
            <Skeleton v-for="n in 2" :key="n" height="56px" border-radius="14px" />
          </div>
          <span v-else-if="!upcoming.length" class="app-card__hint">Когда координатор примет заявку, дело появится здесь</span>
          <div v-else class="app-mini-list">
            <NuxtLink
              v-for="item in upcoming"
              :key="item.id"
              :to="`${ROUTES.volunteer.events}/${item.eventId}`"
              class="app-mini-list__item">
              <img :src="pickIllustration(EVENT_ILLUSTRATIONS, item.eventId)" alt="" class="app-thumb app-thumb--sm" />
              <span class="app-mini-list__text">
                <b>{{ event(item.eventId)?.title }}</b>
                <span>{{ formatEventPeriod(event(item.eventId)?.startsAt ?? null, null) }}</span>
              </span>
            </NuxtLink>
          </div>
        </div>

        <div class="app-card">
          <div class="app-card__head">
            <span class="app-card__title">Мои координаторы</span>
            <NuxtLink :to="ROUTES.volunteer.coordinators" class="app-link">Все</NuxtLink>
          </div>
          <div v-if="loading" class="flex flex-col gap-3">
            <Skeleton v-for="n in 2" :key="n" height="44px" border-radius="14px" />
          </div>
          <span v-else-if="!coordinatorIds.length" class="app-card__hint">Подайте заявку координатору, чтобы видеть его мероприятия</span>
          <div v-else class="app-mini-list">
            <NuxtLink v-for="id in coordinatorIds" :key="id" :to="`${ROUTES.volunteer.coordinators}/${id}`" class="app-mini-list__item">
              <AppAvatar :name="coordinatorName(id)" size="sm" />
              <span class="app-mini-list__text">
                <b>{{ coordinatorName(id) }}</b>
              </span>
            </NuxtLink>
          </div>
        </div>
      </aside>
    </div>
  </div>
</template>
