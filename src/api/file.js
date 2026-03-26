import { http } from '../utils/request';
// 获取文件列表
export const getFileList = (params) => {
  return http.get('/files', params);
};

// 上传文件
export const uploadFile = (formData) => {
  return http.post('/files/upload',formData);
};

// 删除文件
export const deleteFile = (id) => {
  return http.delete(`/files/${id}`);
};

// 下载文件
export const downloadFile = (id) => {
  return `/files/download/${id}`;
};
