<script setup lang="ts">
import { EVENT_REQUIREMENT_LABELS, EventCandidateRequirement } from '~/core/event/services/event/event.types';
import type { IEventForm } from '~/composables/event/useEvents';
import type { IAddressSuggestion, IGeoPoint } from '~/core/geo/services/geo/geo.types';

const form = defineModel<IEventForm>({ required: true });

withDefaults(
  defineProps<{
    disabled?: boolean;
    errors?: Partial<Record<keyof IEventForm, string>>;
  }>(),
  { errors: () => ({}) },
);

const formId = useId();
const errorId = (field: keyof IEventForm) => `${formId}-${field}-error`;

const geoStore = useGeoStore();
const { city } = useSearchArea();
const { toast } = useApp();

const addressQuery = ref<IAddressSuggestion | string>(form.value.location);
const {
  suggestions,
  searching,
  searched,
  emptyMessage,
  search,
  reset: resetSearch,
} = useSuggestionSearch((query) => geoStore.service.suggest(query, 'city', 'house', 7), 'Начните вводить адрес');
const minDate = new Date();

const requirementOptions = Object.values(EventCandidateRequirement)
  .filter(
    (value): value is EventCandidateRequirement =>
      typeof value === 'number' && value !== EventCandidateRequirement.UNSPECIFIED,
  )
  .map((value) => ({ value, label: EVENT_REQUIREMENT_LABELS[value] }));

const mapCenter = computed<IGeoPoint>(() => form.value.point ?? city.value.center ?? DEFAULT_CITY.center!);

const selectAddress = async ({ value }: { value: IAddressSuggestion }) => {
  let point = value.point;
  if (!point && value.fiasId) {
    const [precise] =
      (await handleAsyncOperation(() => geoStore.service.suggest(value.fiasId!, 'city', 'house', 1))) ?? [];
    point = precise?.point ?? null;
  }
  form.value.location = value.value;
  addressQuery.value = value.value;
  if (point) form.value.point = point;
  if (!point || (value.qcGeo ?? 0) > 1) toast.warning('Уточните точку на карте: координаты адреса приблизительные');
};

const onAddressInput = (value: IAddressSuggestion | string) => {
  if (typeof value === 'string') {
    form.value.location = value;
    if (!value.trim()) resetSearch();
  }
};

watch(
  () => form.value.location,
  (location) => {
    const current = typeof addressQuery.value === 'string' ? addressQuery.value : addressQuery.value?.value;
    if (location !== current) addressQuery.value = location;
  },
);

const { resolving, resolve } = useReverseGeocode();

const setImage = (path: string, url: string | null) => {
  form.value.imagePath = path;
  form.value.imageUrl = url;
};

const setPoint = async (point: IGeoPoint) => {
  form.value.point = point;
  const result = await resolve(point);
  if (!result?.address) return;
  form.value.location = result.address;
  addressQuery.value = result.address;
};
</script>

<template>
  <div class="app-form">
    <label class="app-field" :data-invalid="!!errors.title">
      <span class="app-field__label">Название</span>
      <InputText
        v-model="form.title"
        :disabled="disabled"
        fluid
        :invalid="!!errors.title"
        :aria-describedby="errors.title ? errorId('title') : undefined" />
      <small v-if="errors.title" :id="errorId('title')" class="app-field__error">{{ errors.title }}</small>
    </label>

    <div class="app-field">
      <span class="app-field__label">Фото</span>
      <AppImageUpload
        :image="form.imageUrl"
        shape="cover"
        :disabled="disabled"
        @uploaded="setImage($event.path, $event.url)"
        @remove="setImage('', null)" />
    </div>

    <label class="app-field">
      <span class="app-field__label">Описание</span>
      <Textarea v-model="form.description" :disabled="disabled" rows="4" auto-resize fluid />
    </label>

    <div class="app-field" :data-invalid="!!errors.point">
      <span class="app-field__label">Место проведения</span>
      <AutoComplete
        v-model="addressQuery"
        :suggestions="suggestions"
        option-label="value"
        :placeholder="resolving ? 'Определяем адрес…' : 'Улица и номер дома или точка на карте'"
        :delay="300"
        :loading="resolving || searching"
        :show-empty-message="searched && !searching"
        :empty-search-message="emptyMessage"
        :disabled="disabled"
        fluid
        :invalid="!!errors.point"
        :aria-describedby="errors.point ? errorId('point') : undefined"
        :pt="{ pcInputText: { root: { 'aria-describedby': errors.point ? errorId('point') : undefined } } }"
        @complete="search"
        @option-select="selectAddress"
        @update:model-value="onAddressInput" />
      <small v-if="errors.point" :id="errorId('point')" class="app-field__error">{{ errors.point }}</small>
      <GeoMap
        :center="mapCenter"
        :zoom="form.point ? 15 : 12"
        :show-marker="!!form.point"
        class="app-form__map"
        @update:center="setPoint" />
      <span v-if="form.point" class="app-card__hint">
        {{ form.point.latitude.toFixed(5) }}, {{ form.point.longitude.toFixed(5) }}
      </span>
    </div>

    <div class="app-form-grid">
      <label class="app-field" :data-invalid="!!errors.startsAt">
        <span class="app-field__label">Начало</span>
        <DatePicker
          v-model="form.startsAt"
          :min-date="minDate"
          show-time
          hour-format="24"
          :step-minute="30"
          date-format="dd.mm.yy"
          :disabled="disabled"
          show-icon
          fluid
          :invalid="!!errors.startsAt"
          :aria-describedby="errors.startsAt ? errorId('startsAt') : undefined"
          :pt="{ pcInputText: { root: { 'aria-describedby': errors.startsAt ? errorId('startsAt') : undefined } } }" />
        <small v-if="errors.startsAt" :id="errorId('startsAt')" class="app-field__error">{{ errors.startsAt }}</small>
      </label>
      <div class="app-field" :data-invalid="!!errors.durationMinutes">
        <span class="app-field__label">Длительность</span>
        <EventDurationInput
          v-model="form.durationMinutes"
          :disabled="disabled"
          :invalid="!!errors.durationMinutes"
          :aria-describedby="errors.durationMinutes ? errorId('durationMinutes') : undefined" />
        <small v-if="errors.durationMinutes" :id="errorId('durationMinutes')" class="app-field__error">{{
          errors.durationMinutes
        }}</small>
      </div>
      <label class="app-field" :data-invalid="!!errors.volunteerLimit">
        <span class="app-field__label">Мест для волонтёров</span>
        <InputNumber
          v-model="form.volunteerLimit"
          :min="1"
          :max="1000"
          show-buttons
          :disabled="disabled"
          fluid
          :invalid="!!errors.volunteerLimit"
          :aria-describedby="errors.volunteerLimit ? errorId('volunteerLimit') : undefined"
          :pt="{
            pcInputText: {
              root: { 'aria-describedby': errors.volunteerLimit ? errorId('volunteerLimit') : undefined },
            },
          }" />
        <small v-if="errors.volunteerLimit" :id="errorId('volunteerLimit')" class="app-field__error">{{
          errors.volunteerLimit
        }}</small>
      </label>
      <label class="app-field" :data-invalid="!!errors.requirements">
        <span class="app-field__label">Требования к кандидатам</span>
        <MultiSelect
          v-model="form.requirements"
          :options="requirementOptions"
          option-label="label"
          option-value="value"
          display="chip"
          placeholder="Без требований"
          :disabled="disabled"
          fluid
          :invalid="!!errors.requirements"
          :aria-describedby="errors.requirements ? errorId('requirements') : undefined" />
        <small v-if="errors.requirements" :id="errorId('requirements')" class="app-field__error">{{
          errors.requirements
        }}</small>
      </label>
    </div>
  </div>
</template>
