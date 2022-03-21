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
 * @LastEditTime: 2022-03-16 15:20:31
 */
import { request } from './request'
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