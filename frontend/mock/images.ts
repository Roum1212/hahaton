import type http from 'node:http';
import { Code } from '@connectrpc/connect';
import { db, fail, newId, type ImageRec } from './db';

export const STORAGE_PREFIX = '/storage/';
export const MAX_IMAGE_BYTES = 5 * 1024 * 1024;
const MAX_IMAGE_SIDE = 4096;
const PATH_PATTERN = /^[0-9a-f-]{36}\/[0-9a-f-]{36}\.(jpg|png)$/;

interface ImageFormat {
  ext: 'jpg' | 'png';
  contentType: string;
  width: number;
  height: number;
}

const readPng = (data: Buffer): ImageFormat | null => {
  if (data.length < 24 || data.readUInt32BE(0) !== 0x89504e47 || data.toString('ascii', 12, 16) !== 'IHDR') return null;
  return { ext: 'png', contentType: 'image/png', width: data.readUInt32BE(16), height: data.readUInt32BE(20) };
};

const readJpeg = (data: Buffer): ImageFormat | null => {
  if (data.length < 4 || data[0] !== 0xff || data[1] !== 0xd8) return null;
  let offset = 2;
  while (offset + 9 < data.length) {
    if (data[offset] !== 0xff) return null;
    const marker = data[offset + 1]!;
    const length = data.readUInt16BE(offset + 2);
    if (marker >= 0xc0 && marker <= 0xcf && ![0xc4, 0xc8, 0xcc].includes(marker)) {
      return { ext: 'jpg', contentType: 'image/jpeg', width: data.readUInt16BE(offset + 7), height: data.readUInt16BE(offset + 5) };
    }
    offset += 2 + length;
  }
  return null;
};

export const saveImage = (ownerId: string, data: Buffer): ImageRec => {
  if (!data.length || data.length > MAX_IMAGE_BYTES) return fail(Code.InvalidArgument, 'Изображение должно быть не больше 5 МиБ');
  const format = readPng(data) ?? readJpeg(data);
  if (!format) return fail(Code.InvalidArgument, 'Поддерживаются изображения JPEG и PNG');
  if (!format.width || !format.height || format.width > MAX_IMAGE_SIDE || format.height > MAX_IMAGE_SIDE) {
    return fail(Code.InvalidArgument, 'Каждая сторона изображения — не больше 4096 пикселей');
  }
  const image: ImageRec = { path: `${ownerId}/${newId()}.${format.ext}`, ownerId, contentType: format.contentType, data };
  db.images.set(image.path, image);
  return image;
};

export const imageUrl = (path: string) => (path ? `${STORAGE_PREFIX}${path}` : '');

export const resolveImagePath = (path: string, ownerId: string): string => {
  if (!path) return '';
  if (!PATH_PATTERN.test(path)) return fail(Code.InvalidArgument, 'Некорректный путь изображения');
  const image = db.images.get(path);
  if (!image) return fail(Code.NotFound, 'Изображение не найдено');
  if (image.ownerId !== ownerId) return fail(Code.PermissionDenied, 'Изображение принадлежит другому участнику');
  return path;
};

export const serveStorage = (req: http.IncomingMessage, res: http.ServerResponse): boolean => {
  if (req.method !== 'GET' || !req.url?.startsWith(STORAGE_PREFIX)) return false;
  const image = db.images.get(decodeURIComponent(req.url.slice(STORAGE_PREFIX.length).split('?')[0]!));
  if (!image) {
    res.writeHead(404).end();
    return true;
  }
  res.writeHead(200, {
    'Content-Type': image.contentType,
    'Content-Length': image.data.length,
    'Cache-Control': 'public, max-age=31536000',
  });
  res.end(image.data);
  return true;
};
