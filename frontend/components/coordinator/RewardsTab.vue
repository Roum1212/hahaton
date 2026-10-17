<script setup lang="ts">
import type { IReward, IRewardDto } from '~/core/reward/services/reward/reward.types';

const { store, loading, load, save, remove } = useOwnedRewards();
const { toast } = useApp();

const editingId = ref<string | null>(null);
const dialogVisible = ref(false);
const form = ref<IRewardDto>({ title: '', description: '', priceCoins: 20, quantityLimit: 10 });
const imagePreview = ref<string | null>(null);

const openCreate = () => {
  editingId.value = null;
  form.value = { title: '', description: '', priceCoins: 20, quantityLimit: 10 };
  imagePreview.value = null;
  dialogVisible.value = true;
};

const openEdit = (reward: IReward) => {
  editingId.value = reward.id;
  form.value = { title: reward.title, description: reward.description, priceCoins: reward.priceCoins, quantityLimit: reward.quantityLimit };
  imagePreview.value = reward.imageUrl;
  dialogVisible.value = true;
};

const setImage = (path: string, url: string | null) => {
  form.value.imagePath = path;
  imagePreview.value = url;
};

const onSave = async () => {
  if (!form.value.title.trim()) {
    toast.error('Укажите название');
    return;
  }
  if (form.value.priceCoins <= 0 || form.value.quantityLimit < (editingId.value ? 0 : 1)) {
    toast.error('Цена и количество должны быть больше нуля');
    return;
  }
  if (await save(editingId.value, { ...form.value, title: form.value.title.trim() })) dialogVisible.value = false;
};

onMounted(load);
</script>

<template>
  <div class="app-card">
    <div class="app-card__head">
      <span class="app-card__title">Что можно получить за койны</span>
      <Button label="Добавить" icon="pi pi-plus" size="small" @click="openCreate" />
    </div>

    <div v-if="loading" class="app-grid">
      <Skeleton v-for="n in 2" :key="n" height="200px" border-radius="20px" />
    </div>

    <AppEmpty v-else-if="!store.list.length" title="Добавьте первое поощрение для волонтёров" image="/images/landing/step-reward.jpg" />

    <div v-else class="app-grid">
      <article v-for="reward in store.list" :key="reward.id" class="app-reward app-reward--bordered">
        <div class="app-reward__image" :class="{ 'app-reward__image--photo': reward.imageUrl }">
          <img :src="reward.imageUrl ?? pickIllustration(REWARD_ILLUSTRATIONS, reward.id)" alt="" />
        </div>
        <div class="app-reward__body">
          <span class="app-reward__title">{{ reward.title }}</span>
          <p v-if="reward.description" class="app-reward__text">{{ reward.description }}</p>
          <div class="app-reward__footer">
            <div class="flex flex-col">
              <span class="app-reward__price">{{ formatCoins(reward.priceCoins) }}</span>
              <span class="app-card__hint">Осталось {{ reward.availableQuantity }} из {{ reward.quantityLimit }}</span>
            </div>
            <div class="flex gap-1">
              <Button icon="pi pi-pencil" severity="secondary" text rounded size="small" aria-label="Изменить" @click="openEdit(reward)" />
              <Button icon="pi pi-trash" severity="danger" text rounded size="small" aria-label="Удалить" @click="remove(reward)" />
            </div>
          </div>
        </div>
      </article>
    </div>

    <Dialog
      v-model:visible="dialogVisible"
      modal
      :header="editingId ? 'Изменение поощрения' : 'Новое поощрение'"
      :draggable="false"
      :style="{ width: '480px' }">
      <form class="app-form" @submit.prevent="onSave">
        <div class="app-field">
          <span class="app-field__label">Изображение</span>
          <AppImageUpload :image="imagePreview" shape="cover" @uploaded="setImage($event.path, $event.url)" @remove="setImage('', null)" />
        </div>
        <label class="app-field">
          <span class="app-field__label">Название</span>
          <InputText v-model="form.title" fluid />
        </label>
        <label class="app-field">
          <span class="app-field__label">Описание</span>
          <Textarea v-model="form.description" rows="3" auto-resize fluid />
        </label>
        <div class="app-form-grid">
          <label class="app-field">
            <span class="app-field__label">Цена, койнов</span>
            <InputNumber v-model="form.priceCoins" :min="1" fluid />
          </label>
          <label class="app-field">
            <span class="app-field__label">Количество</span>
            <InputNumber v-model="form.quantityLimit" :min="0" fluid />
          </label>
        </div>
        <div class="app-form__actions">
          <Button label="Отмена" severity="secondary" text @click="dialogVisible = false" />
          <Button type="submit" label="Сохранить" />
        </div>
      </form>
    </Dialog>
  </div>
</template>
