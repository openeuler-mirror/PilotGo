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
 * @Date: 2022-01-19 17:30:12
 * @LastEditTime: 2022-03-04 11:34:02
 * @Description: provide agent log manager of pilotgo
 */
import Cookies from 'js-cookie'

const TokenKey = 'Admin-Token'

const Username = 'Username'

const Roles = 'Roles'

const UserDepartId = "UserDepartId"

const UserDepartName = 'UserDepartName'


const whileList = [
    "/login",
    "/401",
    "/404"
]

export function getToken() {
    return Cookies.get(TokenKey)
}

export function setToken(token) {
    if (token) {
        return Cookies.set(TokenKey, token)
    }
}

export function removeToken() {
    return Cookies.remove(TokenKey)
}

export function getUsername() {
    return Cookies.get(Username)
}

export function setUsername(username) {
    if (username) {
        return Cookies.set(Username, username)
    }
}

export function removeUsername() {
    return Cookies.remove(Username)
}

export function getRoles() {
    return Cookies.get(Roles)
}

export function removeRoles() {
    return Cookies.remove(Roles)
}
export function setRoles(roles) {
    if (roles) {
        return Cookies.set(Roles, roles)
    }
}

export function getUserDepartId() {
    return Cookies.get(UserDepartId)
}

export function setUserDepartId(userDepartId) {
    if (userDepartId) {
        return Cookies.set(UserDepartId, userDepartId)
    }
}

export function removeUserDepartId() {
    return Cookies.remove(UserDepartId)
}

export function getUserDepartName() {
    return Cookies.get(UserDepartName)
}

export function setUserDepartName(userDepartName) {
    if(userDepartName) {
        return Cookies.set(UserDepartName, userDepartName)
    }
}

export function removeUserDepartName() {
    return Cookies.remove(UserDepartName)
}

export function hasPermission(menus, to) {
    if (whileList.includes(to.path)) return true
    if (!to.meta) return true
    return Array.isArray(menus) && menus.includes(to.meta.panel)
}