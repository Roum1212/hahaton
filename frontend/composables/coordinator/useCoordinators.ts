export const useCoordinatorNames = () => {
  const store = useCoordinatorsStore();

  const ensure = async (ids: string[]) => {
    const missing = [...new Set(ids)].filter((id) => id && !store.getById(id));
    await Promise.all(missing.map((id) => handleAsyncOperation(() => store.fetchOne(id))));
  };

  const name = (id: string) => store.getById(id)?.displayName ?? 'Координатор';
  const avatar = (id: string) => store.getById(id)?.avatarUrl ?? null;

  return { ensure, name, avatar };
};
