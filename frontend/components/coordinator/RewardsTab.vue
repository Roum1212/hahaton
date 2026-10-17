<script setup lang="ts">
import type { IReward, IRewardDto } from '~/core/reward/services/reward/reward.types';

const { store, loading, load, save, remove } = useOwnedRewards();

const editingId = ref<string | null>(null);
const dialogVisible = ref(false);
const form = ref<IRewardDto>({ title: '', description: '', priceCoins: 20, quantityLimit: 10 });
const imagePreview = ref<string | null>(null);
const formElement = ref<HTMLFormElement | null>(null);
const { errors, errorId, validate, reset } = useFormValidation(formElement, () => {
  const errors: Partial<Record<'title' | 'priceCoins' | 'quantityLimit', string>> = {};
  if (!form.value.title.trim()) errors.title = 'Укажите название';
  if (!Number.isInteger(form.value.priceCoins) || form.value.priceCoins < 1)
    errors.priceCoins = 'Укажите целую цену от 1 койна';
  if (!Number.isInteger(form.value.quantityLimit) || form.value.quantityLimit < (editingId.value ? 0 : 1)) {
    errors.quantityLimit = editingId.value
      ? 'Количество не может быть отрицательным'
      : 'Укажите хотя бы одно поощрение';
  }
  return errors;
});

const openCreate = () => {
  reset();
  editingId.value = null;
  form.value = { title: '', description: '', priceCoins: 20, quantityLimit: 10 };
  imagePreview.value = null;
  dialogVisible.value = true;
};

const openEdit = (reward: IReward) => {
  reset();
  editingId.value = reward.id;
  form.value = {
    title: reward.title,
    description: reward.description,
    priceCoins: reward.priceCoins,
    quantityLimit: reward.quantityLimit,
  };
  imagePreview.value = reward.imageUrl;
  dialogVisible.value = true;
};

const setImage = (path: string, url: string | null) => {
  form.value.imagePath = path;
  imagePreview.value = url;
};

const onSave = async () => {
  if (!(await validate())) return;
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

    <AppEmpty
      v-else-if="!store.list.length"
      title="Добавьте первое поощрение для волонтёров"
      image="/images/landing/step-reward.jpg" />

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
              <Button
                icon="pi pi-pencil"
                severity="secondary"
                text
                rounded
                size="small"
                aria-label="Изменить"
                @click="openEdit(reward)" />
              <Button
                icon="pi pi-trash"
                severity="danger"
                text
                rounded
                size="small"
                aria-label="Удалить"
                @click="remove(reward)" />
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
      <form ref="formElement" class="app-form" @submit.prevent="onSave">
        <div class="app-field">
          <span class="app-field__label">Изображение</span>
          <AppImageUpload
            :image="imagePreview"
            shape="cover"
            @uploaded="setImage($event.path, $event.url)"
            @remove="setImage('', null)" />
        </div>
        <label class="app-field" :data-invalid="!!errors.title">
          <span class="app-field__label">Название</span>
          <InputText
            v-model="form.title"
            :invalid="!!errors.title"
            :aria-describedby="errors.title ? errorId('title') : undefined"
            fluid />
          <small v-if="errors.title" :id="errorId('title')" class="app-field__error">{{ errors.title }}</small>
        </label>
        <label class="app-field">
          <span class="app-field__label">Описание</span>
          <Textarea v-model="form.description" rows="3" auto-resize fluid />
        </label>
        <div class="app-form-grid">
          <label class="app-field" :data-invalid="!!errors.priceCoins">
            <span class="app-field__label">Цена, койнов</span>
            <InputNumber
              v-model="form.priceCoins"
              :invalid="!!errors.priceCoins"
              :aria-describedby="errors.priceCoins ? errorId('priceCoins') : undefined"
              :min="1"
              fluid
              :pt="{
                pcInputText: { root: { 'aria-describedby': errors.priceCoins ? errorId('priceCoins') : undefined } },
              }" />
            <small v-if="errors.priceCoins" :id="errorId('priceCoins')" class="app-field__error">{{
              errors.priceCoins
            }}</small>
          </label>
          <label class="app-field" :data-invalid="!!errors.quantityLimit">
            <span class="app-field__label">Количество</span>
            <InputNumber
              v-model="form.quantityLimit"
              :invalid="!!errors.quantityLimit"
              :aria-describedby="errors.quantityLimit ? errorId('quantityLimit') : undefined"
              :min="0"
              fluid
              :pt="{
                pcInputText: {
                  root: { 'aria-describedby': errors.quantityLimit ? errorId('quantityLimit') : undefined },
                },
              }" />
            <small v-if="errors.quantityLimit" :id="errorId('quantityLimit')" class="app-field__error">{{
              errors.quantityLimit
            }}</small>
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
