import router from './router'
import store from './store'
import { getToken } from '@/utils/auth'

router.beforeEach((to, from, next) => {
    if(to.path === '/login') {
        store.dispatch('logOut')
        next()
    } else {
        if(getToken()) {
            store.dispatch('SetActivePanel', to.meta.panel)
            document.title = to.meta.header_title
            next()
        } else {
            next('/')
        }
        
    }
})