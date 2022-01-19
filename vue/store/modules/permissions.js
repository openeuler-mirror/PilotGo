import { constantRouterMap, routes } from '@/router'
// import { getPermission } from "@/api/role"
import { hasPermission } from "@/utils/auth";


function filterAsyncRouter(routers, menus) {
    routers.forEach((route) => {
        if (!hasPermission(menus, route)) {
            route.hidden = true;
        }
        route.children && filterAsyncRouter(route.children, menus)
    })
    return routers
}

const permission = {
    state: {
        routers: constantRouterMap,
        routes: routes,
        hostRouters: [],
        notfound: [],
        menus: [],
        operations: [],
        activePanel: ''
    },
    mutations: {
        SET_HOSTROUTERS: (state, routers) => {
            state.hostRouters = routers;
        },
        SET_MENUS: (state, menus) => {
            state.menus = menus
        },
        SET_OPERATIONS: (state, operations) => {
            state.operations = operations;
        },
        SET_LOGROUTERS: (state, routers) => {
            state.logRouters = routers;
        },
        SET_ACtiVEPANEL: (state, panel) => {
            state.activePanel = panel;
        },
        SET_NOTFOUND: (state, routers) => {
            state.notfound = routers;
        }
    },
    actions: {
        GenerateRoutes({ commit, state }) {
            return new Promise(resolve => {
                const menus = state.menus;
                let hostRouters;
                hostRouters = filterAsyncRouter(JSON.parse(JSON.stringify(routes)), menus)
                commit('SET_HOSTROUTERS', hostRouters)
                resolve()
            })
        },
       /*  getPermission({ commit }, ids) {
            return getPermission({ ids }).then(res => {
                return new Promise((resolve, reject) => {
                    if (res.data.code === "0") {
                        let data = res.data.data;
                        let { menu, operation } = data;
                        let arr = []
                        operation.forEach(item => {
                            arr.push(item.name)
                        })
                        commit("SET_MENUS", menu);
                        commit("SET_OPERATIONS", arr);
                        resolve()
                    } else {
                        reject()
                    }
                })
            })
        }, */
        SetMenus({ commit }, menus) {
            commit("SET_MENUS", menus)
        },
        SetActivePanel({ commit }, panel) {
            commit("SET_ACtiVEPANEL", panel)
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
        getOperations: state => {
            return state.operations
        },
        getPaths: state => {
            return state.routes[2].children.filter(item => {
                return item.meta != undefined;
            }).map(item => item.meta)
        }
    }
}

export default permission