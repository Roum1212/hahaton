<script setup lang="ts">
const { city, cities, dialogVisible, select } = useCity();

const query = ref('');
</script>

<template>
  <Dialog
    v-model:visible="dialogVisible"
    modal
    dismissable-mask
    header="Город или регион"
    :draggable="false"
    class="landing-city-dialog"
    :style="{ width: '640px' }">
    <IconField class="landing-city-dialog__search">
      <InputText v-model="query" placeholder="Введите название" fluid />
      <InputIcon class="pi pi-search" />
    </IconField>
    <ul class="landing-city-dialog__list">
      <li v-for="item in cities" :key="item.name">
        <button
          type="button"
          class="landing-city-dialog__city"
          :class="{ 'landing-city-dialog__city--active': item.name === city }"
          @click="select(item.name)">
          {{ item.name }}
        </button>
      </li>
    </ul>
  </Dialog>
</template>
