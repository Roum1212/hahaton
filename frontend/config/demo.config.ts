import type { CabinetMode } from './menu.config';

export interface IDemoAccount {
  mode: CabinetMode;
  label: string;
  name: string;
  subject: string;
}

export const DEMO_ISSUER = 'seed';

export const DEMO_ACCOUNTS: IDemoAccount[] = [
  { mode: 'volunteer', label: 'Волонтёр', name: 'Мария Иванова', subject: 'maria.ivanova@example.com' },
  { mode: 'coordinator', label: 'Координатор', name: 'Ирина Ковалёва', subject: 'irina.kovaleva@example.com' },
];
