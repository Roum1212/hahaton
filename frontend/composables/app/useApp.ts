/* eslint-disable @typescript-eslint/no-explicit-any */
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core';

const useApp = () => {
  const toast = useToast();
  const twBreakpoints = useBreakpoints(breakpointsTailwind);

  const useDebounceFn = (fn: (...args: any[]) => void, delay: number) => {
    let timeoutId: ReturnType<typeof setTimeout>;
    return (...args: any[]) => {
      clearTimeout(timeoutId);
      timeoutId = setTimeout(() => fn(...args), delay);
    };
  };

  const useThrottleFn = <Args extends any[]>(fn: (...args: Args) => void, delay = 300) => {
    let lastCall = 0;
    return (...args: Args) => {
      const now = Date.now();
      if (now - lastCall >= delay) {
        lastCall = now;
        fn(...args);
      }
    };
  };

  return {
    toast: getToasts(toast),
    confirm: getConfirm(useConfirm()),
    loader: useLoader(),
    loaderProgress: useLoaderProgress(),
    route: useRoute(),
    breakpoints: twBreakpoints,
    debounce: useDebounceFn,
    throttle: useThrottleFn,
  };
};
export default useApp;
