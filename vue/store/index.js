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
 * @Date: 2022-05-30 13:55:21
 * @LastEditTime: 2022-06-17 15:45:49
 */
import Vue from 'vue'
import Vuex from 'vuex'
import user from './modules/user'
import cluster from './modules/cluster'
import config from './modules/config'
import message from './modules/message'
import tagsView from './modules/tagsView'
import permissions from './modules/permissions'
import getters from './getters'

Vue.use(Vuex);

const store =  new Vuex.Store({
    modules: {
      user,
      permissions,
      cluster,
      config,
      tagsView,
      message,
    },
    getters,
})
export default store;