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
 * @LastEditTime: 2022-06-22 15:55:16
 * @Description: provide agent log manager of pilotgo
 */
import request from './request'

// 请求组织树接口
export function getDeparts(data) {
  return request({
    url: '/macList/departinfo',
    method: 'get',
    params: data
  })
}
// 点击获取当前子节点接口
export function getChildNode(data) {
  return request({
    url: 'macList/depart',
    method: 'get',
    params: data
  })
}
// 添加节点
export function addDepart(data) {
  return request({
    url: 'macList/adddepart',
    method: 'post',
    data
  })
}
// 编辑节点
export function updateDepart(data) {
  return request({
    url: 'macList/updatedepart',
    method: 'post',
    data
  })
}
// 删除节点
export function deleteDepart(data) {
  return request({
    url: 'macList/deletedepartdata',
    method: 'post',
    data
  })
}
// 拖拽节点

// 点击部门刷新列表接口
export function getClusters(data) {
  return request({
    url: 'macList/machineinfo',
    method: 'get',
    params: data
  })
}

// 获取给定主机列表的tag标签
export function getTags(data) {
  return request({
    url: 'macList/gettags',
    method: 'get',
    params: data
  })
}

// 获取资源池列表接口
export function getSourceMac(data) {
  return request({
    url: 'macList/sourcepool',
    method: 'get',
    params: data
  })
}

// 更换机器所属部门
export function changeMacDept(data) {
  return request({
    url: 'macList/modifydepart',
    method: 'post',
    data
  })
}

// 添加ip接口
export function insertIp(data) {
  return request({
    url: 'macList/addmachine',
    method: 'post',
    data
  })
}
// 编辑ip接口
export function updateIp({ ip, ...data }) {
  return request({
    url: `/hosts/${ip}`,
    method: 'post',
    data
  })
}
// 删除ip接口
export function deleteIp(data) {
  return request({
    url: 'macList/deletemachine',
    method: 'post',
    data
  })
}

// 获取所有的机器列表
export function getallMacIps() {
  return request({
    url: '/macList/machinealldata',
    method: 'get',
  })
}

// 创建批次时通过deptId获取IP列表  macList/machinealldata
export function getMacIps(data) {
  return request({
    url: 'macList/selectmachine',
    method: 'get',
    params: data
  })
}

// repo源
export function repoAll(data) {
  return request({
    url: 'api/repos',
    method: 'get',
    params: data
  })
}

// rpm列表
export function rpmAll(data) {
  return request({
    url: 'api/rpm_all',
    method: 'get',
    params: data
  })
}

// rpm详情
export function getDetail(data) {
  return request({
    url: 'api/rpm_info',
    method: 'get',
    params: data
  })
}

// rpm下发
export function rpmIssue(data) {
  return request({
    url: 'agent/rpm_install',
    method: 'post',
    data
  })
}

// rpm卸载
export function rpmUnInstall(data) {
  return request({
    url: 'agent/rpm_remove',
    method: 'post',
    data
  })
}

// 根据ip获取机器信息
export function getDeviceInfo(data) {
  return request({
    url: 'macList/deviceinfo',
    method: 'get',
    params: data
  })
}

// 获取机器基本信息
export function getBasicInfo(data) {
  return request({
    url: 'api/os_basic',
    method: 'get',
    params: data
  })
}

// 获取机器overview信息
export function getOverview(data) {
  return request({
    url: 'api/agent_overview',
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
    method: 'post',
    data
  })
}

// 关闭一项服务
export function serviceStop(data) {
  return request({
    url: 'agent/service_stop',
    method: 'post',
    data
  })
}

// 重启一项服务
export function serviceRestart(data) {
  return request({
    url: 'agent/service_restart',
    method: 'post',
    data
  })
}

// 获取内核信息
export function getSyskernels(data) {
  return request({
    url: 'api/sysctl_info',
    method: 'get',
    params: data
  })
}

// 获取一个内核信息
export function getOneSyskernel(data) {
  return request({
    url: 'api/sysctl_view',
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

// 获取网络信息io
export function getNetwork(data) {
  return request({
    url: 'api/net',
    method: 'get',
    params: data
  })
}

// 编辑网络配置
export function updateNet(data) {
  return request({
    url: 'agent/network',
    method: 'post',
    data
  })
}


// 获取网络信息nic
export function getNetNic(data) {
  return request({
    url: 'api/net_nic',
    method: 'get',
    params: data
  })
}

// 获取网络信息tcp
export function getNetTcp(data) {
  return request({
    url: 'api/net_tcp',
    method: 'get',
    params: data
  })
}

// 获取网络信息udp
export function getNetUdp(data) {
  return request({
    url: 'api/net_udp',
    method: 'get',
    params: data
  })
}

// 获取任务列表
export function getCronList(data) {
  return request({
    url: 'agent/cron_list',
    method: 'get',
    params: data
  })
}

// 新建任务信息
export function createCron(data) {
  return request({
    url: 'agent/cron_new',
    method: 'post',
    data
  })
}

// 修改任务信息
export function updateCron(data) {
  return request({
    url: 'agent/cron_update',
    method: 'post',
    data
  })
}

// 开启关闭任务状态
export function changeCStatus(data) {
  return request({
    url: 'agent/cron_status',
    method: 'post',
    data
  })
}

// 删除任务信息
export function delCron(data) {
  return request({
    url: 'agent/cron_del',
    method: 'post',
    data
  })
}

// 执行一条任务