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
 * @LastEditTime: 2022-03-04 11:13:58
 */
import router from './router'
import store from './store'
import { getToken } from '@/utils/auth'

const whiteList = ['/login'];

router.beforeEach((to, from, next) => {
    if (to.meta && to.meta.header_title) {
        document.title = to.meta.header_title
    }
    
    if(getToken()) {
        if(to.path === '/login') {
            next()
        } else {
            if(to.matched.length === 0) {
                next({ path: "/404", replace: true })
            } else {
                store.dispatch('SetActivePanel', to.meta.panel)
                next()
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