<script setup lang="ts">
import { USER_ROLE_LABELS } from '~/core/users/services/users/users.types';

withDefaults(defineProps<{ compact?: boolean }>(), { compact: false });

const authStore = useAuthStore();
const { logout } = useAuthActions();

const name = computed(() => authStore.user?.fullName || authStore.user?.email || 'Пользователь');
const roleLabel = computed(() => (authStore.user ? USER_ROLE_LABELS[authStore.user.role] : ''));
const letter = computed(() => name.value.slice(0, 1).toLocaleUpperCase());

const menu = ref();
const menuId = useId();

const items = [
  {
    label: 'Профиль',
    icon: 'pi pi-user',
    command: () => navigateTo(`${ROUTES.cabinet}/${ROUTES.settings}/${ROUTES.profile}`),
  },
  {
    label: 'Выход',
    icon: 'pi pi-sign-out',
    command: () => logout(),
  },
];

const toggle = (event: Event) => menu.value.toggle(event);
</script>

<template>
  <button
    v-tooltip.right="compact ? name : undefined"
    type="button"
    class="app-user"
    :class="{ 'app-user--compact': compact }"
    aria-haspopup="true"
    :aria-controls="menuId"
    @click="toggle">
    <span class="app-user__avatar">{{ letter }}</span>
    <template v-if="!compact">
      <span class="app-user__info">
        <span class="app-user__name">{{ name }}</span>
        <span class="app-user__role">{{ roleLabel }}</span>
      </span>
      <i class="app-user__caret pi pi-angle-up" />
    </template>
  </button>

  <Menu :id="menuId" ref="menu" :model="items" :popup="true" />
</template>
