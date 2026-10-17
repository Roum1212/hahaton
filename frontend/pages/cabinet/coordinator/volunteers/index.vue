<script setup lang="ts">
import {
  COOPERATION_STATUS_LABELS,
  COOPERATION_STATUS_SEVERITY,
  CooperationStatus,
} from '~/core/cooperation/services/cooperation/cooperation.types';

definePageMeta({
  title: 'Волонтёры',
  permissions: [ROLES.coordinator],
});

const { store, loading, load, approve, reject, terminate, createInvitation, revokeInvitation } = useCoordinatorCooperations();
const { name: volunteerName, view: volunteerView, avatar: volunteerAvatar } = useVolunteerNames();
const { copy } = useClipboard('Ссылка скопирована');

const tab = ref<CooperationStatus>(CooperationStatus.PENDING);
const expiresAt = ref<Date | null>(null);
const qrUrl = ref<string | null>(null);
const minDate = new Date();

const tabs = computed(() =>
  [CooperationStatus.PENDING, CooperationStatus.ACTIVE, CooperationStatus.REJECTED, CooperationStatus.TERMINATED].map((value) => ({
    value,
    label: `${COOPERATION_STATUS_LABELS[value]} ${store.asCoordinator.filter((item) => item.status === value).length}`,
  })),
);

const visible = computed(() => store.asCoordinator.filter((item) => item.status === tab.value));

const onCreate = async () => {
  await createInvitation(expiresAt.value);
  expiresAt.value = null;
};

onMounted(load);
</script>

<template>
  <AppPage>
    <template #title>Волонтёры</template>

    <div class="app-columns">
      <div class="app-card">
        <AppFilter v-model="tab" :options="tabs" />

        <div v-if="loading" class="flex flex-col gap-3">
          <Skeleton v-for="n in 3" :key="n" height="64px" border-radius="16px" />
        </div>

        <AppEmpty v-else-if="!visible.length" title="Здесь пока никого нет" image="/images/landing/step-request.jpg" />

        <div v-else class="app-list">
          <div v-for="item in visible" :key="item.volunteerId" class="app-list__item">
            <div class="app-list__lead">
              <NuxtLink :to="`${ROUTES.coordinator.volunteers}/${item.volunteerId}`"
                ><AppAvatar :name="volunteerName(item.volunteerId)" :image="volunteerAvatar(item.volunteerId)"
              /></NuxtLink>
              <div class="app-list__main">
                <span class="app-list__title">
                  <NuxtLink :to="`${ROUTES.coordinator.volunteers}/${item.volunteerId}`">{{ volunteerName(item.volunteerId) }}</NuxtLink>
                  <AppStatus :label="COOPERATION_STATUS_LABELS[item.status]" :severity="COOPERATION_STATUS_SEVERITY[item.status]" />
                </span>
                <div class="app-list__meta">
                  <span v-if="volunteerView(item.volunteerId)?.phone"
                    ><i class="pi pi-phone" />{{ volunteerView(item.volunteerId)?.phone }}</span
                  >
                  <span v-if="volunteerView(item.volunteerId)?.email"
                    ><i class="pi pi-envelope" />{{ volunteerView(item.volunteerId)?.email }}</span
                  >
                  <span><i class="pi pi-clock" />{{ formatDateTime(item.updatedAt) }}</span>
                </div>
              </div>
            </div>
            <div class="app-list__actions">
              <template v-if="item.status === CooperationStatus.PENDING">
                <Button label="Принять" size="small" @click="approve(item.volunteerId)" />
                <Button
                  label="Отклонить"
                  severity="secondary"
                  text
                  size="small"
                  @click="reject(item.volunteerId, volunteerName(item.volunteerId))" />
              </template>
              <Button
                v-if="item.status === CooperationStatus.ACTIVE"
                label="Прекратить"
                severity="danger"
                text
                size="small"
                @click="terminate(item.volunteerId, volunteerName(item.volunteerId))" />
            </div>
          </div>
        </div>
      </div>

      <div class="app-card app-columns__aside">
        <img src="/images/landing/step-handover.jpg" alt="" class="app-promo__image" />
        <span class="app-card__title">Пригласить в команду</span>
        <div class="flex flex-wrap items-end gap-2">
          <label class="app-field grow">
            <span class="app-field__label">Действует до</span>
            <DatePicker
              v-model="expiresAt"
              :min-date="minDate"
              show-time
              hour-format="24"
              date-format="dd.mm.yy"
              placeholder="Бессрочно"
              show-clear
              fluid />
          </label>
          <Button label="Создать ссылку" @click="onCreate" />
        </div>

        <div v-for="invitation in store.invitations" :key="invitation.id" class="app-invitation">
          <div class="flex min-w-0 flex-col">
            <span class="app-invitation__url">{{ invitation.url }}</span>
            <span class="app-card__hint">
              {{ invitation.expiresAt ? `До ${formatDateTime(invitation.expiresAt)}` : 'Бессрочное' }}
            </span>
          </div>
          <div class="flex shrink-0 gap-1">
            <Button v-tooltip.top="'Скопировать'" icon="pi pi-copy" severity="secondary" text size="small" @click="copy(invitation.url)" />
            <Button v-tooltip.top="'QR-код'" icon="pi pi-qrcode" severity="secondary" text size="small" @click="qrUrl = invitation.url" />
            <Button
              v-tooltip.top="'Отозвать'"
              icon="pi pi-trash"
              severity="danger"
              text
              size="small"
              @click="revokeInvitation(invitation.id)" />
          </div>
        </div>
      </div>
    </div>

    <Dialog
      :visible="!!qrUrl"
      modal
      header="QR-код приглашения"
      :draggable="false"
      dismissable-mask
      :style="{ width: '360px' }"
      @update:visible="(value: boolean) => !value && (qrUrl = null)">
      <div class="flex flex-col items-center gap-3">
        <CooperationInvitationQr v-if="qrUrl" :url="qrUrl" />
        <Button label="Скопировать ссылку" icon="pi pi-copy" severity="secondary" outlined size="small" @click="qrUrl && copy(qrUrl)" />
      </div>
    </Dialog>
  </AppPage>
</template>
