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
 * @Date: 2022-04-18 15:18:50
 * @LastEditTime: 2022-06-20 11:33:02
 */
import request from './request'
// 重启防火墙
export function reStart(data) {
  return request({
    url: 'agent/firewall_restart',
    method: 'get',
    params: data
  })
}

// 关闭防火墙
export function close(data) {
  return request({
    url: 'agent/firewall_stop',
    method: 'get',
    params: data
  })
}

// 指定区域开放端口
export function openPort(data) {
  return request({
    url: 'agent/firewall_addzp',
    method: 'post',
    data
  })
}


// 删除开放的端口
export function deleteOpenPort(data) {
  return request({
    url: 'agent/firewall_delzp',
    method: 'post',
    data
  })
}

// 获取防火墙配置
export function getConfig(data) {
  return request({
    url: 'api/firewall_config',
    method: 'get',
    params: data
  })
}
