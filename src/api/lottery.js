import request from '../utils/request';

// 获取抽奖活动信息及奖品列表
export const getLotteryInfo = () => {
  return request.get('/lottery/info');
};

// 进行抽奖
export const drawLottery = () => {
  return request.post('/lottery/draw');
};

// 获取我的抽奖记录
export const getMyLotteryRecords = (params) => {
  return request.get('/lottery/records/my', { params });
};

// 获取公开中奖记录跑马灯
export const getLotteryPublicRecords = (params) => {
  return request.get('/lottery/records/public', { params });
};

// =========== 管理员后台 API ============

// 获取抽奖活动列表
export const getLotteryActivities = () => {
  return request.get('/lottery/admin/activities');
};

// 保存抽奖活动与奖品设置
export const saveLotteryActivity = (data) => {
  return request.post('/lottery/admin/activities', data);
};

// 切换状态
export const toggleLotteryActivityStatus = (id, status) => {
  return request.put(`/lottery/admin/activities/${id}/status`, { status });
};

// 删除活动
export const deleteLotteryActivity = (id) => {
  return request.delete(`/lottery/admin/activities/${id}`);
};

// 获取抽奖活动相关所有记录流水（分页、筛选）
export const getLotteryAdminRecords = (params) => {
  return request.get('/lottery/admin/records', { params });
};
