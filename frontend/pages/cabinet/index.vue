<script setup lang="ts">
import {
  ParticipationListScope,
  ParticipationStatus,
  type IParticipation,
} from '~/core/participation/services/participation/participation.types';
import { RewardOrderListScope } from '~/core/reward-order/services/reward-order/reward-order.types';

definePageMeta({
  title: 'Мой профиль',
});

const authStore = useAuthStore();
const volunteerStore = useVolunteerStore();
const participationsStore = useParticipationsStore();
const ordersStore = useRewardOrdersStore();
const { store: cooperationsStore, load: loadCooperations } = useVolunteerCooperations();
const { ensure: ensureCoordinators, name: coordinatorName, avatar: coordinatorAvatar } = useCoordinatorNames();
const { ensure: ensureEvents, event } = useEventLookup();
const { label: cityLabel } = useSearchArea();
const { active } = useProfileTabs(VOLUNTEER_TABS);

const loading = ref(true);
const upcoming = ref<IParticipation[]>([]);
const doneCount = ref(0);
const purchasesCount = ref(0);

const profile = computed(() => volunteerStore.profile);
const stats = computed(() => [
  { label: 'Койны', value: profile.value?.coinBalance ?? 0 },
  { label: 'Часы', value: formatHours(profile.value?.confirmedMinutes ?? 0) },
  { label: 'Мероприятия', value: doneCount.value },
  { label: 'Всего заказов', value: purchasesCount.value },
]);
const coordinatorIds = computed(() => cooperationsStore.activeCoordinatorIds);

const loadUpcoming = async () => {
  const approved = await participationsStore.service.find(
    { scope: ParticipationListScope.MY_APPLICATIONS, status: ParticipationStatus.APPROVED },
    { limit: 20, offset: 0 },
  );
  await ensureEvents(approved.list.map((item) => item.eventId));
  upcoming.value = approved.list.filter((item) => !isPast(event(item.eventId)?.endsAt)).slice(0, 4);
};

const load = async (silent = false) => {
  if (!silent) loading.value = true;
  await handleAsyncOperation(
    async () => {
      const [, confirmed, purchases] = await Promise.all([
        volunteerStore.fetchProfile(),
        participationsStore.service.find(
          { scope: ParticipationListScope.MY_APPLICATIONS, status: ParticipationStatus.CONFIRMED },
          { limit: 1, offset: 0 },
        ),
        ordersStore.service.find({ scope: RewardOrderListScope.MY_PURCHASES }, { limit: 1, offset: 0 }),
        loadCooperations(silent),
        loadUpcoming(),
      ]);
      doneCount.value = confirmed.total;
      purchasesCount.value = purchases.total;
      await ensureCoordinators(coordinatorIds.value);
    },
    () => {
      loading.value = false;
    },
  );
};

onMounted(() => load());
useLiveRefresh(() => load(true));

watch(
  () => participationsStore.list.map((item) => `${item.id}:${item.status}`).join(),
  (value, previous) => {
    if (previous && !loading.value) handleAsyncOperation(loadUpcoming);
  },
);
</script>

<template>
  <div class="app-page">
    <AppProfileHero
      :name="profile?.fullName ?? ''"
      :avatar="profile?.avatarUrl"
      :cover="PROFILE_COVER_VOLUNTEER"
      :stats="stats"
      :loading="!profile">
      <template #chips>
        <span class="app-chip app-chip--accent"><i class="pi pi-heart-fill" />Волонтёр</span>
        <span v-if="authStore.hasCoordinatorProfile" class="app-chip app-chip--violet"
          ><i class="pi pi-briefcase" />Координатор</span
        >
      </template>
      <template #meta>
        <span><i class="pi pi-map-marker" />{{ cityLabel }}</span>
        <span v-if="profile?.email"><i class="pi pi-envelope" />{{ profile.email }}</span>
        <span v-if="profile?.phone"><i class="pi pi-phone" />{{ profile.phone }}</span>
      </template>
      <template #actions>
        <NuxtLink :to="ROUTES.volunteer.events">
          <Button label="Найти мероприятие" icon="pi pi-search" />
        </NuxtLink>
      </template>
    </AppProfileHero>

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
            <span class="app-card__title">Ближайшие мероприятия</span>
            <NuxtLink :to="ROUTES.volunteer.events" class="app-link">Каталог</NuxtLink>
          </div>
          <div v-if="loading" class="flex flex-col gap-3">
            <Skeleton v-for="n in 2" :key="n" height="56px" border-radius="14px" />
          </div>
          <span v-else-if="!upcoming.length" class="app-card__hint"
            >Когда координатор примет заявку, мероприятие появится здесь</span
          >
          <div v-else class="app-mini-list">
            <NuxtLink
              v-for="item in upcoming"
              :key="item.id"
              :to="`${ROUTES.volunteer.events}/${item.eventId}`"
              class="app-mini-list__item">
              <img :src="eventImage(event(item.eventId), item.eventId)" alt="" class="app-thumb app-thumb--sm" />
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
          <span v-else-if="!coordinatorIds.length" class="app-card__hint"
            >Подайте заявку координатору, чтобы видеть его мероприятия</span
          >
          <div v-else class="app-mini-list">
            <NuxtLink
              v-for="id in coordinatorIds"
              :key="id"
              :to="`${ROUTES.volunteer.coordinators}/${id}`"
              class="app-mini-list__item">
              <AppAvatar :name="coordinatorName(id)" :image="coordinatorAvatar(id)" size="sm" />
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
