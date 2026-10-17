export interface IProfileTab {
  key: string;
  label: string;
}

export const VOLUNTEER_TABS: IProfileTab[] = [
  { key: 'applications', label: 'Мои заявки' },
  { key: 'purchases', label: 'Покупки' },
  { key: 'wallet', label: 'История койнов' },
  { key: 'about', label: 'Анкета' },
];

export const COORDINATOR_TABS: IProfileTab[] = [
  { key: 'overview', label: 'Нужно внимание' },
  { key: 'rewards', label: 'Поощрения' },
  { key: 'about', label: 'О вас' },
];

export const COORDINATOR_PUBLIC_TABS: IProfileTab[] = [
  { key: 'events', label: 'Мероприятия' },
  { key: 'rewards', label: 'Поощрения' },
];

export const VOLUNTEER_PUBLIC_TABS: IProfileTab[] = [
  { key: 'participations', label: 'Участия' },
  { key: 'rejections', label: 'Отказы' },
  { key: 'orders', label: 'Заказы' },
];

export const EVENT_TABS: IProfileTab[] = [
  { key: 'about', label: 'О мероприятии' },
  { key: 'chat', label: 'Чат' },
];

export const OWNED_EVENT_TABS: IProfileTab[] = [
  { key: 'applications', label: 'Заявки волонтёров' },
  { key: 'chat', label: 'Чат' },
];

export const EVENT_ILLUSTRATIONS = [
  '/images/landing/step-catalog.jpg',
  '/images/landing/step-request.jpg',
  '/images/landing/step-handover.jpg',
  '/images/landing/step-results.jpg',
  '/images/landing/step-accrual.jpg',
];

export const REWARD_ILLUSTRATIONS = [
  '/images/landing/step-reward.jpg',
  '/images/landing/feature-rewards.webp',
  '/images/landing/feature-coins.webp',
];

export const pickIllustration = (list: string[], id: string): string => {
  const hash = [...id].reduce((sum, char) => (sum * 31 + char.charCodeAt(0)) >>> 0, 7);
  return list[hash % list.length]!;
};
