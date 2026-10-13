export function isValuable<T>(value: T): value is NonNullable<T> {
  return value !== null && value !== undefined;
}

export function isValuableArray<T>(value: T): value is T & unknown[] {
  return Array.isArray(value) && value.length > 0;
}

export function isValuableString(value: unknown): value is string {
  return typeof value === 'string' && value.trim().length > 0;
}

export function isValuableInt(value: unknown): value is number {
  return typeof value === 'number' && Number.isSafeInteger(value);
}
