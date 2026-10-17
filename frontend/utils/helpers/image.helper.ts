export const IMAGE_ACCEPT = 'image/jpeg,image/png';
export const MAX_IMAGE_BYTES = 5 * 1024 * 1024;

export function validateImageFile(file: File): string | null {
  if (file.type !== 'image/jpeg' && file.type !== 'image/png') return 'Поддерживаются изображения JPEG и PNG';
  if (file.size === 0 || file.size > MAX_IMAGE_BYTES) return 'Выберите изображение размером до 5 МиБ';
  return null;
}

export function validateSourceImage(file: File): string | null {
  if (file.type !== 'image/jpeg' && file.type !== 'image/png') return 'Поддерживаются изображения JPEG и PNG';
  if (file.size === 0 || file.size > MAX_SOURCE_IMAGE_BYTES) return 'Выберите изображение размером до 20 МиБ';
  return null;
}

export function resolveImageUrl(url: string | undefined, baseUrl: string): string | null {
  return url ? new URL(url, baseUrl).href : null;
}
