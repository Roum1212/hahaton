import type { ToastMessageOptions } from 'primevue/toast';
import type { ToastServiceMethods } from 'primevue/toastservice';
type ToastButton = {
  label: string;
  action: () => void;
  severity?: string;
};

type ToastOptions = {
  summary?: string;
  buttons?: ToastButton[];
};

const toastError = (toast: ToastServiceMethods) => (message: string, options?: ToastOptions) =>
  toast.add({
    severity: 'error',
    summary: options?.summary ?? 'Ошибка',
    detail: message,
    life: 3000,
    data: {
      buttons: options?.buttons ?? [],
    },
  } as ToastMessageOptions);
const toastWarning = (toast: ToastServiceMethods) => (message: string, options?: ToastOptions) =>
  toast.add({
    severity: 'warn',
    summary: options?.summary ?? 'Внимание',
    detail: message,
    life: 3000,
    data: {
      buttons: options?.buttons ?? [],
    },
  } as ToastMessageOptions);
const toastSuccess = (toast: ToastServiceMethods) => (message: string, options?: ToastOptions) =>
  toast.add({
    severity: 'success',
    summary: options?.summary ?? 'Успешно',
    detail: message,
    life: 3000,
    data: {
      buttons: options?.buttons ?? [],
    },
  } as ToastMessageOptions);
const toastInfo = (toast: ToastServiceMethods) => (message: string, options?: ToastOptions) =>
  toast.add({
    severity: 'info',
    summary: options?.summary ?? 'Успешно',
    detail: message,
    life: 3000,
    data: {
      buttons: options?.buttons ?? [],
    },
  } as ToastMessageOptions);

export function getToasts(toast: ToastServiceMethods) {
  return {
    error: toastError(toast),
    warning: toastWarning(toast),
    success: toastSuccess(toast),
    info: toastInfo(toast),
  };
}
