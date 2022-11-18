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
 * @LastEditTime: 2022-06-27 15:20:32
 */
import { constantRouterMap, routes } from '@/router'
import router from '@/router';
import { getPermission } from "@/request/user"
import { hasPermission } from "@/utils/auth";

function filterAsyncRouter(routers, menus) {
  routers.forEach((route) => {
    if (!hasPermission(menus, route)) {
      route.meta.hidden = true;
    }
    route.children && filterAsyncRouter(route.children, menus)
  })
  return routers
}

const permission = {
  state: {
    // 最终路由列表
    routers: [],
    // 动态路由列表
    dynamicRoutes: [],
    // router组件routes
    routes: routes,
    notfound: [],
    // 基本页面菜单面板
    baseMenus: [],
    // 最终页面菜单列表
    menus: [],
    //
    operations: [],
    // 当前激活面板
    activePanel: ''
  },
  mutations: {
    SET_ROUTERS: (state, routers) => {
      state.routers = [...routers, ...constantRouterMap, ...state.dynamicRoutes];
    },
    SET_DYNAMIC_ROUTERS: (state, routers) => {
      state.dynamicRoutes = routers;
    },
    SET_MENUS: (state, menus) => {
      state.menus = menus
    },
    SET_BASE_MENUS: (state, menus) => {
      state.baseMenus = menus
    },
    SET_OPERATIONS: (state, operations) => {
      state.operations = operations;
    },
    SET_LOGROUTERS: (state, routers) => {
      state.logRouters = routers;
    },
    SET_ACTIVE_PANEL: (state, panel) => {
      state.activePanel = panel;
    },
    SET_NOTFOUND: (state, routers) => {
      state.notfound = routers;
    }
  },
  actions: {
    GenerateRoutes({ commit, state }) {
      return new Promise(resolve => {
        let menus = [...state.baseMenus];
        state.dynamicRoutes.forEach(() => {
          menus.push("plugin3");
        })
        commit("SET_MENUS", menus)

        router.updateRoutes(state.dynamicRoutes);
        let routers = filterAsyncRouter(JSON.parse(JSON.stringify(routes)), menus)
        commit('SET_ROUTERS', routers)

        resolve()
      })
    },
    SetDynamicRouters({ commit }, routers) {
      commit("SET_DYNAMIC_ROUTERS", routers)
    },
    getPermission({ commit }, roles) {
      let roleId = roles.split(',').map(Number)
      return getPermission({ roleId: roleId }).then(res => {
        return new Promise((resolve, reject) => {
          if (res.data.code === 200) {
            let data = res.data.data;
            let { menu, button } = data;
            button.push("default_all");
            commit("SET_BASE_MENUS", menu.split(','));
            commit("SET_OPERATIONS", button);
            resolve()
          } else {
            reject()
          }
        })
      })
    },
    SetMenus({ commit }, menus) {
      commit("SET_MENUS", menus)
    },
    SetBaseMenus({ commit }, menus) {
      commit("SET_BASE_MENUS", menus)
    },
    SetActivePanel({ commit }, panel) {
      commit("SET_ACTIVE_PANEL", panel)
    },
    addRoute({ commit, state }, route) {
      let r = filterAsyncRouter(JSON.parse(JSON.stringify(routes)), state.menus)
      r[1].children.push(...state.dynamicRoutes);
      commit('SET_ROUTERS', r)
    }
  },
  getters: {
    addRoutes: state => {
      const { hostRouters, notfound } = state;
      return [...hostRouters, ...notfound];
    },
    getMenus: state => {
      return state.menus
    },
    getBaseMenus: state => {
      return state.baseMenus
    },
    getOperations: state => {
      return state.operations
    },
    getPaths: state => {
      return state.routers[1].children.filter(item => {
        return item.meta != undefined;
      }).map(item => item)
    },
    getDynamicRoutes: state => {
      return state.dynamicRoutes
    }
  }
}

export default permission