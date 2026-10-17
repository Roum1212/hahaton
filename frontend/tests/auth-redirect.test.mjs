import assert from 'node:assert/strict';
import { test } from 'node:test';
import { loadTypeScript } from './helpers/load-typescript.mjs';

for (const path of ['/auth/login', '/auth/callback']) {
  test(`an existing session redirects from ${path} without component injection`, async () => {
    const redirects = [];
    let loaderVisible = false;
    const globals = {
      AUTH_PAGES: ['/auth/login', '/auth/callback'],
      GUEST_PAGES: ['/auth/login', '/auth/callback'],
      useAuthStore: () => ({ init: async () => {}, isLoggedIn: true, hasCoordinatorProfile: false }),
      useLoader: () => ({ show: () => (loaderVisible = true), hide: () => (loaderVisible = false) }),
      useApp: () => {
        throw new Error('No PrimeVue Toast provided!');
      },
      useCookie: () => ({ value: null }),
      takeLoginRedirect: () => null,
      normalizeLoginRedirect: () => null,
      useCabinetMode: () => ({ homeFor: () => '/cabinet' }),
      navigateTo: (target) => redirects.push(target),
      defineNuxtRouteMiddleware: (handler) => handler,
    };
    const actions = loadTypeScript(new URL('../composables/auth/useAuth.ts', import.meta.url), globals);
    const middleware = loadTypeScript(new URL('../middleware/0.auth.global.ts', import.meta.url), {
      ...globals,
      ...actions,
    }).default;

    await middleware({ path, fullPath: path, query: {} });

    assert.deepEqual(redirects, ['/cabinet']);
    assert.equal(loaderVisible, false);
  });
}
