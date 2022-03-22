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
  { path: '/', component: _import('errorPage/404') },
]
export const routes = [
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
        meta: { title: 'overview', header_title: "概览", panel: "overview", icon_class: 'el-icon-s-help', icon_code:'&#xe612;' }
      },
      {
        path: '/cluster',
        name: 'Cluster',
        component:  _import('Cluster/Cluster'),
        meta: { title: 'cluster', header_title: "机器管理", panel: "cluster", icon_class: 'el-icon-s-platform' },
        children:[
          {
            path: '/firewall',
            name: 'Firewall',
            component:  _import('Firewall/Firewall'),
            meta: {  
              header_title: "防火墙配置", 
              panel: "cluster", 
              breadcrumb: [
                { name: '机器管理', path: '/cluster' },
                { name: '防火墙配置'}], 
            }
          },
          {
            path: '/cluster:uuid',
            name: 'MacDetail',
            component: _import('Cluster/detail/index'),
            meta: {
              header_title: "机器详情", 
              panel: "cluster", 
              breadcrumb: [
                  { name: '机器管理', path: '/cluster' },
                  { name: '机器详情'}
              ],
              icon_class: ''
            }
        },

        ]
      },
      {
        path: '/batch',
        component:  _import('Batch/Batch'),
        meta: { title: 'batch', header_title: "批次管理", panel: "batch", icon_class: 'el-icon-menu' },
        children:[
          {
            path: '/batch:id',
            name: 'BatchDetail',
            component: _import('Batch/detail/index'),
            meta: {
              header_title: "批次详情", 
              panel: "batch", 
              breadcrumb: [
                  { name: '批次管理', path: '/batch' },
                  { name: '批次详情'}
              ],
              icon_class: ''
            }
        },
        ]
      }, 
          
      /*{
        path: '/plug_in',
        name: 'PlugIn',
        component:  _import('Plug-in/Plug-in'),
        meta: { title: 'plug_in', header_title: "插件管理", panel: "plug_in", icon_class: 'el-icon-document' }
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
     /*  {
        path: '/firewall',
        name: 'Firewall',
        component:  _import('Firewall/Firewall'),
        meta: { title: 'firewall', header_title: "防火墙配置", panel: "firewall", icon_class: 'el-icon-s-home' }
      }, */
      {
        path: '/log',
        name: 'Log',
        component:  _import('Log/Log'),
        meta: { title: 'log', header_title: "日志管理", panel: "log", icon_class: 'el-icon-s-order' }
      },
      {
        path: '', 
        redirect: '/overview'
      },
    ]
  },
]

const router = new Router({
  mode: 'history',
  routes: [ ...routes, ...constantRouterMap],
})

export default router;
