import { http } from '../utils/request';

// 获取用户菜单
export const getUserMenus = () => {
  return http.get('/menus/user');
};

// 获取菜单树
export const getMenuTree = () => {
  return http.get('/menus/tree');
};

// 创建菜单
export const createMenu = (data) => {
  return http.post('/menus', data);
};

// 更新菜单
export const updateMenu = (id, data) => {
  return http.put(`/menus/${id}`, data);
};

// 删除菜单
export const deleteMenu = (id) => {
  return http.delete(`/menus/${id}`);
};
