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
        meta: { title: 'overview', header_title: "概览", panel: "overview", icon_class: 'el-icon-s-data', submenu: false ,
        breadcrumb: [
          { name: '概览' },
        ]}
      },
      {
        path: '/cluster',
        component: _import('Cluster/Cluster'),
        meta: { title: 'cluster', header_title: "系统", panel: "cluster", icon_class: 'el-icon-s-platform', submenu: true},
        children:[
          {
            path: '',
            redirect: 'macList'
          },
          {
            path: '/cluster/macList',
            name: 'macList',
            component:  _import('Cluster/Cluster'),
            meta: {
              header_title: "机器列表", 
              panel: "/cluster/macList",
              breadcrumb: [
                { name: '系统', path: '/cluster', children: [
                  {name:'createBatch',menuName:'创建批次'},
                  {name:'Prometheus',menuName:'性能监控'},
                ]},
                { name: '机器列表' },
              ],
              icon_class: ''
            },
            children: [
              {
                path: '/cluster/macList/:detail',
                name: 'MacDetail',
                component:_import('Cluster/detail/index'),
                meta: {
                  header_title: "机器详情", 
                  panel: "/cluster/macList", 
                  breadcrumb: [
                      { name: '系统', path: '/cluster', children:[
                        {name:'createBatch',menuName:'创建批次'},
                        {name:'prometheus',menuName:'性能监控'},
                      ]},
                      { name: '机器列表', path: '/cluster/' },
                      { name: '机器详情'}
                  ],
                  icon_class: ''
                }
              },
            ]
          },
          {
            path: '/cluster/createBatch',
            name: 'createBatch',
            component:_import('Cluster/createBatch/index'),
            meta: {
              header_title: "创建批次", 
              panel: "/cluster/createBatch", 
              breadcrumb: [
                  { name: '系统', path: '/cluster', children:[
                    {name:'macList',menuName:'机器列表'},
                    {name:'Prometheus',menuName:'性能监控'},
                  ]},
                  { name: '创建批次'}
              ],
              icon_class: ''
            }
          },
          {
            path: '/cluster/prometheus',
            name: 'Prometheus',
            component: _import('Prometheus/index'),
            meta: {  
              header_title: "性能监控", 
              panel: "/cluster/prometheus", 
              breadcrumb: [
                { name: '系统', path: '/cluster', children:[
                  {name:'macList',menuName:'机器列表'},
                  {name:'createBatch',menuName:'创建批次'},
                ]},
                { name: '性能监控'}
            ],
              icon_class: ''
            },
          },
          
        ],
        submenu: [
          {
            name: '/cluster/macList',
            menuName: "机器列表", 
            panel: 'cluster'
          },
          {
            name: '/cluster/createBatch',
            menuName: "创建批次", 
            panel: 'cluster'
          },
          {
            name: '/cluster/prometheus',
            menuName: "性能监控", 
            panel: 'cluster'
          },
        ]
      },
      {
        path: '/batch',
        name: 'Batch',
        component:  _import('Batch/Batch'),
        meta: { title: 'batch', header_title: "批次", panel: "batch", icon_class: 'el-icon-menu',breadcrumb: [
          { name: '批次' },
        ] },
        children:[
          {
            path: '/batch:id',
            name: 'BatchDetail',
            component: _import('Batch/detail/index'),
            meta: {
              header_title: "批次详情", 
              panel: "batch", 
              breadcrumb: [
                  { name: '批次', path: '/batch' },
                  { name: '批次详情'}
              ],
              icon_class: ''
            }
        },
        ]
      }, 
      {
        path: '/usermanager',
        name: 'UserManager',
        component:  _import('UserManager/UserMan'),
        meta: { title: 'usermanager', header_title: "用户管理", panel: "usermanager", icon_class: 'el-icon-user-solid' ,
        breadcrumb: [
          { name: '用户管理' },
        ],}
      },
      {
        path: '/rolemanager',
        name: 'RoleManager',
        component:  _import('RoleManager/RoleMan'),
        meta: { title: 'rolemanager', header_title: "角色管理", panel: "rolemanager", icon_class: 'el-icon-s-custom' ,
        breadcrumb: [
          { name: '角色管理' },
        ],}
      },
      {
        path: '/config',
        component:  _import('Config/libfile'),
        meta: { title: 'config', header_title: "配置管理", panel: "config", icon_class: 'el-icon-s-tools' , submenu: true,
        breadcrumb: [
          { name: '配置管理'}
        ],},
        children:[
          {
            path: '',
            redirect: 'libconfig'
          },
          {
            path: '/config/libconfig',
            name: 'libconfig',
            component:_import('Config/libfile'),
            meta: {
              header_title: "库配置文件", 
              panel: "/config/libconfig", 
              breadcrumb: [
                  { name: '配置管理', path: '/config/libconfig'},
                  { name: '库配置文件'}
              ],
              icon_class: ''
            }
          },
        ],
        submenu: [
          {
            name: '/config/libconfig',
            menuName: "库配置文件", 
          },
        ]

      },
      {
        path: '/log',
        name: 'Log',
        component:  _import('Log/Log'),
        meta: { title: 'log', header_title: "日志管理", panel: "log", icon_class: 'el-icon-s-order' ,
        breadcrumb: [
          { name: '日志管理' },
        ],}
      },
      {
        path: '/plugin',
        name: 'Plugin',
        component:  _import('Plugin/Plugin'),
        meta: { title: 'plugin', header_title: "插件管理", panel: "plugin", icon_class: 'el-icon-s-order' ,
        breadcrumb: [
          { name: '插件管理' },
        ],}
      },
    ]
  },
]

const router = new Router({
  mode: 'history',
  routes: [ ...routes, ...constantRouterMap],
})

export default router;
