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
 * @Date: 2022-03-15 10:12:29
 * @LastEditTime: 2022-05-16 15:24:37
 */
import request from './request'
// 获取角色列表
export function getRoles() {
  return request({
    url: '/user/roles',
    method: 'get',
  })
}

// 获取角色菜单
export function getMenu(data) {
  return request({
    url: '/role/permission',
    method: 'post',
    data
  })
}

// 添加角色
export function addRole(data) {
  return request({
    url: '/user/addRole',
    method: 'post',
    data
  })
}

// 编辑角色
export function updateRole(data) {
  return request({
    url: '/user/updateRole',
    method: 'post',
    data
  })
}

// 给角色赋予权限
export function roleAuth(data) {
  return request({
    url: '/user/roleChange',
    method: 'post',
    data
  })
}

// 删除角色
export function delRole(data) {
  return request({
    url: '/user/delRole',
    method: 'post',
    data
  })
}