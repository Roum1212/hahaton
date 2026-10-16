<script setup lang="ts">
import {
  PARTICIPATION_CANCELLATION_LABELS,
  PARTICIPATION_STATUS_LABELS,
  PARTICIPATION_STATUS_SEVERITY,
  type IParticipationChange,
} from '~/core/participation/services/participation/participation.types';

defineProps<{
  loading: boolean;
  changes: IParticipationChange[];
}>();

const visible = defineModel<boolean>('visible', { required: true });

const ACTOR_LABELS: Record<IParticipationChange['actor'], string> = {
  volunteer: 'Волонтёр',
  coordinator: 'Координатор',
  system: 'Система',
};
</script>

<template>
  <Dialog v-model:visible="visible" modal header="История заявки" :draggable="false" dismissable-mask :style="{ width: '520px' }">
    <div v-if="loading" class="flex flex-col gap-3">
      <Skeleton v-for="n in 3" :key="n" height="48px" />
    </div>
    <ol v-else class="app-timeline">
      <li v-for="change in changes" :key="change.sequence" class="app-timeline__item">
        <div class="flex flex-wrap items-center gap-2">
          <Tag :value="PARTICIPATION_STATUS_LABELS[change.status]" :severity="PARTICIPATION_STATUS_SEVERITY[change.status]" />
          <span class="app-card__hint">{{ ACTOR_LABELS[change.actor] }} · {{ formatDateTime(change.createdAt) }}</span>
        </div>
        <span v-if="change.cancellationReason" class="app-card__hint">{{
          PARTICIPATION_CANCELLATION_LABELS[change.cancellationReason]
        }}</span>
        <span v-if="change.rejectionReason" class="app-card__hint">{{ change.rejectionReason }}</span>
      </li>
    </ol>
  </Dialog>
</template>
