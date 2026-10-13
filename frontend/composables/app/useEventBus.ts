/* eslint-disable @typescript-eslint/no-explicit-any */
import { reactive } from 'vue';

type EventMap = Record<string, ((...args: any[]) => void)[]>;

export const useEventBus = () => {
  const listeners: EventMap = reactive({});

  const emit = (event: string, ...args: any[]) => {
    listeners[event]?.forEach((fn) => fn(...args));
  };

  const on = (event: string, fn: (...args: any[]) => void) => {
    listeners[event] = listeners[event] || [];
    listeners[event].push(fn);
  };

  const off = (event: string, fn: (...args: any[]) => void) => {
    listeners[event] = listeners[event]?.filter((l) => l !== fn) || [];
  };

  return { emit, on, off };
};
