<script setup lang="ts">
export interface IProfileStat {
  label: string;
  value: string | number;
}

withDefaults(
  defineProps<{
    name: string;
    avatar?: string | null;
    cover: string;
    about?: string | null;
    stats?: IProfileStat[];
    loading?: boolean;
  }>(),
  {
    avatar: null,
    about: null,
    stats: () => [],
    loading: false,
  },
);
</script>

<template>
  <section class="app-profile">
    <div class="app-profile__cover">
      <img :src="cover" alt="" class="app-profile__cover-image" />
      <div class="app-profile__shade" />
      <div class="app-profile__head">
        <AppAvatar :name="name" :image="avatar" size="xl" class="app-profile__avatar" />
        <div class="app-profile__info">
          <Skeleton v-if="loading" width="16rem" height="2rem" />
          <h1 v-else class="app-profile__name">{{ name }}</h1>
          <div v-if="$slots.chips" class="app-profile__chips"><slot name="chips" /></div>
          <div v-if="$slots.meta" class="app-profile__meta"><slot name="meta" /></div>
        </div>
      </div>
    </div>
    <div class="app-profile__bar">
      <dl v-if="stats.length" class="app-profile__stats">
        <div v-for="stat in stats" :key="stat.label" class="app-profile__stat">
          <dt>{{ stat.label }}</dt>
          <dd>{{ stat.value }}</dd>
        </div>
      </dl>
      <div v-if="$slots.actions" class="app-profile__actions"><slot name="actions" /></div>
      <p v-if="about" class="app-profile__about">{{ about }}</p>
    </div>
  </section>
</template>
