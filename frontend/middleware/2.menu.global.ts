export default defineNuxtRouteMiddleware(() => {
  const authStore = useAuthStore();
  if (!authStore.isLoggedIn) return;

  const { setMenu } = useAppLayout();
  const roles = authStore.roles;
  const isAllowed = (permissions?: string[]) => !permissions?.length || permissions.some((permission) => roles.includes(permission));

  setMenu(
    MENU.filter((section) => isAllowed(section.permissions))
      .map((section) => ({
        ...section,
        groups: section.groups
          .map((group) => ({ ...group, items: group.items.filter((item) => item.visible !== false && isAllowed(item.permissions)) }))
          .filter((group) => group.items.length),
      }))
      .filter((section) => section.groups.length),
  );
});
