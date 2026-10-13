/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable @typescript-eslint/no-empty-object-type */
import { onBeforeUnmount, ref, watchEffect } from 'vue';

interface Options extends IntersectionObserverInit {}

export const useIntersectionObserver = (target: HTMLElement | null | Ref<HTMLElement | null>, options: Options = {}) => {
  const isIntersecting = ref(false);
  let observer: IntersectionObserver | null = null;

  const stopWatch = watchEffect(() => {
    const el = (target as any)?.value ?? target;
    if (!el) return;

    observer?.disconnect();

    observer = new IntersectionObserver((entries) => {
      const entry = entries[0];
      if (!entry) return;
      isIntersecting.value = entry.isIntersecting;
    }, options);

    observer.observe(el);
  });

  onBeforeUnmount(() => {
    observer?.disconnect();
    stopWatch();
  });

  return isIntersecting;
};
