/* eslint-disable @typescript-eslint/no-explicit-any */
import { Button } from '#components';
import { markRaw } from 'vue';

export interface DialogOptions {
  type?: 'modal' | 'sidebar';
  header?: string;
  class?: string;
  isModal?: boolean;
  position?: 'center' | 'top' | 'bottom' | 'left' | 'right' | 'topleft' | 'topright' | 'bottomleft' | 'bottomright' | undefined;
  size?: 'small' | 'medium' | 'large';
  minimizable?: boolean;
  detailUrl?: string;
}

export interface DialogInstance<PROPS, COMPONENT_PROPS> {
  id: string;
  component: any;
  componentProps?: COMPONENT_PROPS;
  componentRef: any;
  visible: boolean;
  options?: DialogOptions;
  props?: PROPS;
  headerContent?: (modal: DialogInstance<PROPS, COMPONENT_PROPS>) => any;
  footerContent?: (modal: DialogInstance<PROPS, COMPONENT_PROPS>) => any;
  onOpen?: (props: PROPS, options: DialogOptions) => COMPONENT_PROPS | Promise<COMPONENT_PROPS>;
}
export const dialogs = ref<Record<string, DialogInstance<any, any>>>({});
export function useDialogManager() {
  const openDialog = async <P, CP>(component: any, options: Partial<DialogInstance<P, CP>>) => {
    const defaultId = `dlg_${Date.now()}`;
    const id = options.id || defaultId;

    if (options?.onOpen) {
      options.componentProps = await options.onOpen(options.props as P, options.options || {});
    } else {
      options.componentProps = options.props as CP;
    }
    dialogs.value[id] = {
      id,
      visible: true,
      component: markRaw(component),
      componentProps: options.componentProps,
      options: {
        minimizable: true,
        ...options.options,
      },
      footerContent: options.footerContent,
      headerContent: options.headerContent,
      props: options.props || {},
      componentRef: ref(null),
      onOpen: options.onOpen,
    };
    return id;
  };

  const close = async (id: string) => {
    const dlg = dialogs.value[id];
    if (!dlg) return;

    dlg.visible = false;

    await setTimeout(() => {
      dialogs.value[id] = undefined as unknown as DialogInstance<any, any>;
    }, 200);
  };
  const update = (id: string, data: Partial<DialogInstance<any, any>>) => {
    const dlg = dialogs.value[id];
    if (dlg) Object.assign(dlg, data);
  };
  const revertAll = () => {
    Object.values(dialogs.value).forEach((dlg) => dlg && (dlg.visible = true));
  };
  const minimizeAll = () => {
    Object.values(dialogs.value).forEach((dlg) => dlg && (dlg.visible = false));
  };
  const minimize = (id: string) => {
    const dlg = dialogs.value[id];
    if (!dlg) return;

    dlg.visible = false;
  };
  const revert = (id: string) => {
    const dlg = dialogs.value[id];
    if (!dlg) return;

    dlg.visible = true;
  };
  const generateButton = (props: {
    severity: 'secondary' | 'primary' | 'danger' | 'success' | 'warning' | 'info';
    label: string;
    icon?: string;
    loading?: boolean;
    disabled?: boolean;
    text?: boolean;
    classList?: string;
    fn?: () => void;
  }) => {
    const { label, icon, text, severity, classList, fn, disabled, loading } = props;
    return h(Button, {
      label,
      severity,
      icon,
      text,
      loading,
      class: classList,
      disabled,
      onClick: fn,
    });
  };
  const defaultFooter = (
    modal: DialogInstance<any, any>,
    callback?: (res: any) => void,

    buttonParams: {
      closeLabel: string;
      succecssLabel: string;
    } = { closeLabel: 'Закрыть', succecssLabel: 'Сохранить' },
  ) =>
    h('div', { class: 'flex justify-end gap-2' }, [
      generateButton({
        severity: 'secondary',
        label: buttonParams.closeLabel,
        fn: () => {
          close(modal.id);
        },
      }),
      generateButton({
        severity: 'primary',
        label: buttonParams.succecssLabel,
        fn: async () => {
          if (!modal.componentRef) return;
          const res = await modal.componentRef?.save();

          if (res) {
            if (callback) callback(res);
            close(modal.id);
          }
        },
      }),
    ]);

  return { generateButton, defaultFooter, dialogs, openDialog, minimize, revertAll, minimizeAll, close, update, revert };
}
