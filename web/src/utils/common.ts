/**
 * 通用的方法
 *
 * @author xiaoRui
 */

interface TreeNode {
  [key: string]: any;
  children?: TreeNode[];
}

interface TreeConfig {
  id: string;
  parentId: string;
  childrenList: string;
}

export default {
  // 展开树形数据方法
  handleTree(data: TreeNode[], id?: string, parentId?: string, children?: string): TreeNode[] {
    const config: TreeConfig = {
      id: id || 'id',
      parentId: parentId || 'parentId',
      childrenList: children || 'children'
    };
    
    const childrenListMap: { [key: string]: TreeNode[] } = {};
    const nodeIds: { [key: string]: TreeNode } = {};
    const tree: TreeNode[] = [];

    for (const d of data) {
      const parentId = d[config.parentId];
      if (childrenListMap[parentId] == null) {
        childrenListMap[parentId] = [];
      }
      nodeIds[d[config.id]] = d;
      childrenListMap[parentId].push(d);
    }

    for (const d of data) {
      const parentId = d[config.parentId];
      if (nodeIds[parentId] == null) {
        tree.push(d);
      }
    }

    for (const t of tree) {
      adaptToChildrenList(t);
    }

    function adaptToChildrenList(o: TreeNode) {
      if (childrenListMap[o[config.id]] !== null) {
        o[config.childrenList] = childrenListMap[o[config.id]];
      }
      if (o[config.childrenList]) {
        for (const c of o[config.childrenList]) {
          adaptToChildrenList(c);
        }
      }
    }

    return tree;
  }
}