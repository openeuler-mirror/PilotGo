/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
 */
import request from "./request";

// 分页获取部门机器信息
export function getPagedDepartMachines(data: any) {
  return request({
    url: "macList/machineinfo",
    method: "get",
    params: data,
  });
}

// 获取部门所有机器信息
export function getDepartMachines(data: any) {
  return request({
    url: "macList/selectmachine",
    method: "get",
    params: data,
  });
}

// 获取平台所有机器
export function getAllMachines() {
  return request({
    url: "macList/machineinfo_nopage",
    method: "get",
  });
}

// 删除机器接口
export function deleteMachine(data: any) {
  return request({
    url: "macList/deletemachine",
    method: "post",
    data,
  });
}

// 获取指定部门下的子部门
export function getSubDepartment(data: any) {
  return request({
    url: "macList/depart",
    method: "get",
    params: data,
  });
}

// 编辑部门节点
export function updateDepartment(data: any) {
  return request({
    url: "macList/updatedepart",
    method: "post",
    data,
  });
}

// 删除部门节点
export function deleteDepartment(data: any) {
  return request({
    url: "macList/deletedepartdata",
    method: "post",
    data,
  });
}

// 添加部门节点
export function addDepartment(data: any) {
  return request({
    url: "macList/adddepart",
    method: "post",
    data,
  });
}

// 更换机器所属部门
export function changeDepartment(data: any) {
  return request({
    url: "macList/modifydepart",
    method: "post",
    data,
  });
}

// 获取机器overview信息
export function getMachineOverview(data: any) {
  return request({
    url: "api/agent_overview",
    method: "get",
    params: data,
  });
}

// 获取所有服务
export function getserviceList(data: any) {
  return request({
    url: "api/service_list",
    method: "get",
    params: data,
  });
}

// 获取当前机器登录user
export function getCurrentUser(data: any) {
  return request({
    url: "api/user_info",
    method: "get",
    params: data,
  });
}

// 获取机器上所有user
export function getMachineAllUser(data: any) {
  return request({
    url: "api/user_all",
    method: "get",
    params: data,
  });
}

// 获取所有服务
export function getServiceList(data: any) {
  return request({
    url: "api/service_list",
    method: "get",
    params: data,
  });
}

// 关闭一项服务
export function stopService(data: any) {
  return request({
    url: "agent/service_stop",
    method: "post",
    data,
  });
}

// 开启一项服务
export function startService(data: any) {
  return request({
    url: "agent/service_start",
    method: "post",
    data,
  });
}

// 重启一项服务
export function restartService(data: any) {
  return request({
    url: "agent/service_restart",
    method: "post",
    data,
  });
}

// 获取网络信息
export function getNetworkInfo(data: any) {
  return request({
    url: "api/net",
    method: "get",
    params: data,
  });
}

// 获取内核信息
export function getSysctlInfo(data: any) {
  return request({
    url: "api/sysctl_info",
    method: "get",
    params: data,
  });
}

// 获取所有repo源
export function getRepos(data: any) {
  return request({
    url: "api/repos",
    method: "get",
    params: data,
  });
}

// 获取所有已安装的package
export function getInstalledPackages(data: any) {
  return request({
    url: "api/rpm_all",
    method: "get",
    params: data,
  });
}

// 获取单个packge的详情
export function getPackageDetail(data: any) {
  return request({
    url: "api/rpm_info",
    method: "get",
    params: data,
  });
}

// 安装软件包
export function installPackage(data: any) {
  return request({
    url: "agent/rpm_install",
    method: "post",
    data,
  });
}

// 卸载软件包
export function removePackage(data: any) {
  return request({
    url: "agent/rpm_remove",
    method: "post",
    data,
  });
}

// 获取给定主机列表的tag标签
export function getMachineTags(data: any) {
  return request({
    url: "macList/gettags",
    method: "post",
    data,
  });
}
