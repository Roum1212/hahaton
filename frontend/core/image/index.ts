import { ImageApi } from './api/image/image.api';
import { ImageService } from './services/image/image.service';

export const IMAGE_API = () => {
  const { $transport } = useNuxtApp();
  return {
    image: new ImageApi($transport),
  };
};

export const IMAGE_SERVICES = () => {
  return {
    ImageService: new ImageService(),
  };
};
