<script setup lang="ts">
import type { MenuItem } from 'primevue/menuitem';
import { CABINET_MODE_LABELS } from '~/config/menu.config';

const volunteerStore = useVolunteerStore();
const { logout } = useAuthActions();
const { mode, canCoordinate, switchMode, homeFor } = useCabinetMode();

const name = computed(() => volunteerStore.profile?.fullName || 'Участник');

const menu = ref();
const menuId = useId();

const items = computed<MenuItem[]>(() => {
  const result: MenuItem[] = [
    { label: name.value, items: [] },
    { label: 'Мой профиль', icon: 'pi pi-user', command: () => navigateTo(homeFor('volunteer')) },
  ];
  if (canCoordinate.value) {
    const target = mode.value === 'coordinator' ? 'volunteer' : 'coordinator';
    result.push({
      label: `Режим: ${CABINET_MODE_LABELS[target].toLocaleLowerCase()}`,
      icon: target === 'coordinator' ? 'pi pi-briefcase' : 'pi pi-heart',
      command: () => switchMode(target),
    });
  }
  result.push(
    { label: 'Анкета', icon: 'pi pi-pencil', command: () => navigateTo({ path: homeFor('volunteer'), query: { tab: 'about' } }) },
    { separator: true },
    { label: 'Выйти', icon: 'pi pi-sign-out', command: () => logout() },
  );
  return result;
});

const toggle = (event: Event) => menu.value.toggle(event);
</script>

<template>
  <button type="button" class="app-user-button" aria-haspopup="true" :aria-controls="menuId" @click="toggle">
    <AppAvatar :name="name" :image="volunteerStore.profile?.avatarUrl" size="sm" />
    <span class="app-user-button__name">{{ name }}</span>
    <i class="pi pi-angle-down" />
  </button>

  <Menu :id="menuId" ref="menu" :model="items" :popup="true" class="app-user-menu" />
</template>
