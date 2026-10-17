import type { ConnectRouter } from '@connectrpc/connect';
import { ImageService } from '../../genproto/image/v1/image_service_pb';
import { requireUser } from '../auth';
import { imageUrl, saveImage } from '../images';

export const registerImageService = (router: ConnectRouter) =>
  router.service(ImageService, {
    uploadImage(req, ctx) {
      const user = requireUser(ctx);
      const image = saveImage(user.id, Buffer.from(req.data));
      return { imagePath: image.path, imageUrl: imageUrl(image.path) };
    },
  });
