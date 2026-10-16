import { createClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { ImageService } from '../../genproto/image/v1/image_service_pb';

const maxImageBytes = 5 * 1024 * 1024;

export async function uploadImage(file: File, rpcBaseUrl: string) {
  if (file.size === 0 || file.size > maxImageBytes) {
    throw new Error('Выберите изображение размером до 5 МиБ.');
  }

  if (file.type !== 'image/jpeg' && file.type !== 'image/png') {
    throw new Error('Поддерживаются изображения JPEG и PNG.');
  }

  const client = createClient(
    ImageService,
    createConnectTransport({
      baseUrl: rpcBaseUrl,
      fetch: (input, init) => globalThis.fetch(input, { ...init, credentials: 'include' }),
    }),
  );
  const response = await client.uploadImage({ data: new Uint8Array(await file.arrayBuffer()) });

  return {
    imagePath: response.imagePath,
    imageUrl: new URL(response.imageUrl, rpcBaseUrl).href,
  };
}
