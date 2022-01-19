import Vue from 'vue'
import Router from 'vue-router'
const _import = require('./_import')
import Home from '@/views/Home/Home'

Vue.use(Router)
const originalPush = Router.prototype.push
const originalReplace = Router.prototype.replace
Router.prototype.push = function push(location) {
    return originalPush.call(this, location).catch(err => err)
}
Router.prototype.replace = function replace(location) {
    return originalReplace.call(this, location).catch(err => err)
}
export const constantRouterMap = [
  { path: '/401', component: _import('errorPage/401') },
  { path: '/404', component: _import('errorPage/404') },
]
export const routes = [
  {
    path: '/', 
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: _import('Login'),
    meta: { title: 'login', header_title: "登录", panel: "login" }
  },
  {
    path: '/home',
    component: Home,
    children: [
      {
        path: '/overview',
        name: 'Overview',
        component: _import('Overview/Overview'),
        meta: { title: 'overview', header_title: "概览", panel: "overview", icon_class: 'el-icon-location' }
      },
      {
        path: '/cluster',
        name: 'Cluster',
        component:  _import('Cluster/Cluster'),
        meta: { title: 'cluster', header_title: "机器管理", panel: "cluster", icon_class: 'el-icon-s-platform' }
      },
      {
        path: '/cluster/ip',
        name: 'cluter_detail',
        component:  _import('Cluster/detail/index'),
      },
      {
        path: '/batch',
        name: 'Batch',
        component:  _import('Batch/Batch'),
        meta: { title: 'batch', header_title: "批次管理", panel: "batch", icon_class: 'el-icon-menu' }
      },     
      {
        path: '/plug_in',
        name: 'PlugIn',
        component:  _import('Plug-in/Plug-in'),
        meta: { title: 'plug_in', header_title: "插件管理", panel: "plug_in", icon_class: 'el-icon-document' }
      },
      /* {
        path: '/prometheus',
        name: 'Prometheus',
        component: Prometheus,
        // icon_class: 'el-icon-odometer'
      }, 
      {
        path: '/cockpit',
        name: 'Cockpit',
        component: Cockpit
        //icon_class: 'el-icon-setting'
      },*/
      {
        path: '/usermanager',
        name: 'UserManager',
        component:  _import('UserManager/UserMan'),
        meta: { title: 'usermanager', header_title: "用户管理", panel: "usermanager", icon_class: 'el-icon-user-solid' }
      },
      {
        path: '/rolemanager',
        name: 'RoleManager',
        component:  _import('RoleManager/RoleMan'),
        meta: { title: 'rolemanager', header_title: "角色管理", panel: "rolemanager", icon_class: 'el-icon-s-custom' }
      },
      {
        path: '/firewall',
        name: 'Firewall',
        component:  _import('Firewall/Firewall'),
        meta: { title: 'firewall', header_title: "防火墙配置", panel: "firewall", icon_class: 'el-icon-s-home' }
      },
      {
        path: '', 
        redirect: '/overview'
      },
    ]
  },
]

const router = new Router({
  mode: 'hash',
  routes: [...constantRouterMap, ...routes],
})

export default router;
