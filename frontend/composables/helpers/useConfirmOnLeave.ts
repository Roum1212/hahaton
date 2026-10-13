/* eslint-disable @typescript-eslint/no-explicit-any */
import { onBeforeRouteLeave } from 'vue-router';

export function onConfirmOnLeave(isShow: boolean, save: () => Promise<any>) {
  const confirm = useConfirm();
  onBeforeRouteLeave((_, __, next) => {
    if (isShow) {
      confirm.require({
        group: 'confirmSave',
        message: 'У вас есть несохраненные изменения. Вы уверены, что хотите покинуть страницу?',
        header: 'Подтверждение',
        icon: 'pi pi-exclamation-triangle',
        rejectClass: 'p-button-secondary p-button-outlined',
        rejectLabel: 'Уйти',
        acceptLabel: 'Сохранить',
        accept: async () => {
          await save();
          next();
        },
        reject: () => {
          next();
        },
      });
    } else {
      next();
    }
  });
}
