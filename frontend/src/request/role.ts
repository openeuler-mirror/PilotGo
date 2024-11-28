/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
 */
import request from './request'

// 分页获取角色列表
export function getRolesPaged(data: any) {
    return request({
        url: '/user/roles_paged',
        method: 'get',
        params: data
    })
}

// 获取所有角色列表
export function getRoles() {
    return request({
        url: '/user/roles',
        method: 'get',
    })
}

// 修改角色权限
export function changeRolePermission(data: any) {
    return request({
        url: '/user/roleChange',
        method: 'post',
        data
    })
}

// 删除角色
export function deleteRole(data: any) {
    return request({
        url: '/user/delRole',
        method: 'post',
        data
    })
}

// 添加角色
export function addRole(data: any) {
    return request({
        url: '/user/addRole',
        method: 'post',
        data
    })
}

// 更新角色信息
export function updateRole(data: any) {
    return request({
        url: '/user/updateRole',
        method: 'post',
        data
    })
}