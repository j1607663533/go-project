import { http } from '../utils/request';

// 获取全部订单列表
export const getAllOrders = () => {
  return http.get('/orders/all');
};

/**
 * 分页获取订单列表（支持筛选）
 * @param {number} page 页码
 * @param {number} pageSize 每页数量
 * @param {string|null} productId 商品ID (支持模糊搜索)
 * @param {string|null} status 订单状态
 */
export const getOrdersByPage = (page = 1, pageSize = 10, productId = null, status = null) => {
  const params = {
    page,
    page_size: pageSize
  };

  if (productId !== null && productId !== '') params.product_id = productId;
  if (status) params.status = status;

  return http.get('/orders', params);
};

// 获取订单详情
export const getOrderDetail = (id) => {
  return http.get(`/orders/${id}`);
};

// 创建订单
export const createOrder = (orderData) => {
  return http.post('/orders', orderData);
};

// 更新订单
export const updateOrder = (id, orderData) => {
  return http.put(`/orders/${id}`, orderData);
};

// 删除订单
export const deleteOrder = (id) => {
  return http.delete(`/orders/${id}`);
};
