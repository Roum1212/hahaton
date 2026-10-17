<script setup lang="ts">
import type { IImageSize } from '~/config/image.config';

type Corner = 'nw' | 'ne' | 'sw' | 'se';

interface IRect {
  x: number;
  y: number;
  width: number;
  height: number;
}

interface IDrag {
  mode: 'move' | Corner;
  startX: number;
  startY: number;
  start: IRect;
}

const props = withDefaults(
  defineProps<{
    file: File | null;
    size: IImageSize;
    round?: boolean;
  }>(),
  {
    round: false,
  },
);

const emit = defineEmits<{
  crop: [file: File];
}>();

const visible = defineModel<boolean>('visible', { required: true });
const { toast } = useApp();

const CORNERS: Corner[] = ['nw', 'ne', 'sw', 'se'];

const source = ref<string | null>(null);
const image = ref<HTMLImageElement | null>(null);
const natural = ref({ width: 0, height: 0 });
const scale = ref(1);
const selection = ref<IRect>({ x: 0, y: 0, width: 0, height: 0 });
const saving = ref(false);

let drag: IDrag | null = null;

const aspect = computed(() => props.size.width / props.size.height);

const frameStyle = computed(() => ({
  left: `${selection.value.x * scale.value}px`,
  top: `${selection.value.y * scale.value}px`,
  width: `${selection.value.width * scale.value}px`,
  height: `${selection.value.height * scale.value}px`,
}));

const clamp = (value: number, min: number, max: number) => Math.min(Math.max(value, min), max);

const updateScale = () => {
  if (image.value && natural.value.width) scale.value = image.value.clientWidth / natural.value.width;
};

const onLoad = () => {
  if (!image.value) return;
  const { naturalWidth, naturalHeight } = image.value;
  natural.value = { width: naturalWidth, height: naturalHeight };
  const width = Math.min(naturalWidth, naturalHeight * aspect.value);
  const height = width / aspect.value;
  selection.value = { x: (naturalWidth - width) / 2, y: (naturalHeight - height) / 2, width, height };
  updateScale();
};

const onError = () => {
  toast.error('Не удалось прочитать изображение, выберите другой файл');
  visible.value = false;
};

const moveRect = (start: IRect, dx: number, dy: number): IRect => ({
  ...start,
  x: clamp(start.x + dx, 0, natural.value.width - start.width),
  y: clamp(start.y + dy, 0, natural.value.height - start.height),
});

const resizeRect = (start: IRect, corner: Corner, dx: number, dy: number): IRect => {
  const east = corner.endsWith('e');
  const south = corner.startsWith('s');
  const anchorX = east ? start.x : start.x + start.width;
  const anchorY = south ? start.y : start.y + start.height;
  const maxWidth = Math.min(
    east ? natural.value.width - anchorX : anchorX,
    (south ? natural.value.height - anchorY : anchorY) * aspect.value,
  );
  const wanted = Math.max(start.width + (east ? dx : -dx), (start.height + (south ? dy : -dy)) * aspect.value);
  const width = clamp(wanted, Math.min(MIN_CROP_SIDE, maxWidth), maxWidth);
  const height = width / aspect.value;
  return { x: east ? anchorX : anchorX - width, y: south ? anchorY : anchorY - height, width, height };
};

const onPointerMove = (event: PointerEvent) => {
  if (!drag) return;
  const dx = (event.clientX - drag.startX) / scale.value;
  const dy = (event.clientY - drag.startY) / scale.value;
  selection.value = drag.mode === 'move' ? moveRect(drag.start, dx, dy) : resizeRect(drag.start, drag.mode, dx, dy);
};

const stopDrag = () => {
  drag = null;
  window.removeEventListener('pointermove', onPointerMove);
  window.removeEventListener('pointerup', stopDrag);
};

const startDrag = (event: PointerEvent, mode: IDrag['mode']) => {
  event.preventDefault();
  event.stopPropagation();
  drag = { mode, startX: event.clientX, startY: event.clientY, start: { ...selection.value } };
  window.addEventListener('pointermove', onPointerMove);
  window.addEventListener('pointerup', stopDrag);
};

const save = async () => {
  if (!image.value) return;
  saving.value = true;
  try {
    const canvas = document.createElement('canvas');
    canvas.width = props.size.width;
    canvas.height = props.size.height;
    const context = canvas.getContext('2d');
    if (!context) return;
    context.fillStyle = '#ffffff';
    context.fillRect(0, 0, canvas.width, canvas.height);
    context.imageSmoothingQuality = 'high';
    const { x, y, width, height } = selection.value;
    context.drawImage(image.value, x, y, width, height, 0, 0, canvas.width, canvas.height);
    const blob = await new Promise<Blob | null>((resolve) => canvas.toBlob(resolve, 'image/jpeg', CROPPED_IMAGE_QUALITY));
    if (!blob) return;
    emit('crop', new File([blob], 'image.jpg', { type: 'image/jpeg' }));
    visible.value = false;
  } finally {
    saving.value = false;
  }
};

watch(
  () => [visible.value, props.file] as const,
  ([open, file]) => {
    if (source.value) URL.revokeObjectURL(source.value);
    natural.value = { width: 0, height: 0 };
    source.value = open && file ? URL.createObjectURL(file) : null;
  },
  { immediate: true },
);

useResizeObserver(image, updateScale);

onBeforeUnmount(() => {
  stopDrag();
  if (source.value) URL.revokeObjectURL(source.value);
});
</script>

<template>
  <Dialog
    v-model:visible="visible"
    modal
    header="Выбор миниатюры"
    :draggable="false"
    :closable="!saving"
    class="app-crop-dialog"
    :style="{ width: '760px' }"
    :breakpoints="{ '767px': '100vw' }">
    <div class="app-crop">
      <div class="app-crop__stage">
        <img v-if="source" ref="image" :src="source" alt="" class="app-crop__image" draggable="false" @load="onLoad" @error="onError" />
        <div
          v-if="natural.width"
          class="app-crop__frame"
          :class="{ 'app-crop__frame--round': round }"
          :style="frameStyle"
          @pointerdown="startDrag($event, 'move')">
          <span class="app-crop__mask" />
          <span
            v-for="corner in CORNERS"
            :key="corner"
            class="app-crop__handle"
            :class="`app-crop__handle--${corner}`"
            @pointerdown="startDrag($event, corner)" />
        </div>
      </div>
    </div>
    <template #footer>
      <Button label="Отмена" severity="secondary" text :disabled="saving" @click="visible = false" />
      <Button label="Сохранить" :loading="saving" :disabled="!natural.width" @click="save" />
    </template>
  </Dialog>
</template>
