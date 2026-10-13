/* eslint-disable @typescript-eslint/no-explicit-any */
import { ref } from 'vue';
export const useEntityForm = <T extends Record<string, any>>(initDefaultData?: () => T) => {
  const formRef = ref<{ save: () => Promise<boolean> } | null>(null);
  const isSubmitting = ref(false);

  const initialData: T = initDefaultData?.() || ({} as T);

  const submitForm = async (callback?: () => void | Promise<void>) => {
    if (formRef.value && (await formRef.value.save())) {
      if (callback) await callback();
    }
  };

  return {
    formRef,
    submitting: isSubmitting,
    initialData,
    submitForm,
  };
};
