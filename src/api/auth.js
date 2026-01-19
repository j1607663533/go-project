import { http } from '../utils/request';

// 获取验证码
export const getCaptcha = (refreshId = null) => {
  return http.get('/captcha', refreshId ? { refresh: refreshId } : {});
};

// 用户登录
export const login = async (username, password, captchaId, captcha) => {
  const data = await http.post('/login', {
    username,
    password,
    captcha_id: captchaId,
    captcha,
  });

  localStorage.setItem("token", data.token);
  localStorage.setItem("user", JSON.stringify(data.user));
  localStorage.setItem("menus", JSON.stringify(data.menus || []));
  return data;
};

// 用户注册
export const register = (username, email, password, nickname, captchaId, captcha) => {
  return http.post('/register', {
    username,
    email,
    password,
    nickname,
    captcha_id: captchaId,
    captcha,
  });
};

// 获取用户信息
export const getProfile = () => {
  return http.get('/auth/profile');
};

// 获取所有用户列表
export const getAllUsers = () => {
  return http.get('/users');
};

// 更新用户
export const updateUser = (id, data) => {
  return http.put(`/users/${id}`, data);
};

// 退出登录
export const logout = async () => {
  try {
    const token = localStorage.getItem("token");
    if (token) {
      await http.post('/logout');
    }
  } catch (error) {
    console.error("退出登录请求失败:", error);
  } finally {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    localStorage.removeItem("menus");
  }
};

// 检查是否已登录
export const isAuthenticated = () => {
  return !!localStorage.getItem("token");
};

// 获取当前用户
export const getCurrentUser = () => {
  const userStr = localStorage.getItem("user");
  return userStr ? JSON.parse(userStr) : null;
};

// 获取用户菜单
export const getUserMenus = () => {
  const menusStr = localStorage.getItem("menus");
  return menusStr ? JSON.parse(menusStr) : [];
};
