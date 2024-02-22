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
  {
    path: '/:catchAll(.*)',
    redirect: '/404'
  }
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

import { ref, onMounted, watchEffect, shallowRef } from "vue";
import { routerStore, type Menu } from '@/stores/router';
import { iframeComponents } from "@/views/Plugin/plugin";
import PluginFrame from "@/views/Plugin/PluginFrame.vue";

import { app } from "@/main";
import { useRouter } from "vue-router";

export function updateSidebarItems() {
  let menus = generateLocalMenus();
  for (let item of iframeComponents.value) {
    // console.log(item)
    // 添加插件路由信息
    addPluginRoute(item);
    // 更新侧边菜单
    let obj: Menu = {
      path: item.path,
      title: item.name,
      hidden: false,
      panel: item.name,
      icon: "Menu",
      subMenus: null,
    }
    let subMenus = [] as any;
    if (item.subMenus.length > 0) {
      item.subMenus.forEach((subItem: any) => {
        subMenus.push({
          hidden: false,
          icon: "",
          panel: subItem.subRoute,
          title: subItem.title,
          path: item.path + subItem.subRoute,
        })
      })
    }
    obj.subMenus = subMenus;
    menus.push(obj)

    // app.component(item.name, PluginFrame);

  }
  routerStore().menus = menus;
}

function addPluginRoute(item: any) {
  let childrens = [] as any;
  if (item.subMenus && item.subMenus.length > 0) {
    // 如果有多级菜单
    item.subMenus.forEach((subItem: any) => {
      childrens.push({
        path: item.path + subItem.subRoute,
        component: () => import('@/views/Plugin/PluginFrame.vue'),
        meta: {
          title: subItem.title,
          subRoute: subItem.path + subItem.subRoute,
          breadcrumb: []
        }
      })
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
      title: item.name
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

router.beforeEach((to, from) => {
  if (to.meta && to.meta.title) {
    document.title = to.meta.title as string
  }
})

export function directTo(to: any) {
  router.push(to)
}