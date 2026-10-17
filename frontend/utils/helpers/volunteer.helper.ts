export const COINS_PER_HALF_HOUR = 5;
export const MINUTES_STEP = 30;

export function pluralize(count: number, forms: [string, string, string]): string {
  const abs = Math.abs(count) % 100;
  const last = abs % 10;
  if (abs > 10 && abs < 20) return forms[2];
  if (last > 1 && last < 5) return forms[1];
  if (last === 1) return forms[0];
  return forms[2];
}

export function formatCoins(count: number): string {
  return `${count} ${pluralize(count, ['койн', 'койна', 'койнов'])}`;
}

export function formatSignedCoins(count: number): string {
  return `${count > 0 ? '+' : count < 0 ? '−' : ''}${formatCoins(Math.abs(count))}`;
}

export function coinsForMinutes(minutes: number): number {
  return Math.floor(minutes / MINUTES_STEP) * COINS_PER_HALF_HOUR;
}

export function formatHours(minutes: number): string {
  return (minutes / 60).toLocaleString('ru-RU', { maximumFractionDigits: 1 });
}

export function formatDuration(minutes: number): string {
  const hours = Math.floor(minutes / 60);
  const rest = minutes % 60;
  if (!hours) return `${rest} мин`;
  return rest ? `${hours} ч ${rest} мин` : `${hours} ч`;
}

export function formatDateTime(value: Date | null | undefined): string {
  return value ? formatDate(value, 'd MMMM yyyy, HH:mm') : '—';
}

export function formatEventPeriod(startsAt: Date | null, endsAt: Date | null): string {
  if (!startsAt) return '—';
  const start = formatDate(startsAt, 'd MMMM yyyy, HH:mm');
  if (!endsAt) return start;
  const sameDay = formatDate(startsAt, 'yyyy-MM-dd') === formatDate(endsAt, 'yyyy-MM-dd');
  return sameDay ? `${start}–${formatDate(endsAt, 'HH:mm')}` : `${start} — ${formatDate(endsAt, 'd MMMM yyyy, HH:mm')}`;
}

export function isPast(value: Date | null | undefined): boolean {
  return !!value && value.getTime() <= Date.now();
}
