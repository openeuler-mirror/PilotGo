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
 * @Date: 2022-02-25 16:33:46
 * @LastEditTime: 2022-03-18 14:52:50
 */
import router from './router'
import store from './store'
import { getRoles, hasPermission } from '@/utils/auth'

const whiteList = ['/login']

router.beforeEach((to, from, next) => {
    if (to.meta && to.meta.header_title) {
        document.title = to.meta.header_title
    }
    if (getRoles()) {
        if (to.path === '/login') {
            next({ path: '/' })
            
        } else {
            if (!store.getters.getMenus || store.getters.getMenus.length === 0) {
                store.dispatch('getPermission', store.getters.roles).then(res => {
                    store.dispatch('GenerateRoutes').then(() => {
                        next({ ...to, replace: true })
                    })
                })
            } else {
                if (to.path === "/") {
                    let paths = store.getters.getPaths.filter(path => !path.hidden)
                    let to = paths.length > 0 ? paths[0].panel : "/404"
                    next({ path: to, replace: true })
                } else {
                    if (hasPermission(store.getters.getMenus, to)) {
                        store.dispatch('SetActivePanel', to.meta.panel)
                        next()
                    } else {
                        next({ path: '/404', replace: true })
                        
                    }
                }
            }
        }
    } else {
        if (whiteList.indexOf(to.path) !== -1) {
            next()
        } else {
            next('/login')
            
        }
    }
})

router.afterEach(route => {
    
})