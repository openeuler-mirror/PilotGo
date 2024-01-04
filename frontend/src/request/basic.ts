import request from './request';

// 在组件或其他地方调用request方法
export const platformVersion = async () => {
  const response = await request({ url: '/version', method: 'GET' });
  return response;
};
