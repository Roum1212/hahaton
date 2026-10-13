import AppPage from '~/layouts/cabinet/AppPage.vue';

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.vueApp.component('AppPage', AppPage);
});
