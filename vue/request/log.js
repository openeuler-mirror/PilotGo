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
 * @Date: 2022-02-28 14:26:23
 * @LastEditTime: 2022-03-03 13:48:35
 * @Description: provide agent log manager of pilotgo
 */
import { request } from './request'
// 日志列表
export function getLogs(data) {
  return request({
    url: 'agent/log_all',
    method: 'get',
    params: data
  })
}

// 日志详情
export function getLogDetail(data) {
  return request({
    url: 'agent/logs',
    method: 'get',
    params: data
  })
}

// 日志删除
export function deleteLog(data) {
  return request({
    url: 'agent/delete',
    method: 'post',
    data
  })
}
