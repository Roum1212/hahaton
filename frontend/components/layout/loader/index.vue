<script setup lang="ts">
const { isShow } = useLoader();
const { isShowProgress } = useLoaderProgress();
</script>
<template>
  <div v-if="isShow" ref="globalLoader" class="loader-container">
    <ProgressSpinner />
  </div>
  <div v-if="isShowProgress" ref="globalLoader" class="loader-container">
    <Card class="flex flex-col gap-2 max-w-[315px] w-full">
      <template #title>
        <span class="font-semibold"> Вносим изменения </span>
      </template>
      <template #content>
        <ProgressBar mode="indeterminate" class="w-full !h-3 my-4" :show-value="true" />
        <div class="text-sm">Пожалуйста, не закрывайте и не обновляйте страницу, пока процесс не завершится.</div>
      </template>
    </Card>
  </div>
</template>

<style scoped>
.loader-container {
  z-index: 9999;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: rgba(0, 0, 0, 0.5);
}
</style>
