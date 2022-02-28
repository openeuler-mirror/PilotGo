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
            store.dispatch('logOut')
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