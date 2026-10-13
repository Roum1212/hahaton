import { onBeforeUnmount, onMounted, ref } from 'vue';
export const useWindowSize = () => {
  const width = ref(window.innerWidth);
  const height = ref(window.innerHeight);

  const onResize = () => {
    width.value = window.innerWidth;
    height.value = window.innerHeight;
  };

  onMounted(() => {
    window.addEventListener('resize', onResize);
  });

  onBeforeUnmount(() => {
    window.removeEventListener('resize', onResize);
  });

  return { width, height };
};
