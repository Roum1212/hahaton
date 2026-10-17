<script setup lang="ts">
import { CooperationListScope, CooperationStatus } from '~/core/cooperation/services/cooperation/cooperation.types';
import type { IEvent } from '~/core/event/services/event/event.types';
import { ParticipationListScope, ParticipationStatus } from '~/core/participation/services/participation/participation.types';
import { RewardOrderListScope, RewardOrderStatus } from '~/core/reward-order/services/reward-order/reward-order.types';

definePageMeta({
  title: 'Кабинет координатора',
  permissions: [ROLES.coordinator],
});

const eventsStore = useEventsStore();
const { open: openCreate } = useEventCreateDialog();
const participationsStore = useParticipationsStore();
const ordersStore = useRewardOrdersStore();
const cooperationsStore = useCooperationsStore();
const coordinatorsStore = useCoordinatorsStore();
const { active } = useProfileTabs(COORDINATOR_TABS);

const loading = ref(true);
const awaiting = ref<IEvent[]>([]);
const upcoming = ref<IEvent[]>([]);
const pendingApplications = ref(0);
const pendingOrders = ref(0);
const pendingCooperations = ref(0);
const activeVolunteers = ref(0);

const me = computed(() => coordinatorsStore.me);
const attention = computed(() => pendingApplications.value + pendingCooperations.value + pendingOrders.value + awaiting.value.length);

onMounted(async () => {
  await handleAsyncOperation(
    async () => {
      const page = { limit: 5, offset: 0 };
      const single = { limit: 1, offset: 0 };
      const [, waiting, next, applications, orders, cooperations, volunteers] = await Promise.all([
        coordinatorsStore.fetchMe(),
        eventsStore.service.find(toOwnedFilter('awaiting'), page),
        eventsStore.service.find(toOwnedFilter('upcoming'), page),
        participationsStore.service.find({ scope: ParticipationListScope.MY_EVENTS, status: ParticipationStatus.PENDING }, single),
        ordersStore.service.find({ scope: RewardOrderListScope.INCOMING, status: RewardOrderStatus.PENDING }, single),
        cooperationsStore.service.find({ scope: CooperationListScope.AS_COORDINATOR, status: CooperationStatus.PENDING }, single),
        cooperationsStore.service.find({ scope: CooperationListScope.AS_COORDINATOR, status: CooperationStatus.ACTIVE }, single),
      ]);
      awaiting.value = waiting.list;
      upcoming.value = next.list;
      pendingApplications.value = applications.total;
      pendingOrders.value = orders.total;
      pendingCooperations.value = cooperations.total;
      activeVolunteers.value = volunteers.total;
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
      <div class="app-hero__cover app-hero__cover--violet" />
      <div class="app-hero__card">
        <AppAvatar :name="me?.displayName ?? ''" :image="me?.avatarUrl || null" size="lg" class="app-hero__avatar" />
        <div class="app-hero__body">
          <Skeleton v-if="!me" width="18rem" height="2rem" />
          <template v-else>
            <h1 class="app-hero__name">{{ me.displayName }}</h1>
            <div class="app-hero__chips">
              <span class="app-chip app-chip--violet"><i class="pi pi-briefcase" />Координатор</span>
            </div>
            <p v-if="me.description" class="app-hero__about">{{ me.description }}</p>
          </template>
        </div>
        <div class="app-hero__side">
          <div class="app-hero__stats">
            <div class="app-hero__stat">
              <b>{{ activeVolunteers }}</b>
              <span>{{ pluralize(activeVolunteers, ['волонтёр', 'волонтёра', 'волонтёров']) }}</span>
            </div>
            <div class="app-hero__stat">
              <b>{{ upcoming.length }}</b>
              <span>впереди</span>
            </div>
            <div class="app-hero__stat">
              <b>{{ attention }}</b>
              <span>ждут вас</span>
            </div>
          </div>
          <div class="app-hero__actions">
            <Button label="Создать мероприятие" class="flex-1" fluid @click="openCreate" />
            <NuxtLink :to="ROUTES.coordinator.volunteers">
              <Button
                v-tooltip.bottom="'Пригласить волонтёров'"
                icon="pi pi-user-plus"
                severity="secondary"
                outlined
                aria-label="Пригласить волонтёров" />
            </NuxtLink>
          </div>
        </div>
      </div>
    </section>

    <div class="app-columns">
      <div class="app-columns__main">
        <AppTabs v-model="active" :tabs="COORDINATOR_TABS" :counts="{ overview: attention }" />

        <template v-if="active === 'overview'">
          <div v-if="loading" class="app-card">
            <Skeleton v-for="n in 3" :key="n" height="56px" border-radius="16px" />
          </div>

          <template v-else>
            <div class="app-tiles">
              <NuxtLink :to="ROUTES.coordinator.events" class="app-tile">
                <span class="app-tile__value">{{ pendingApplications }}</span>
                <span class="app-tile__label">{{ pluralize(pendingApplications, ['заявка', 'заявки', 'заявок']) }} на участие</span>
              </NuxtLink>
              <NuxtLink :to="ROUTES.coordinator.volunteers" class="app-tile">
                <span class="app-tile__value">{{ pendingCooperations }}</span>
                <span class="app-tile__label">{{ pluralize(pendingCooperations, ['заявка', 'заявки', 'заявок']) }} на сотрудничество</span>
              </NuxtLink>
              <NuxtLink :to="ROUTES.coordinator.orders" class="app-tile">
                <span class="app-tile__value">{{ pendingOrders }}</span>
                <span class="app-tile__label">{{ pluralize(pendingOrders, ['заказ', 'заказа', 'заказов']) }} к выдаче</span>
              </NuxtLink>
            </div>

            <div v-if="awaiting.length" class="app-card">
              <span class="app-card__title">Подведите итоги</span>
              <div class="app-list">
                <div v-for="event in awaiting" :key="event.id" class="app-list__item">
                  <div class="app-list__lead">
                    <img :src="eventImage(event, event.id)" alt="" class="app-thumb" />
                    <div class="app-list__main">
                      <span class="app-list__title">{{ event.title }}</span>
                      <div class="app-list__meta">
                        <span><i class="pi pi-calendar" />{{ formatEventPeriod(event.startsAt, event.endsAt) }}</span>
                        <span
                          ><i class="pi pi-users" />{{ event.approvedCount }}
                          {{ pluralize(event.approvedCount, ['участник', 'участника', 'участников']) }}</span
                        >
                      </div>
                    </div>
                  </div>
                  <NuxtLink :to="`${ROUTES.coordinator.events}/${event.id}/complete`">
                    <Button label="Подвести итоги" size="small" />
                  </NuxtLink>
                </div>
              </div>
            </div>

            <div class="app-card">
              <div class="app-card__head">
                <span class="app-card__title">Ближайшие мероприятия</span>
                <NuxtLink :to="ROUTES.coordinator.events" class="app-link">Все</NuxtLink>
              </div>
              <AppEmpty v-if="!upcoming.length" title="Предстоящих мероприятий нет">
                <Button label="Создать мероприятие" icon="pi pi-plus" @click="openCreate" />
              </AppEmpty>
              <div v-else class="app-list">
                <NuxtLink v-for="event in upcoming" :key="event.id" :to="`${ROUTES.coordinator.events}/${event.id}`" class="app-list__item">
                  <div class="app-list__lead">
                    <img :src="eventImage(event, event.id)" alt="" class="app-thumb" />
                    <div class="app-list__main">
                      <span class="app-list__title">{{ event.title }}</span>
                      <div class="app-list__meta">
                        <span><i class="pi pi-calendar" />{{ formatEventPeriod(event.startsAt, event.endsAt) }}</span>
                        <span><i class="pi pi-users" />{{ event.approvedCount }} из {{ event.volunteerLimit }}</span>
                      </div>
                    </div>
                  </div>
                  <span v-if="event.pendingCount" class="app-chip app-chip--warn">
                    {{ event.pendingCount }} {{ pluralize(event.pendingCount, ['новая заявка', 'новые заявки', 'новых заявок']) }}
                  </span>
                </NuxtLink>
              </div>
            </div>
          </template>
        </template>

        <CoordinatorRewardsTab v-else-if="active === 'rewards'" />
        <CoordinatorAboutTab v-else />
      </div>

      <aside class="app-columns__aside">
        <div class="app-card app-promo">
          <img src="/images/landing/step-handover.jpg" alt="" class="app-promo__image" />
          <span class="app-card__title">Пригласите волонтёров</span>
          <span class="app-card__hint"
            >Отправьте ссылку или покажите QR-код — сотрудничество начнётся сразу после того, как волонтёр примет приглашение</span
          >
          <NuxtLink :to="ROUTES.coordinator.volunteers">
            <Button label="Пригласить" severity="secondary" outlined fluid />
          </NuxtLink>
        </div>
      </aside>
    </div>
  </div>
</template>
