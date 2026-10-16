export const ROUTES = {
  home: '/',
  cabinet: '/cabinet',
  login: '/auth/login',
  authCallback: '/auth/callback',
  access: '/access',
  invite: '/invite',
  volunteer: {
    home: '/cabinet',
    events: '/cabinet/events',
    coordinators: '/cabinet/coordinators',
    rewards: '/cabinet/rewards',
  },
  coordinator: {
    home: '/cabinet/coordinator',
    events: '/cabinet/coordinator/events',
    volunteers: '/cabinet/coordinator/volunteers',
    orders: '/cabinet/coordinator/orders',
  },
} as const;

export const ROLES = {
  volunteer: 'volunteer',
  coordinator: 'coordinator',
} as const;

export const AUTH_PAGES: string[] = [ROUTES.login, ROUTES.authCallback];
export const GUEST_PAGES: string[] = [ROUTES.home, ...AUTH_PAGES, ROUTES.access];
