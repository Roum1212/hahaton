import type { RouteLocationNormalized } from 'vue-router';

export default defineNuxtRouteMiddleware((to: RouteLocationNormalized) => {
  const pagePermissions = (to.meta?.permissions || []) as string[];
  if (!pagePermissions.length) return true;

  const roles = useAuthStore().roles;
  if (pagePermissions.some((permission) => roles.includes(permission))) return true;

  return navigateTo(ROUTES.access);
});
