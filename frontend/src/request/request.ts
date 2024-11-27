/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
 */
import axios, { type AxiosRequestConfig } from 'axios';
import { directTo } from "@/router/index"

// 公共定义
export const RespCodeOK = 200
export interface RespInterface {
  code?: number;
  data?: any[];
  msg?: string;
  ok?:boolean;
  page?:number;
  size?: number;
  total?: number;
}
  

// 创建一个axios实例
const instance = axios.create({
  baseURL: '/api/v1', // 设置你的API的基本URL
  timeout: 10000, // 设置请求超时时间
});

// 添加请求拦截器，你可以在这里添加请求头等配置
instance.interceptors.request.use(
  (config) => {
    // 在发送请求之前做些事情
    return config;
  },
  (error) => {
    // 处理请求错误
    return Promise.reject(error);
  }
);

// 添加响应拦截器，你可以在这里处理响应数据
instance.interceptors.response.use(
  (response) => {
    // 在响应之前做些事情
    return response.data; // 只返回响应数据部分
  },
  (error) => {
    // 处理响应错误
    if (error.response.status === 401) {
      // 登录过期
      // TODO: 其他清理工作
      directTo('/login')
      return Promise.reject(error);
    }
    if (error.response.data) {
      // 存在data，后端正常处理
      return error.response.data;
    }
    return {
      code: error.response.status,
      msg: error.response.statusText,
    }
  }
);

// 封装通用的request函数
export default function request(config: AxiosRequestConfig) {
  return instance(config);
}