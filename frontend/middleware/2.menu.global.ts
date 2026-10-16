export default defineNuxtRouteMiddleware((to) => {
  if (!useAuthStore().isLoggedIn) return;
  const pathMode = modeFromPath(to.path);
  if (pathMode) useCabinetMode().setMode(pathMode);
});
