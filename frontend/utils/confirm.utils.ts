import type { ConfirmationOptions } from 'primevue/confirmationoptions';
type ConfirmType = { require: (option: ConfirmationOptions) => void; close: () => void };
const confirmDelete = (confirm: ConfirmType) => (message: string, options?: ConfirmationOptions & { loader?: 'base' | 'progress' }) => {
  const confirmLoader = options?.loader === 'progress' ? useLoaderProgress() : useLoader();
  return confirm.require({
    group: 'confirmDelete',
    message,
    header: options?.header || 'Подтверждение',
    icon: 'pi pi-exclamation-triangle',
    rejectClass: 'p-button-secondary p-button-outlined',
    rejectLabel: options?.rejectLabel || 'Закрыть',
    acceptLabel: options?.acceptLabel || 'Удалить',
    accept: async () =>
      await handleAsyncOperation(
        async () => {
          confirmLoader.show();
          if (options?.accept) {
            await options.accept();
          }
        },
        () => {
          confirmLoader.hide();
        },
        (e) => {
          console.error(e);
        },
      ),
    reject: options?.reject,
  });
};

const confirmAccept = (confirm: ConfirmType) => (message: string, options?: ConfirmationOptions & { loader?: 'base' | 'progress' }) => {
  const confirmLoader = options?.loader === 'progress' ? useLoaderProgress() : useLoader();
  return confirm.require({
    group: 'confirmSave',
    message,
    header: options?.header || 'Подтверждение',
    icon: 'pi pi-exclamation-triangle',
    rejectClass: 'p-button-secondary p-button-outlined',
    rejectLabel: 'Отменить',
    acceptLabel: 'Подтвердить',
    accept: async () =>
      await handleAsyncOperation(
        async () => {
          confirmLoader.show();
          if (options?.accept) await options.accept();
        },
        () => {
          confirmLoader.hide();
        },
      ),
  });
};

export function getConfirm(confirm: ConfirmType) {
  return {
    delete: confirmDelete(confirm),
    accept: confirmAccept(confirm),
  };
}
