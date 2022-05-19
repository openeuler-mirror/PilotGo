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
 * @Date: 2022-05-12 15:28:36
 * @LastEditTime: 2022-05-16 15:23:59
 * @Description: provide agent log manager of pilotgo
 */
import request from './request'
// 获取所有repo源
export function getRepos(data) {
  return request({
    url: '/config/allRepos',
    method: 'get',
    params: data
  })
}

// 增加repo源
export function createRepo(data) {
  return request({
    url: '/config/createRepo',
    method: 'post',
    data
  })
}

// 删除repo源
export function delRepos(data) {
  return request({
    url: '/config/delRepos',
    method: 'post',
    data
  })
}

// 编辑repo源
export function updateRepo(data) {
  return request({
    url: '/config/updateRepo',
    method: 'post',
    data
  })
}