export default defineNuxtPlugin(() => {
  const authStore = useAuthStore();
  const volunteerStore = useVolunteerStore();
  const updatesStore = useUpdatesStore();
  const { start, stop } = useUpdatesConnection();

  watch(
    () => authStore.isLoggedIn,
    (loggedIn) => (loggedIn ? start() : stop()),
    { immediate: true },
  );

  watch(
    () => updatesStore.revision,
    () => {
      if (!authStore.isLoggedIn) return;
      handleAsyncOperation(() => Promise.all([authStore.refreshSession(), volunteerStore.fetchProfile()]));
    },
  );
});
