import { TreeNode } from 'primeng/api';

export const getFullPath = (node: TreeNode): string => {
  const path = [];
  let current: TreeNode | null = node;
  while (current) {
    path.unshift(current.label);
    current = current.parent as TreeNode;
  }
  return path.join('/');
};
