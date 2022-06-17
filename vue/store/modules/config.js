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
 * @LastEditTime: 2022-06-17 15:59:16
 */
import { libFileList } from "@/request/config"
const cluster = {
    state: {
        agentPath: [],
        configType: [],
    },
    mutations: {
        SET_AGENTPATH(state, path) {
            state.agentPath = path;
        },
        SET_CONFIGTYPE(state,type) {
            state.configType = type;
        },
    },
    actions: {
        getLibFileInfo({ commit },params) {
          libFileList(params).then(res => {
            if(res.data.code === 200) {
              let {agentPath, type} = res.data;
              commit('SET_AGENTPATH',agentPath);
              commit('SET_CONFIGTYPE', type);
            }
          })
        },
    },
    
}

export default cluster;