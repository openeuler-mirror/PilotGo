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
 * @Date: 2022-03-03 17:00:36
 * @LastEditTime: 2022-04-15 15:37:25
 * @Description: provide agent log manager of pilotgo
 */
import { request } from './request'
// 用户登录
export function loginByEmail(data) {
  return request({
    url: '/user/login',
    method: 'post',
    data
  })
}
// 用户退出
export function logout(data) {
  return request({
    url: '/user/logout',
    method: 'get',
    params: data
  })
}
// 获取当前登录用户信息
export function getCurUser() {
  return request({
    url: '/user/info',
    method: 'get'
  })
}
// 获取用户权限
export function getPermission(data) {
  return request({
    url: '/user/permission',
    method: 'post',
    data
  })
}
// 获取全部用户信息
export function getUsers(data) {
  return request({
    url: '/user/searchAll',
    method: 'get',
    params: data
  })
}
// 添加用户
export function addUser(data) {
  return request({
    url: '/user/register',
    method: 'post',
    data
  })
}
// 编辑用户
export function updateUser(data) {
  return request({
    url: '/user/update',
    method: 'post',
    data
  })
}
// 删除用户
export function delUser(data) {
  return request({
    url: '/user/delete',
    method: 'post',
    data
  })
}
// 重置密码
export function resetPwd(data) {
  return request({
    url: '/user/reset',
    method: 'post',
    data
  })
}
// 批量导入用户
export function importUser(data) {
  return request({
    url: '/user/import',
    method: 'post',
    data
  })
}
// 按邮箱查找用户
export function searchUser(data) {
  return request({
    url: '/user/userSearch',
    method: 'post',
    data
  })
}