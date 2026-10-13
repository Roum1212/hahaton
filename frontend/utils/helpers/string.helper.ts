export function trimText(text: string | undefined, length: number = 120) {
  if (!text) return '';
  if (text.length > length) {
    return text.substring(0, length) + '...';
  }
  return text;
}
export function truncateText(text: string | undefined, startLength: number, endLength: number) {
  if (!text) {
    return '';
  }
  if (text.length <= startLength + endLength) {
    return text;
  }

  const start = text.slice(0, startLength);
  const end = text.slice(-endLength);
  return `${start}...${end}`;
}

export function stringToSNumber(string: string) {
  return string.replace(',', '.');
}
export function normalizeHtml(value: string | number | boolean | null) {
  if (!value) return value;
  if (typeof value === 'string') {
    return value
      .replace(/&nbsp;/g, ' ')
      .replace(/\s+/g, ' ')
      .trim();
  }
  return value;
}
