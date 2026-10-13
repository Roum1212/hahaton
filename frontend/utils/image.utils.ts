export const PLACEHOLDER_IMAGE = '/images/common/image-placeholder.svg';
export const PLACEHOLDER_AVATAR = '/images/common/avatar_placeholder.webp';
export const PLACEHOLDER_FILE = '/images/files/file-default.svg';

export function getImageUrl(): string {
  return useRuntimeConfig().public.imagesPath;
}

export async function fileToBase64(file: File): Promise<string | ArrayBuffer | null> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();

    reader.onload = () => {
      resolve(reader.result);
    };

    reader.onerror = (error) => {
      reject(error);
    };

    reader.readAsDataURL(file);
  });
}
export function base64ToFile(base64: string, filename = 'file.png'): File {
  const arr = base64.split(',');
  const mime = arr[0]!.match(/:(.*?);/)?.[1] || 'image/png';
  const bstr = atob(arr[arr.length - 1] ?? '');
  let n = bstr.length;
  const u8arr = new Uint8Array(n);

  while (n--) {
    u8arr[n] = bstr.charCodeAt(n);
  }

  return new File([u8arr], filename, { type: mime });
}
export function convertToServerBase64(base64: string): string {
  const extensionData = base64.split(';')[0]!.split('/');
  const extensionString = extensionData[extensionData.length - 1];
  const ivb = base64.split('base64,')[1];
  return `b64:${extensionString}:${ivb}`;
}

export function convertFromServerBase64(imageB64: string | undefined): string {
  if (imageB64) {
    const imageExtensions = ['jpg', 'png', 'webp', 'jpeg'];
    const _temp_image_data = imageB64.split(':');
    const extension = _temp_image_data[1] as string;
    if (imageExtensions.includes(extension)) {
      return `data:image/${extension};base64,${_temp_image_data[2]}`;
    }
  }
  return PLACEHOLDER_IMAGE;
}

export function isBase64(str: string): boolean {
  if (str === '' || str.trim() === '') return false;
  try {
    return btoa(atob(str)) === str;
  } catch {
    return false;
  }
}
export function isServerBase64(str: string): boolean {
  return str.startsWith('b64:');
}
export function getExtension(file: string): string {
  let extension: string = '';
  if (isServerBase64(file)) {
    extension = file.split(':')![1] as string;
  } else {
    if (file.includes('.')) {
      const parts = file.split('.');
      extension = parts[parts.length - 1] as string;
    }
  }
  const availableExtensions = ['doc', 'pdf', 'csv', 'pla', 'txt', 'excel-filled'];
  return availableExtensions.includes(extension) ? extension : 'default';
}
export function validateImage(file: string): boolean {
  const allowedExtensions = ['png', 'jpg', 'jpeg', 'webp'];
  const extension = file.split('.').pop() || '';
  return allowedExtensions.includes(extension.toLowerCase());
}
