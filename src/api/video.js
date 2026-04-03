import { http } from '../utils/request';

/**
 * 视频去水印
 * @param {string} url 视频链接
 */
export const removeWatermark = (url) => {
  return http.post('/video/remove-watermark', { url });
};
