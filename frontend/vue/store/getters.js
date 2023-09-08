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
 * @LastEditTime: 2022-07-01 16:15:57
 */
const getters = {
  token: state => state.user.token,
  userName: state => state.user.username,
  roles: state => state.user.roles,
  userType: state => parseInt(state.user.userType),
  UserDepartId: state => parseInt(state.user.departId),
  UserDepartName: state => state.user.departName,
  activePanel: state => state.permissions.activePanel,
  iframeComponents: state => state.permissions.iframeComponents,
  selectIp: state => state.cluster.selectIp,
  tableTitle: state => state.cluster.tableTitle,
  immutable: state => state.cluster.immutable,
  visitedViews: state => state.tagsView.visitedViews,
  cachedViews: state => state.tagsView.cachedViews,
  messageNum: state => state.message.messageNum,
  agentPath: state => state.config.agentPath,
  configType: state => state.config.configType,
}

export default getters