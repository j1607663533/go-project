import { http } from '../utils/request';

// 获取所有角色
export const getAllRoles = () => {
  return http.get('/roles');
};

// 获取角色详情
export const getRoleById = (id) => {
  return http.get(`/roles/${id}`);
};

// 创建角色
export const createRole = (data) => {
  return http.post('/roles', data);
};

// 更新角色
export const updateRole = (id, data) => {
  return http.put(`/roles/${id}`, data);
};

// 删除角色
export const deleteRole = (id) => {
  return http.delete(`/roles/${id}`);
};

// 为角色分配菜单
export const assignMenus = (id, menuIds) => {
  return http.post(`/roles/${id}/menus`, { menu_ids: menuIds });
};
