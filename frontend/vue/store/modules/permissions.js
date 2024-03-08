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
 * @LastEditTime: 2023-09-02 11:09:25
 */
import Vue from 'vue';
import { constantRouterMap, routes } from '@/router'
import router from '@/router';
import { getPermission } from "@/request/user"
import { hasPermission } from "@/utils/auth";
import { getPlugins } from "@/request/plugin";
import _import from '../../router/_import';

// 过滤有展示权限的路由
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
    activePanel: '',
    // iframe的组件数组
    iframeComponents: [],
  },
  mutations: {
    SET_ROUTERS: (state, routers) => {
      state.routers = [...routers, ...constantRouterMap, ...state.dynamicRoutes];
    },
    SET_DYNAMIC_ROUTERS: (state, routers) => {
      state.dynamicRoutes = routers;
    },
    ADD_DYNAMIC_ROUTERS: (state, newRoute) => {
      state.dynamicRoutes = [...state.dynamicRoutes, ...newRoute];
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
        state.dynamicRoutes.forEach((item, index) => {
          menus.push('plugin' + index);
        })
        commit("SET_MENUS", menus)

        router.updateRoutes(state.dynamicRoutes);
        let routers = filterAsyncRouter(JSON.parse(JSON.stringify(routes)), menus)
        commit('SET_ROUTERS', routers)
        resolve()
      })
    },
    SetDynamicRouters({ commit, state }, routers) {
      // if (routers.length == 0) {
      // 初始化动态路由表
      return new Promise(resolve => {
        // 获取动态插件路由
        let p = [];
        state.iframeComponents = []
        getPlugins().then((res) => {
          if (res.data.code === 200) {
            res.data.data.forEach((item, index) => {
              if (item.enabled === 0) {
                // 0:禁用，1：启用
                return;
              }
              p.push({
                path: '/plugin' + index,
                name: 'Plugin' + index,
                iframeComponent: '',
                meta: {
                  title: 'plugin', header_title: item.name, panel: "plugin" + index, icon_class: 'el-icon-s-ticket', url: item.url,
                  breadcrumb: [
                    { name: item.name },
                  ],
                }
              })
              let iframeObj = {
                path: '/plugin' + index,
                name: 'Plugin' + index,
                src: '/plugin/' + item.name,
                component: _import('IFrame/IFrame'), // 组件文件的引用
                url: item.url,
                plugin_type: item.plugin_type
              }
              state.iframeComponents.push(iframeObj);
              Vue.component('Plugin' + index, _import('IFrame/IFrame'));
            });
            commit("SET_DYNAMIC_ROUTERS", p)
            resolve()
          } else {
            this.$message.error("查询插件列表错误：", res.data.msg);
          }
        })

      })
      // } else {
      //   return new Promise(resolve => {
      //     commit("ADD_DYNAMIC_ROUTERS", routers)
      //     resolve()
      //   })
      // }
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