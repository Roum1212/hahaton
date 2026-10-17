import { defineStore } from 'pinia';

export const useImagesStore = defineStore('images', {
  getters: {
    service: () => useServices().IMAGE.ImageService,
  },
  actions: {
    async upload(file: File) {
      return this.service.upload(file);
    },
  },
});
