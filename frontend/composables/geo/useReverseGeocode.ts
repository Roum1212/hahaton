import type { IGeoPoint, IReverseGeocode } from '~/core/geo/services/geo/geo.types';

export const useReverseGeocode = () => {
  const geoStore = useGeoStore();
  const { toast } = useApp();
  const resolving = ref(false);
  let controller: AbortController | null = null;

  const findCity = async (name: string | null) => {
    if (!name) return null;
    const cities = (await handleAsyncOperation(() => geoStore.searchCities(name, 10))) ?? [];
    const normalized = name.toLocaleLowerCase();
    return cities.find((item) => item.name.toLocaleLowerCase() === normalized) ?? null;
  };

  const resolve = async (point: IGeoPoint): Promise<IReverseGeocode | null | undefined> => {
    controller?.abort();
    const current = new AbortController();
    controller = current;
    resolving.value = true;
    try {
      const result = await geoStore.service.reverse(point, current.signal);
      if (current.signal.aborted) return undefined;
      if (!result?.address) toast.warning('Не удалось определить адрес, укажите его вручную');
      return result;
    } catch {
      if (current.signal.aborted) return undefined;
      toast.warning('Не удалось определить адрес, укажите его вручную');
      return null;
    } finally {
      if (controller === current) resolving.value = false;
    }
  };

  onBeforeUnmount(() => controller?.abort());

  return { resolving, resolve, findCity };
};
