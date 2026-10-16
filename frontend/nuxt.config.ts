import tailwindcss from '@tailwindcss/vite';
import { fileURLToPath } from 'node:url';
import ThemeSettings from './theme.settings';

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  ssr: false,

  alias: {
    '#genproto': fileURLToPath(new URL('./genproto', import.meta.url)),
  },

  components: {
    dirs: ['~/components'],
  },

  devtools: {
    enabled: true,
  },

  modules: ['@vueuse/nuxt', '@nuxt/eslint', '@primevue/nuxt-module', '@pinia/nuxt'],

  css: ['primeicons/primeicons.css', '@/assets/style.scss', '@/assets/tailwind.css'],

  runtimeConfig: {
    public: {
      rpcUrl: process.env.APP_RPC_URL || 'http://localhost:8080',
      imagesPath: process.env.APP_IMAGES_URL || '',
      siteUrl: process.env.APP_PUBLIC_URL || 'http://localhost:3000',
      siteName: process.env.APP_PUBLIC_SITE_NAME || 'Template',
      siteDescription: process.env.APP_PUBLIC_SITE_DESCRIPTION || '',
      titleSeparator: '|',
      demoLogin: process.env.APP_DEMO_LOGIN === 'true',
      yandexMapsApiKey: process.env.APP_YANDEX_MAPS_API_KEY || '',
      flida: {
        clientId: process.env.APP_FLIDA_CLIENT_ID || '',
        redirectUri: process.env.APP_FLIDA_REDIRECT_URI || `${process.env.APP_PUBLIC_URL || 'http://localhost:3000'}/auth/callback`,
        postLogoutRedirectUri: process.env.APP_FLIDA_POST_LOGOUT_REDIRECT_URI || '',
        mode: process.env.APP_FLIDA_MODE || 'redirect',
        fedcm: (process.env.APP_FLIDA_FEDCM ?? 'active,passive')
          .split(',')
          .map((mode) => mode.trim())
          .filter(Boolean),
      },
    },
  },

  app: {
    head: {
      title: process.env.APP_PUBLIC_SITE_NAME || 'Template',
      link: [
        { rel: 'icon', type: 'image/png', href: '/images/common/volonterim_logo_transparent.png' },
        { rel: 'stylesheet', href: 'https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap' },
      ],
      meta: [
        { name: 'viewport', content: 'width=device-width, initial-scale=1.0, user-scalable=no' },
        { name: 'description', content: process.env.APP_PUBLIC_SITE_DESCRIPTION || '' },
      ],
    },
  },

  primevue: {
    usePrimeVue: true,
    options: {
      ripple: true,
      theme: ThemeSettings,
      zIndex: {
        modal: 1100,
        overlay: 1100,
        menu: 1100,
        tooltip: 1200,
      },
    },
  },

  imports: {
    dirs: ['composables', 'composables/**', 'core/**/stores', 'utils', 'utils/**', 'config'],
  },

  vite: {
    plugins: [tailwindcss()],
  },
});
