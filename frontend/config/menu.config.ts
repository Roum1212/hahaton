import { LayoutDashboard, Settings } from 'lucide-vue-next';
import { markRaw, type Component } from 'vue';
import { ROUTES } from './routes.config';

export interface IMenuItem {
  label: string;
  icon?: string;
  to?: string;
  badge?: string;
  disabled?: boolean;
  visible?: boolean;
  permissions?: string[];
}

export interface IMenuGroup {
  key: string;
  label: string;
  items: IMenuItem[];
}

export interface IMenuSection {
  key: string;
  label: string;
  lucide: Component;
  groups: IMenuGroup[];
  permissions?: string[];
}

export const MENU: IMenuSection[] = [
  {
    key: 'main',
    label: 'Главная',
    lucide: markRaw(LayoutDashboard),
    groups: [
      {
        key: 'overview',
        label: 'Обзор',
        items: [{ label: 'Дашборд', icon: 'pi pi-chart-bar', to: ROUTES.cabinet }],
      },
    ],
  },
  {
    key: 'settings',
    label: 'Настройки',
    lucide: markRaw(Settings),
    groups: [
      {
        key: 'account',
        label: 'Аккаунт',
        items: [{ label: 'Профиль', icon: 'pi pi-user', to: `${ROUTES.cabinet}/${ROUTES.settings}/${ROUTES.profile}` }],
      },
      {
        key: 'administration',
        label: 'Администрирование',
        items: [
          {
            label: 'Пользователи',
            icon: 'pi pi-users',
            to: `${ROUTES.cabinet}/${ROUTES.settings}/${ROUTES.users}`,
            permissions: ['admin'],
          },
        ],
      },
    ],
  },
];
