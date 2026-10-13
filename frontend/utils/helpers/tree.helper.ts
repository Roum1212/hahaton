import type { TreeNode } from 'primevue/treenode';
export interface ITreeDto {
  id: number;
  title?: string;
  label?: string;
  name?: string;
  children?: ITreeDto[];
}
export function flattenTree<T extends { id: number; parent?: number | null; children?: T[] }>(tree: T[]): Omit<T, 'children'>[] {
  const result: Omit<T, 'children'>[] = [];

  for (const node of tree) {
    const { children, ...rest } = node;
    result.push(rest);
    if (children && children.length > 0) {
      result.push(...flattenTree(children.map((el) => ({ ...el, parent: node.id }) as T)));
    }
  }

  return result;
}
export function prepareTreeData(data: ITreeDto[], parent: number = 0, draggable: boolean = false): TreeNode[] {
  const treeData: TreeNode[] = [];
  data.forEach((item) => {
    treeData.push({
      key: item.id.toString(),
      label: item.title,
      draggable: !!draggable,
      droppable: !!draggable,
      data: {
        ...filterFieldsNeg(item, 'children'),
        parent,
      },
      children: item.children ? prepareTreeData(item.children, item.id || 0, draggable) : [],
    });
  });
  return treeData;
}
export function filterTreeData(data: TreeNode[], ids: (number | string)[] | undefined): TreeNode[] {
  const convertIds = ids ? ids.map((id) => id.toString()) : undefined;
  if (Array.isArray(convertIds) && convertIds.length === 0) return [];

  return data
    .map((node) => {
      const filteredChildren = filterTreeData(node.children || [], convertIds);

      const isNodeIncluded = convertIds && convertIds.includes(node.data.id.toString());

      if (isNodeIncluded || filteredChildren.length > 0 || isNodeIncluded === undefined) {
        return {
          ...node,
          children: isNodeIncluded ? node.children || [] : filteredChildren,
        };
      }

      return null;
    })
    .filter((node) => node !== null);
}
