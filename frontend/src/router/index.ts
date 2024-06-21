import { createRouter, createWebHistory } from 'vue-router';

const commonRoutes = [
  {
    path: '/',
    redirect: '/home',
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/Login/Login.vue'),
  },
  {
    path: '/404',
    name: 'errorPage',
    component: () => import('@/views/errorPage.vue')
  },
  /* {
    path: '/:catchAll(.*)',
    redirect: '/404'
  } */
];

let sidebarRoutes = [
  {
    path: '/home',
    name: 'home',
    redirect: '/overview',
    component: () => import('@/views/Home/Home.vue'),
    children: [
      {
        path: '/overview',
        name: 'overview',
        component: () => import('@/views/Overview/Overview.vue'),
        meta: {
          title: '概览',
          panel: 'overview',
          icon: 'HomeFilled',
          breadcrumb: [{ name: '概览' }],
        },
      },
      {
        path: '/cluster',
        meta: { title: "系统", panel: "cluster", icon: 'Platform' },
        children: [
          {
            path: '',
            redirect: '/cluster/macList'
          },
          {
            path: '/cluster/macList',
            name: 'macList',
            component: () => import('../views/Cluster/Cluster.vue'),
            meta: {
              title: "机器列表",
              panel: "/cluster/macList",
              breadcrumb: [
                {
                  name: '系统', path: '/cluster', children: [
                    { name: 'createBatch', menuName: '创建批次' },
                  ]
                },
                { name: '机器列表' },
              ],
              icon: '',
            }
          },
          {
            path: '/cluster/machine/:uuid',
            name: 'MacDetail',
            component: () => import('../views/Cluster/MachineDetail/Index.vue'),
            meta: {
              title: "机器详情",
              panel: "/cluster/macList",
              breadcrumb: [
                {
                  name: '系统', path: '/cluster', children: [
                    { name: 'createBatch', menuName: '创建批次' },
                  ]
                },
                { name: '机器列表', path: '/cluster/' },
                { name: '机器详情' }
              ],
              icon: '',
              ignore: true,
            }
          },
          {
            path: '/cluster/createBatch',
            name: 'createBatch',
            component: () => import('../views/Cluster/CreateBatch.vue'),
            meta: {
              title: "创建批次",
              panel: "/cluster/createBatch",
              breadcrumb: [
                {
                  name: '系统', path: '/cluster', children: [
                    { name: 'macList', menuName: '机器列表' },
                  ]
                },
                { name: '创建批次' },
              ],
              icon: ''
            }
          }
        ]
      },
      {
        path: '/batch',
        meta: {
          title: "批次", panel: "batch", icon: 'DocumentCopy',
          breadcrumb: [
            { name: '批次' },
          ]
        },
        children: [
          {
            path: '',
            redirect: '/batch/list'
          },
          {
            path: '/batch/list',
            name: 'BatchList',
            component: () => import('../views/Batch/Batch.vue'),
            meta: {
              title: "批次列表",
              panel: "batch",
              breadcrumb: [
                { name: '批次', path: '/batch' },
                { name: '批次列表' }
              ],
              icon: ''
            },
          },
          {
            path: '/batch/detail/:id',
            name: 'BatchDetail',
            component: () => import('../views/Batch/Detail.vue'),
            meta: {
              ignore: true,
              title: "批次详情",
              panel: "batch",
              breadcrumb: [
                { name: '批次', path: '/batch' },
                { name: '批次详情' }
              ],
              icon: ''
            }
          },
        ]
      },
      {
        path: '/user',
        name: 'User',
        component: () => import('../views/User/User.vue'),
        meta: {
          title: "用户管理", panel: "user", icon: 'UserFilled',
          breadcrumb: [
            { name: '用户管理' },
          ],
        }
      },
      {
        path: '/role',
        name: 'Role',
        component: () => import('../views/Role/Role.vue'),
        meta: {
          title: "角色管理", panel: "role", icon: 'Ticket',
          breadcrumb: [
            { name: '角色管理' },
          ],
        }
      },
      {
        path: '/audit',
        name: 'Audit',
        component: () => import('../views/Audit/Audit.vue'),
        meta: {
          title: "审计日志", panel: "audit", icon: 'View',
          breadcrumb: [
            { name: '审计日志' },
          ],
        }
      },
      {
        path: '/pgPlugin',
        name: 'Plugin',
        component: () => import('../views/Plugin/Plugin.vue'),
        meta: {
          title: "插件管理", panel: "plugin", icon: 'Menu',
          breadcrumb: [
            { name: '插件管理' },
          ],
        }
      }
    ],
  },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [...commonRoutes, ...sidebarRoutes],
});

export default router;

import { routerStore, type Menu } from '@/stores/router';
import { iframeComponents } from "@/views/Plugin/plugin";

import { hasPermisson } from '@/module/permission';
import { getToken } from "@/module/cookie";

export function updateSidebarItems() {
  let menus = generateLocalMenus();
  routerStore().routers = [];
  for (let item of iframeComponents.value) {
    // 添加插件路由信息
    addPluginRoute(item);
    // 更新侧边菜单
    let obj: Menu = {
      type:'plugin',
      path: item.path,
      title: item.name,
      hidden: false,
      panel: item.name,
      icon: "Menu",
      subMenus: null,
    }
    let subMenus = [] as any;
    if (item.subMenus && item.subMenus.length > 0) {
      item.subMenus.forEach((subItem: any) => {
        subMenus.push({
          hidden: false,
          icon: "",
          panel: subItem.url,
          title: subItem.name,
          path: item.path + subItem.url,
        })
      })
    }
    obj.subMenus = subMenus;
    menus.push(obj)

    // app.component(item.name, PluginFrame);

  }
  menus.map((menu: any) => {
    if (menu.type === 'plugin') return; // 插件权限先不做限制
    menu.hidden = !hasPermisson('menu/' + menu.panel);
  })
  routerStore().menus = menus;
}

function addPluginRoute(item: any) {
  routerStore().routers.push(item);
  let childrens = [] as any;
  // 插件至少两级子菜单
  if (item.subMenus && item.subMenus.length > 0) {
    let breadChild = [] as any;
    item.subMenus.forEach((subItem: any) => {
      // 1.创建breadcrumb的children
      breadChild.push({
        path: item.path + subItem.url,
        menuName: subItem.name
      })
      // 2.创建router -> children
      childrens.push({
        path: item.path + subItem.url,
        component: () => import('@/views/Plugin/PluginFrame.vue'),
        meta: {
          title: subItem.name,
          subRoute: '/plugin/' + item.name.split('-')[1] + subItem.url,
          // subRoute: '/plugin/' + item.name.split('plugin-')[1] + subItem.url,
          breadcrumb: [{
            name: item.name,
            path: item.path
          }, {
            name: subItem.name
          }]
        }
      })
    })

    // 3.给子router添加breadCrumb的children
    childrens.forEach((childRoute: any, index: number) => {
      childRoute.meta.breadcrumb[0].children = breadChild.filter((breabItem: any, bIndex: number) => index !== bIndex)
    })
  }
  router.addRoute("home", {
    path: item.path,
    name: item.name,
    component: () => import('@/views/Plugin/PluginFrame.vue'),
    meta: {
      icon: "Menu",
      name: item.name,
      url: item.url,
      plugin_type: item.plugin_type,
      title: item.name,
      breadcrumb: [{
        name: item.name
      }]
    },
    children: childrens
  })
}

function generateLocalMenus() {
  // 迭代 /home 下的所有组件
  let menus = [];
  for (let route of sidebarRoutes[0].children) {
    let subMenus = []
    if (route.children != null) {
      for (let item of route.children) {
        if (item.meta != null) {
          if ('ignore' in item.meta && item.meta.ignore === true) {
            continue
          }
          let obj: Menu = {
            path: item.path,
            title: item.meta.title,
            hidden: false,
            panel: item.meta.panel,
            icon: item.meta.icon,
            subMenus: null,
          }
          subMenus.push(obj)
        }
      }
    }

    let obj: Menu = {
      path: route.path,
      title: route.meta.title,
      hidden: false,
      panel: route.meta.panel,
      icon: route.meta.icon,
      subMenus: subMenus.length > 0 ? subMenus : null,
    };

    menus.push(obj);
  }

  return menus;
}

// 路由导航
const whiteList = ['/login']
router.beforeEach((to, from, next) => {
  if (to.meta && to.meta.title) {
    document.title = to.meta.title as string
  }
  if (getToken()) {
    // 已登录
    if (to.path.includes('plugin-')) {
      // 解决在插件页面一刷新页面空白问题
      let paths = [] as any;
      router.getRoutes().forEach(routeItem => {
        paths.push(routeItem.path);
      })
      if (paths.includes(to.path)) {
        next();
      } else {
        new Promise(async (resolve, rejection) => {
          await routerStore().routers.forEach((route: any) => {
            addPluginRoute(route);
          })
          resolve('addRoute success');
        }).then(res => {
          next({ ...to, replace: true })
        })
    }
    } else {
      next()
    }
  
  } else {
    // 未登录
    if (whiteList.indexOf(to.path) !== -1) {
      // 如果在白名单的话，放行
      next()
    } else {
      // 如果不在白名单的话，转login登录
      next('/login')
    }
  }
  
})

export function directTo(to: any) {
  router.push(to)
}