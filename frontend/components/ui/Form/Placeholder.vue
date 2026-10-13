<script lang="ts" setup>
withDefaults(defineProps<{ rows?: number; count?: number }>(), { rows: 1, count: 1 });
</script>
<template>
  <div :class="`grid gap-2 grid-cols-${rows}`">
    <div v-for="i in count" :key="i" class="flex flex-col gap-2 mb-2">
      <Skeleton height="13px" width="30%" />
      <Skeleton height="37px" width="100%" />
    </div>
    <div v-if="count > 4" class="flex flex-col gap-2">
      <div class="grid grid-cols-2 gap-2">
        <div v-for="i in 4" :key="i" class="flex items-center gap-2">
          <Skeleton height="21px" width="35px" border-radius="30px" />
          <Skeleton height="14px" width="100px" />
        </div>
      </div>
    </div>
  </div>
</template>
