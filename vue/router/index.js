import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/views/Home'
import Login from '@/views/Login'
import UserInfo from "@/views/UserInfo";

Vue.use(Router)

const router = new Router({
  routes: [
    {
      path: '/home',
      name: 'Home',
      component: Home,
    },
    {
      path: '/userinfo',
      name: 'UserInfo',
      component: UserInfo,
    },
    {
      path: '/login',
      name: 'Login',
      component: Login
    }
  ]
})
// router.beforeEach((to, from, next) => {
//
//   const token = sessionStorage.setItem("token")
//   //已经登录
//   if(token){
//     // 登录过就不能访问登录界面，需要中断这一次路由守卫login，执行下一次路由守卫/home
//     if (to.path === '/login') {
//       next({ path: '/home' })
//     }
//     // 动态添加路由
//     // 保存在store中路由不为空则放行 (如果执行了刷新操作，则 store 里的路由为空，此时需要重新添加路由)
//     if (store.getters.getRoutes.length || to.name != null) {
//       //放行
//       next()
//     } else {
//       // 将路由添加到 store 中，用来标记已添加动态路由
//       store.commit('SET_ROUTER', '需要添加的路由')
//       router.addRoutes('需要添加的路由')
//       // 如果 addRoutes 并未完成，路由守卫会一层一层的执行执行，直到 addRoutes 完成，找到对应的路由
//       next({ ...to, replace: true })
//       //replace: true只是一个设置，表示不能通过浏览器后退按钮返回前一个路由。
//     }
//   }else {
//     // 未登录时，注意 ：在这里也许你的项目不只有 logon 不需要登录 ，register 等其他不需要登录的页面也需要处理
//     if (to.path !== '/login') {
//       next({ path: '/login' })
//     } else {
//       next()
//     }
//   }
// })

export default router;
