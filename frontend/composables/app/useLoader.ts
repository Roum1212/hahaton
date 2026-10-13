import { ref } from 'vue';

const isShow = ref(false);
const isShowProgress = ref(false);
export const useLoader = () => {
  const show = () => {
    isShow.value = true;
  };

  const hide = () => {
    isShow.value = false;
  };

  return { isShow, show, hide };
};
export const useLoaderProgress = () => {
  const show = () => {
    isShowProgress.value = true;
  };

  const hide = () => {
    isShowProgress.value = false;
  };

  return { isShowProgress, show, hide };
};
