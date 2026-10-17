export interface IImageSize {
  width: number;
  height: number;
}

export const AVATAR_IMAGE_SIZE: IImageSize = { width: 200, height: 200 };

export const COVER_IMAGE_SIZE: IImageSize = { width: 640, height: 400 };

export const MAX_SOURCE_IMAGE_BYTES = 20 * 1024 * 1024;

export const CROPPED_IMAGE_QUALITY = 0.9;

export const MIN_CROP_SIDE = 40;
