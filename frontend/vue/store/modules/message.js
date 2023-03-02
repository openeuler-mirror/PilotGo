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
 * @Date: 2022-05-25 17:30:12
 * @LastEditTime: 2022-05-26 09:43:27
 */
const message = {
  state: {
      messageNum: 0,
  },
  mutations: {
      SET_MessageNum(state, count) {
          state.messageNum = count;
      },
  },
  actions: {
      setMessageNum({commit}, count) {
          commit('SET_MessageNum', count);
      },
  },
  
}

export default message;