export const useFormValidation = <Field extends string>(
  root: Ref<HTMLElement | null>,
  getErrors: () => Partial<Record<Field, string>>,
) => {
  const submitted = ref(false);
  const formId = useId();
  const errors = computed<Partial<Record<Field, string>>>(() => (submitted.value ? getErrors() : {}));
  const errorId = (field: Field) => `${formId}-${field}-error`;
  const reset = () => {
    submitted.value = false;
  };

  const validate = async () => {
    submitted.value = true;
    if (!Object.values(errors.value).some(Boolean)) return true;
    await nextTick();
    const field = root.value?.querySelector<HTMLElement>('[data-invalid="true"]');
    field?.scrollIntoView({ behavior: 'smooth', block: 'center' });
    field
      ?.querySelector<HTMLElement>(
        'input:not([type="hidden"]):not(:disabled), textarea:not(:disabled), [tabindex="0"], button:not(:disabled)',
      )
      ?.focus({ preventScroll: true });
    return false;
  };

  return { errors, errorId, reset, validate };
};
