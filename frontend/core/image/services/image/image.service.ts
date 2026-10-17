import { ConnectError, Code } from '@connectrpc/connect';
import { IMAGE_API } from '../..';
import type { IUploadedImage } from './image.types';

export class ImageService {
  private service = IMAGE_API().image;
  private baseUrl = String(useRuntimeConfig().public.rpcUrl);

  async upload(file: File, signal?: AbortSignal): Promise<IUploadedImage> {
    const error = validateImageFile(file);
    if (error) throw new ConnectError(error, Code.InvalidArgument);
    const { imagePath, imageUrl } = await this.service.upload(new Uint8Array(await file.arrayBuffer()), signal);
    return { path: imagePath, url: this.resolveUrl(imageUrl) ?? '' };
  }

  resolveUrl(url: string | undefined): string | null {
    return resolveImageUrl(url, this.baseUrl);
  }
}
