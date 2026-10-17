import type { Transport } from '@connectrpc/connect';
import { BaseRpcApi } from '~/core/app/rpc/base-api';
import { ImageService } from '#genproto/image/v1/image_service_pb';
import type { ApiImage } from './image.api.types';

export class ImageApi extends BaseRpcApi<typeof ImageService> {
  constructor(transport: Transport) {
    super(ImageService, transport);
  }

  async upload(data: Uint8Array, signal?: AbortSignal): Promise<ApiImage.UploadResponse> {
    return this.client.uploadImage({ data }, this.options(signal));
  }
}
