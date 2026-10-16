<script setup lang="ts">
import type * as YMaps from 'yandex-maps';
import type { IGeoPoint } from '~/core/geo/services/geo/geo.types';

const props = withDefaults(
  defineProps<{
    center: IGeoPoint;
    radiusMeters?: number | null;
    zoom?: number;
    showMarker?: boolean;
    interactive?: boolean;
  }>(),
  {
    radiusMeters: null,
    zoom: 12,
    showMarker: true,
    interactive: true,
  },
);

const emit = defineEmits<{
  'update:center': [value: IGeoPoint];
}>();

const apiKey = String(useRuntimeConfig().public.yandexMapsApiKey || '');
const container = ref<HTMLElement | null>(null);
const locating = ref(false);
const failed = ref(false);

let maps: YandexMaps | null = null;
let map: YMaps.Map | null = null;
let marker: YMaps.Placemark | null = null;
let circle: YMaps.Circle | null = null;

const fitView = () => {
  if (!map) return;
  const bounds = circle?.geometry?.getBounds();
  if (bounds) map.setBounds(bounds, { checkZoomRange: true, zoomMargin: [24, 24, 24, 24] });
  else map.setCenter(toYandexCoords(props.center), props.zoom);
};

const syncCircle = () => {
  if (!map || !maps) return;
  if (!props.radiusMeters) {
    if (circle) map.geoObjects.remove(circle);
    circle = null;
  } else if (circle) {
    circle.geometry?.setCoordinates(toYandexCoords(props.center));
    circle.geometry?.setRadius(props.radiusMeters);
  } else {
    circle = new maps.Circle(
      [toYandexCoords(props.center), props.radiusMeters],
      {},
      { fillColor: `${MAP_ACCENT_COLOR}26`, strokeColor: MAP_ACCENT_COLOR, strokeOpacity: 0.8, strokeWidth: 2, cursor: 'pointer' },
    );
    if (props.interactive) circle.events.add('click', (event) => emitCoords(event.get('coords')));
    map.geoObjects.add(circle);
  }
  fitView();
};

const syncMarker = () => {
  if (!map || !marker) return;
  marker.geometry?.setCoordinates(toYandexCoords(props.center));
  const attached = map.geoObjects.indexOf(marker) !== -1;
  if (props.showMarker && !attached) map.geoObjects.add(marker);
  if (!props.showMarker && attached) map.geoObjects.remove(marker);
};

const emitCoords = (coords: number[]) => {
  const [latitude, longitude] = coords;
  if (latitude === undefined || longitude === undefined) return;
  emit('update:center', { latitude, longitude });
};

const locate = () => {
  if (!navigator.geolocation) return;
  locating.value = true;
  navigator.geolocation.getCurrentPosition(
    ({ coords }) => {
      locating.value = false;
      emit('update:center', { latitude: coords.latitude, longitude: coords.longitude });
    },
    () => {
      locating.value = false;
    },
    { timeout: 10000 },
  );
};

onMounted(async () => {
  try {
    maps = await loadYandexMaps(apiKey);
  } catch {
    failed.value = true;
    return;
  }
  if (!container.value) return;
  map = new maps.Map(
    container.value,
    { center: toYandexCoords(props.center), zoom: props.zoom, controls: ['zoomControl'] },
    { suppressMapOpenBlock: !props.interactive, yandexMapDisablePoiInteractivity: true, autoFitToViewport: 'always' },
  );
  if (!props.interactive) map.behaviors.disable('scrollZoom');
  marker = new maps.Placemark(
    toYandexCoords(props.center),
    {},
    { preset: 'islands#dotIcon', iconColor: MAP_ACCENT_COLOR, draggable: props.interactive },
  );
  marker.events.add('dragend', () => emitCoords(marker!.geometry!.getCoordinates() ?? []));
  if (props.interactive) map.events.add('click', (event) => emitCoords(event.get('coords')));
  syncMarker();
  syncCircle();
});

watch(
  () => [props.center.latitude, props.center.longitude],
  () => {
    syncMarker();
    syncCircle();
  },
);

watch(() => props.radiusMeters, syncCircle);

watch(() => props.showMarker, syncMarker);

onBeforeUnmount(() => {
  map?.destroy();
  map = null;
});

defineExpose({ invalidate: () => map?.container.fitToViewport() });
</script>

<template>
  <div class="geo-map">
    <div ref="container" class="geo-map__canvas" />
    <div v-if="failed" class="geo-map__error">Не удалось загрузить карту</div>
    <button
      v-if="interactive && !failed"
      type="button"
      class="geo-map__locate"
      aria-label="Моё местоположение"
      :disabled="locating"
      @click="locate">
      <i :class="locating ? 'pi pi-spin pi-spinner' : 'pi pi-compass'" />
    </button>
  </div>
</template>
