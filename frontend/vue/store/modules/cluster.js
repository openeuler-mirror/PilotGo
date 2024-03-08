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
 * @LastEditTime: 2022-04-13 10:20:51
 */
const cluster = {
    state: {
        selectIp: '',
        tableTitle: '',
        immutable: false,
    },
    mutations: {
        SET_SELECTIP(state, ip) {
            state.selectIp = ip;
        },
        SET_TABLETITLE(state,title) {
            state.tableTitle = new String(title);
        },
        SET_IMMUTABLE(state, yes) {
            state.immutable = yes;
        },
    },
    actions: {
        setSelectIp({commit}, ip) {
            commit('SET_SELECTIP',ip + ':9100');
        },
        setTableTitle({commit}, title) {
            commit('SET_TABLETITLE', title);
        }
    },
    
}

export default cluster;