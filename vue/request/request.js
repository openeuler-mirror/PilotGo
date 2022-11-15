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
 * @LastEditTime: 2022-05-16 16:49:18
 * @Description: provide agent log manager of pilotgo
 */
import axios from 'axios'
import store from '@/store'
import router from '@/router'
import { getToken } from '@/utils/auth'


// 1.创建axios实例
const request = axios.create({
  baseURL: '/api/v1',
  timeOut: 5000
})

// 2.1添加请求拦截器
request.interceptors.request.use(config => {
  if (store.getters.token) {
    config.headers['authToken'] = getToken()
  }
  return config
}, error => {
  return Promise.reject(error);
});

// 2.2添加响应拦截器
request.interceptors.response.use(response => {
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
        store.dispatch('logOutFont').then(function() {
          router.push("/login")
      })
    }
    /* Message({
      message: error.message,
      type: 'error',
      duration: 5*1000
    }) */
    return Promise.reject(error.response.data);
  }
});

export default request
