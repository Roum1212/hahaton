/* eslint-disable @typescript-eslint/no-explicit-any */
export function listToMap<T, K extends string | number | symbol>(list: T[], keySelector: (item: T) => K): Record<K, T>;
export function listToMap<T, K extends string | number | symbol, V>(
  list: T[],
  keySelector: (item: T) => K,
  valueSelector?: (item: T) => V,
): Record<K, V>;
export function listToMap<T, K extends string | number | symbol, V = T>(
  list: T[],
  keySelector: (item: T) => K,
  valueSelector?: (item: T) => V,
): Record<K, V> {
  const res = {} as Record<K, V>;
  if (!list) return res;
  for (const item of list) {
    const key = keySelector(item);
    res[key] = valueSelector ? valueSelector(item) : (item as unknown as V);
  }
  return res;
}
export function toArray(value: any): (typeof value)[] {
  return Array.isArray(value) ? value : [value];
}
export function chunkArray<T>(array: T[], chunkSize: number): T[][] {
  if (chunkSize <= 0) throw new Error('Chunk size must be greater than 0');
  const chunks: T[][] = [];
  for (let i = 0; i < array.length; i += chunkSize) {
    chunks.push(array.slice(i, i + chunkSize));
  }

  return chunks;
}
