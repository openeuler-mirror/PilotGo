/*
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND, 
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * @Author: zhaozhenfang
 * @Date: 2022-01-19 17:30:12
 * @LastEditTime: 2022-04-13 11:04:55
 * @Description: provide agent log manager of pilotgo
 */
import axios from 'axios'
import store from '@/store'
import router from '@/router'

export function request(config) {
  // 1.创建axios实例
  const instance = axios.create({
    baseURL: '',
    timeOut: 5000
  })

  // 2.1添加请求拦截器
  axios.interceptors.request.use(config => {
    // 在发送请求之前做些什么
    let pathname = location.pathname;
    if(localStorage.getItem('token')){
      if(pathname != '/' &&  pathname != '/login'){
        config.headers.common['token'] = localStorage.getItem('token');
      }
    }
    return config;
  }, error => {
    // 对请求错误做些什么
    return Promise.reject(error);
  });

  // 2.2添加响应拦截器
  axios.interceptors.response.use(response => {
    return response;
  },error => {
    if (error.response) {
      switch (error.response.status) {
        // 返回401，清除token信息并跳转到登录页面
        case 401:
          localStorage.removeItem('token');
          router.replace({
            path: '/login'
            //登录成功后跳入浏览的当前页面
            // query: {redirect: router.currentRoute.fullPath}
          });
        case 400: console.log(error.response)
      }
      // 返回接口返回的错误信息
      return Promise.reject(error.response.data);
    }
  });

  // 3 发送真正的网络请求
  // 由于instance实例本身返回的就是promise实例
  // 再调用的时 也就可以使用.then() 拿到成功的 .catch() 返回值失败的返回值
  return instance(config)

}
