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
 * @Date: 2022-02-18 17:47:56
 * @LastEditTime: 2022-03-04 15:59:02
 * @Description: provide agent log manager of pilotgo
 */
import { loginByEmail, logout } from '@/request/user'
import { getToken, setToken, removeToken, getUsername, setUsername, removeUsername, 
    getRoles, setRoles, removeRoles, removeUserDepartId, setUserDepartId,
    getUserDepartId, getUserDepartName, removeUserDepartName, setUserDepartName, } from '@/utils/auth'

const user = {
    state: {
        token: getToken(),
        username: getUsername(),
        roles: getRoles() ? JSON.parse(getRoles()) : [],
        departId: getUserDepartId(),
        departName: getUserDepartName(),
    },
    mutations: {
        SET_TOKEN: (state, token) => {
            state.token = token
        },
        SET_NAME: (state, name) => {
            state.username = name
        },
        SET_ROLES: (state, roles) => {
            state.roles = roles
        },
        SET_DEPARTID: (state, departId) => {
            state.departId = departId
        },
        SET_DEPARTNAME: (state, departName) => {
            state.departName = departName
        },
    },
    actions: {
        loginByEmail({ commit }, userInfo) {
            const username = userInfo.username.trim()
            return new Promise((resolve, reject) => {
                loginByEmail({'email':username, 'password':userInfo.password}).then(response => {
                    const res = response.data;
                    if (res.code != "200") {
                        reject(res)
                    } else {
                        commit('SET_TOKEN', res.data.token)
                        commit('SET_NAME', username)
                        // commit('SET_ROLES', username)
                        commit('SET_DEPARTID', res.data.departId)
                        commit('SET_DEPARTNAME', res.data.departName)
                        setToken(res.data.token)
                        // setRoles(res.data.roles)
                        setUsername(username)
                        setUserDepartId(res.data.departId)
                        setUserDepartName(res.data.departName)
                        resolve()
                    }
                }).catch(error => {
                    reject('缺少必要的参数')
                })
            })
        },
        logOut({ commit }, userInfo) {
            return new Promise((resolve, reject) => {
                logout().then(() => {
                    commit('SET_TOKEN', '')
                    commit('SET_ROLES', [])
                    commit('SET_MENUS', [])
                    commit('SET_NAME', '')
                    commit('SET_DEPARTID', '')
                    commit('SET_DEPARTNAME', '')
                    removeRoles();
                    removeUsername();
                    removeToken();
                    removeUserDepartId();
                    removeUserDepartName();
                    localStorage.clear()
                    resolve()
                }).catch(error => {
                    reject(error)
                })
            })
        },
    }
}

export default user