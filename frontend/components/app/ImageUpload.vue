<script setup lang="ts">
import type { IUploadedImage } from '~/core/image/services/image/image.types';

const props = withDefaults(
  defineProps<{
    image: string | null;
    name?: string;
    shape?: 'circle' | 'cover';
    disabled?: boolean;
  }>(),
  {
    name: '',
    shape: 'circle',
    disabled: false,
  },
);

const emit = defineEmits<{
  uploaded: [image: IUploadedImage];
  remove: [];
}>();

const imagesStore = useImagesStore();
const { toast } = useApp();

const input = ref<HTMLInputElement | null>(null);
const uploading = ref(false);
const preview = ref<string | null>(null);

const shown = computed(() => preview.value ?? props.image);

const pick = () => input.value?.click();

const cropFile = ref<File | null>(null);
const cropVisible = ref(false);
const cropSize = computed(() => (props.shape === 'circle' ? AVATAR_IMAGE_SIZE : COVER_IMAGE_SIZE));

const onChange = (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  target.value = '';
  if (!file) return;
  const error = validateSourceImage(file);
  if (error) {
    toast.error(error);
    return;
  }
  cropFile.value = file;
  cropVisible.value = true;
};

const upload = async (file: File) => {
  const local = URL.createObjectURL(file);
  preview.value = local;
  uploading.value = true;
  const image = await handleAsyncOperation(
    () => imagesStore.upload(file),
    () => {
      uploading.value = false;
    },
  );
  preview.value = null;
  URL.revokeObjectURL(local);
  if (image) emit('uploaded', image);
};
</script>

<template>
  <div class="app-upload" :class="`app-upload--${shape}`">
    <button type="button" class="app-upload__preview" :disabled="disabled || uploading" aria-label="Выбрать изображение" @click="pick">
      <img v-if="shown" :src="shown" alt="" />
      <AppAvatar v-else-if="shape === 'circle'" :name="name" size="lg" />
      <span v-else class="app-upload__placeholder"><i class="pi pi-image" /></span>
      <span class="app-upload__overlay">
        <i :class="uploading ? 'pi pi-spin pi-spinner' : 'pi pi-camera'" />
      </span>
    </button>
    <div class="app-upload__actions">
      <Button
        :label="shown ? 'Заменить' : 'Загрузить'"
        severity="secondary"
        outlined
        size="small"
        :loading="uploading"
        :disabled="disabled"
        @click="pick" />
      <Button
        v-if="image && !uploading"
        label="Удалить"
        severity="secondary"
        text
        size="small"
        :disabled="disabled"
        @click="emit('remove')" />
      <i
        v-tooltip.bottom="`JPEG или PNG. После выбора выделите нужную область, сохраним ${cropSize.width}×${cropSize.height} px`"
        class="pi pi-question-circle text-surface-400" />
    </div>
    <input ref="input" type="file" :accept="IMAGE_ACCEPT" hidden @change="onChange" />
    <AppImageCropDialog v-model:visible="cropVisible" :file="cropFile" :size="cropSize" :round="shape === 'circle'" @crop="upload" />
  </div>
</template>
