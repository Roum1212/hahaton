import { ref } from 'vue';

export const useClipboard = (message?: string) => {
  const copied = ref(false);
  const { toast } = useApp();
  const successMessage = message || 'Информация скопирована в буфер обмена';
  const copy = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text);
      copied.value = true;
      toast.info(successMessage);
      setTimeout(() => (copied.value = false), 1000);
    } catch (e) {
      copied.value = false;
      console.error('Clipboard error:', e);
    }
  };

  return { copy, copied };
};
