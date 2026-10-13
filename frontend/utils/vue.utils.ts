export async function handleAsyncOperation<T>(
  asyncFn: () => Promise<T>,
  onFinally?: () => void,
  onError?: (error: unknown) => unknown,
  errorPage?: boolean,
): Promise<T | undefined> {
  try {
    return await asyncFn();
  } catch (error) {
    if (onError) {
      return (await onError(error)) as T;
    }

    if (errorPage) {
      throw showError({ statusCode: 404, fatal: true });
    }

    return undefined;
  } finally {
    onFinally?.();
  }
}
