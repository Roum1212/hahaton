export enum ViewModeEnum {
  VIEW = 'view',
  EDIT = 'edit',
}
export enum ViewModeType {
  BRAND = 'brand',
  COLLECTION = 'collection',
  CATEGORY = 'category',
  PRODUCT = 'product',
  ORDER = 'order',
}
export const useViewMode = (catalogId: number | null, type: ViewModeType) => {
  const mode = ref<ViewModeEnum>((localStorage.getItem(`${catalogId}-${type}-editMode`) as ViewModeEnum) || ViewModeEnum.VIEW);
  const updateMode = () => {
    localStorage.setItem(`${catalogId}-${type}-editMode`, mode.value);
  };
  return {
    mode,
    updateMode,
  };
};
