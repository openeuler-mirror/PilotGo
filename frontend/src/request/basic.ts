/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
 */
import request from './request';

// 在组件或其他地方调用request方法
export const platformVersion = async () => {
  const response = await request({ url: '/version', method: 'GET' });
  return response;
};
