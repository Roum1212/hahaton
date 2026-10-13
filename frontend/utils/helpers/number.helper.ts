export function formatDecimal(value: number | string | null | undefined, precision = 2): string {
  if (value === null || value === undefined || value === '') return '';
  const number = typeof value === 'number' ? value : parseFloat(value);
  if (Number.isNaN(number)) return '';
  return number.toLocaleString('ru-RU', { minimumFractionDigits: precision, maximumFractionDigits: precision });
}

export function convertToNumber(value: string | number): number {
  if (!value) return 0;
  if (typeof value === 'number') return value;
  const cleanedStr = value.replace(/[^\d.-]/g, '');
  const number = parseFloat(cleanedStr);
  return number;
}
