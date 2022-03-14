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
 * @Date: 2022-02-15 17:02:23
 * @LastEditTime: 2022-03-04 16:56:49
 * @Description: provide agent log manager of pilotgo
 */
// api请求接口文件
import { request } from './request'

const API1 = '/plugin/prometheus/api/v1'
export function getChartName(time) {
  return request({
    url: API1 + '/label/__name__/values?_=' + time,
    method: 'get',
  })
}

export function getChart(url) {
  return request({
    url: API1 + url,
    method: 'get',
  })
}

export function FirewallConfig(data) {
  return request({
    url: '/firewall/config',
    method: 'post',
    data
  })
}

export function FirewallStop(data) {
  return request({
    url: '/firewall/stop',
    method: 'post',
    data
  })
}

export function FirewallRestart(data) {
  return request({
    url: '/firewall/restart',
    method: 'post',
    data
  })
}

export function FirewallReload(data) {
  return request({
    url: '/firewall/reload',
    method: 'post',
    data
  })
}

export function FirewallAddZonePort(data) {
  return request({
    url: '/firewall/addzp',
    method: 'post',
    data
  })
}

export function FirewallDelZonePort(data) {
  return request({
    url: '/firewall/delzp',
    method: 'post',
    data
  })
}

export function FirewallAddZonePortPermanent(data) {
  return request({
    url: '/firewall/addzpp',
    method: 'post',
    data
  })
}
