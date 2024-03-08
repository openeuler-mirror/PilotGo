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
 * @Date: 2022-03-04 16:56:07
 * @LastEditTime: 2022-05-25 15:58:46
 */
import request from './request'
import request_prome from './request_prome'
export function getData(data) {
  return request({
    url: '/prometheus/queryrange',
    method: 'get',
    params: data
  })
}
export function getCurrData(data) {
  return request({
    url: '/prometheus/query',
    method: 'post',
    data
  })
}
// 获取所有监控机器的ip列表
export function getPromeIp(data) {
  return request({
    url: '/machinemanager/machinealldata',
    method: 'get',
    params: data
  })
}
// 告警信息列表
// export function getAlerts() {
//   return request_prome({
//     url: '/plugin/Prometheus/api/v1/alerts',
//     method: 'get',
//   })
// }

// 发送告警信息 
export function sendMessage(data) {
  return request_prome({
    url: '/prometheus/alertmanager',
    method: 'post',
    data
  })
}
// 获取首页看板机器数据总量
export function getPanelDatas() {
  return request({
    url: '/overview/info',
    method: 'get',
  })
}
// 获取首页看板给部门机器数据
export function getDeptDatas() {
  return request({
    url: '/overview/depart_info',
    method: 'get',
  })
}