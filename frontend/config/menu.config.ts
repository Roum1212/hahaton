import { ROUTES } from './routes.config';

export type CabinetMode = 'volunteer' | 'coordinator';

export interface INavItem {
  label: string;
  icon: string;
  to: string;
  tab?: string;
  exact?: boolean;
}

export const CABINET_MODE_LABELS: Record<CabinetMode, string> = {
  volunteer: 'Волонтёр',
  coordinator: 'Координатор',
};

export const NAV_BY_MODE: Record<CabinetMode, INavItem[]> = {
  volunteer: [
    { label: 'Профиль', icon: 'pi pi-user', to: ROUTES.volunteer.home, exact: true },
    { label: 'Мероприятия', icon: 'pi pi-calendar', to: ROUTES.volunteer.events },
    { label: 'Поощрения', icon: 'pi pi-gift', to: ROUTES.volunteer.rewards },
    { label: 'Координаторы', icon: 'pi pi-users', to: ROUTES.volunteer.coordinators },
  ],
  coordinator: [
    { label: 'Профиль', icon: 'pi pi-id-card', to: ROUTES.coordinator.home, exact: true },
    { label: 'Мероприятия', icon: 'pi pi-calendar', to: ROUTES.coordinator.events },
    { label: 'Волонтёры', icon: 'pi pi-users', to: ROUTES.coordinator.volunteers },
    { label: 'Заказы', icon: 'pi pi-inbox', to: ROUTES.coordinator.orders },
  ],
};
