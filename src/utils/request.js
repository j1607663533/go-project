import axios from 'axios';

// 使用 Vite 环境变量，默认为本地开发地址
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080/api/v1";

const service = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
});

/**
 * 请求拦截器
 * 自动注入 Token
 */
service.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("token");
    if (token) {
      config.headers["Authorization"] = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    console.error("请求发送失败:", error);
    return Promise.reject(error);
  }
);

/**
 * 响应拦截器
 * 处理数据脱壳、业务代码错误及 HTTP 状态码错误
 */
service.interceptors.response.use(
  (response) => {
    const res = response.data;

    // 根据后端约定的状态码进行判断 (假设 0 为成功)
    if (res.code !== 0) {
      const errorMsg = res.message || "未知业务错误";
      
      // 特殊处理：未授权/登录过期
      if (res.code === 401) {
        handleAuthError();
      }

      return Promise.reject(new Error(errorMsg));
    }

    // 返回脱壳后的业务数据
    return res.data;
  },
  (error) => {
    let message = "连接服务器失败";
    
    if (error.response) {
      const status = error.response.status;
      switch (status) {
        case 401:
          message = "登录已过期，请重新登录";
          handleAuthError();
          break;
        case 403:
          message = "权限不足，拒绝访问";
          break;
        case 404:
          message = "接口请求地址不存在";
          break;
        case 500:
          message = "服务器内部错误";
          break;
        default:
          message = `网络异常 (${status})`;
      }
    } else if (error.request) {
      message = "网络无响应，请查看网络状态";
    }

    console.error("响应错误:", message);
    return Promise.reject(new Error(message));
  }
);

/**
 * 处理授权相关错误（如 401）
 */
function handleAuthError() {
  localStorage.removeItem("token");
  localStorage.removeItem("user");
}

// 封装常用的请求方法
const http = {
  get(url, params = {}, config = {}) {
    return service.get(url, { params, ...config });
  },
  post(url, data = {}, config = {}) {
    return service.post(url, data, config);
  },
  put(url, data = {}, config = {}) {
    return service.put(url, data, config);
  },
  delete(url, config = {}) {
    return service.delete(url, config);
  }
};

export { http };
export default service;
