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
 * @Date: 2022-02-25 16:33:46
 * @LastEditTime: 2022-03-25 10:49:48
 * @Description: provide agent log manager of pilotgo
 */
import { request } from './request'

// 请求组织树接口
export function getDeparts(data) {
  return request({
    url: 'machinemanager/departinfo',
    method: 'get',
    params: data
  })
}
// 点击获取当前子节点接口
export function getChildNode(data) {
  return request({
    url: 'machinemanager/depart',
    method: 'get',
    params: data
  })
}
// 添加节点
export function addDepart(data) {
  return request({
    url: 'machinemanager/adddepart',
    method: 'post',
    params: data
  })
}
// 编辑节点
export function updateDepart(data) {
  return request({
    url: 'machinemanager/updatedepart',
    method: 'get',
    params: data
  })
}
// 删除节点
export function deleteDepart(data) {
  return request({
    url: 'machinemanager/t',
    method: 'get',
    params: data
  })
}
// 拖拽节点

// 点击部门刷新列表接口
export function getClusters(data) {
    return request({
      url: 'machinemanager/machineinfo',
      method: 'get',
      params: data
    })
  }

  // 更换机器所属部门
export function changeMacDept(data) {
  return request({
    url: 'machinemanager/modifydepart',
    method: 'post',
    data
  })
}

// 添加ip接口
export function insertIp(data) {
  return request({
    url: '/machinemanager/addmachine',
    method: 'post',
    data
  })
}
// 编辑ip接口
export function updateIp({ip, ...data}) {
  return request({
    url: `/hosts/${ip}`,
    method: 'post',
    data
  })
}
// 删除ip接口
export function deleteIp(data) {
  return request({
    url: '/machinemanager/deletemachinedata',
    method: 'post',
    params: data
  })
}

// rpm下发
export function rpmIssue(data) {
  return request({
    url: '/agent/rpm_install',
    method: 'post',
    data
  })
}

// rpm卸载
export function rpmUnInstall(data) {
  return request({
    url: '/agent/rpm_remove',
    method: 'post',
    data
  })
}

// 根据ip获取机器信息
export function getDeviceInfo(data) {
  return request({
    url: 'machinemanager/deviceinfo',
    method: 'get',
    params: data
  })
}

// 获取OS
export function getOS(data) {
  return request({
    url: 'api/os_info',
    method: 'get',
    params: data
  })
}

// 获取CPU
export function getCpu(data) {
  return request({
    url: 'api/cpu_info',
    method: 'get',
    params: data
  })
}

// 获取memory
export function getMemory(data) {
  return request({
    url: 'api/memory_info',
    method: 'get',
    params: data
  })
}

// 获取当前user
export function getUser(data) {
  return request({
    url: 'api/user_info',
    method: 'get',
    params: data
  })
}

// 获取所有user
export function getAllUser(data) {
  return request({
    url: 'api/user_all',
    method: 'get',
    params: data
  })
}

// 获取所有服务
export function getserviceList(data) {
  return request({
    url: 'api/service_list',
    method: 'get',
    params: data
  })
}

// 开启一项服务
export function serviceStart(data) {
  return request({
    url: 'agent/service_start',
    method: 'get',
    params: data
  })
}

// 关闭一项服务
export function serviceStop(data) {
  return request({
    url: 'agent/service_stop',
    method: 'get',
    params: data
  })
}

// 获取内核信息
export function getSyskernel(data) {
  return request({
    url: 'api/sysctl_info',
    method: 'get',
    params: data
  })
}

// 修改内核信息
export function changeSyskernel(data) {
  return request({
    url: 'agent/sysctl_change',
    method: 'get',
    params: data
  })
}

// 获取磁盘信息
export function getDisk(data) {
  return request({
    url: 'api/disk_use',
    method: 'get',
    params: data
  })
}