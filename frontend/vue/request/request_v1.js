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
 * @LastEditTime: 2023-09-02 10:12:26
 * @Description:  v1 api client encapsulation
 */
import axios from 'axios'
import store from '@/store'
import router from '@/router'
import { getToken } from '@/utils/auth'

// 1.创建axios实例
export const request_v1 = axios.create({
  baseURL: '/api/v1',
  timeOut: 5000
})

// 2.1添加请求拦截器
request_v1.interceptors.request.use(config => {
  if (store.getters.token) {
    config.headers['authToken'] = getToken()
  }
  config.baseURL = "/api/v1"
  return config
}, error => {
  return Promise.reject(error);
});

// 2.2添加响应拦截器
request_v1.interceptors.response.use(response => {
  if (response.data && response.data.code == '401') {
    store.dispatch('logOutFont').then(function(){
        router.push("/login")
    })
  } else {
    return response;
  }
  return response;
},error => {
  if (error.response) {
    switch (error.response.status) {
      case 401:
        store.dispatch('logOut')
        router.push("/login")
      }
    return Promise.reject(error.response.data);
  }
});
