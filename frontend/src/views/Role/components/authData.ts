/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
 */
// 整个系统需要控制权限的按钮配置项 当前24

// 需要做动态添加插件权限逻辑

export let authData = [
  {
    id: "1",
    label: "概览",
    isMenu: true,
    display: true,
    menuName: "overview",
    operations: [],
  },
  {
    id: "2",
    label: "系统",
    isMenu: true,
    display: true,
    menuName: "cluster",
    operations: [
      {
        id: "8",
        btnId: "1",
        label: "rpm下发",
        menuName: "rpm_install",
      },
      {
        id: "9",
        btnId: "2",
        label: "rpm卸载",
        menuName: "rpm_uninstall",
      },
      {
        id: "22",
        btnId: "15",
        label: "变更部门",
        menuName: "dept_change",
      },
      {
        id: "23",
        btnId: "16",
        label: "机器删除",
        menuName: "machine_delete",
      },
      {
        id: "24",
        btnId: "17",
        label: "创建批次",
        menuName: "batch_create",
      },
      {
        id: "25",
        btnId: "18",
        label: "添加部门",
        menuName: "dept_add",
      },
      {
        id: "26",
        btnId: "19",
        label: "删除部门",
        menuName: "dept_delete",
      },
      {
        id: "27",
        btnId: "20",
        label: "编辑部门",
        menuName: "dept_update",
      },
    ],
  },
  {
    id: "3",
    label: "批次",
    isMenu: true,
    display: true,
    menuName: "batch",
    operations: [
      {
        id: "10",
        btnId: "3",
        label: "编辑批次",
        menuName: "batch_update",
      },
      {
        id: "11",
        btnId: "4",
        label: "删除批次",
        menuName: "batch_delete",
      },
    ],
  },
  {
    id: "4",
    label: "用户管理",
    isMenu: true,
    display: true,
    menuName: "usermanager",
    operations: [
      {
        id: "12",
        btnId: "5",
        label: "添加用户",
        menuName: "user_add",
      },
      {
        id: "13",
        btnId: "6",
        label: "导入用户",
        menuName: "user_import",
      },
      {
        id: "14",
        btnId: "7",
        label: "编辑用户",
        menuName: "user_edit",
      },
      {
        id: "15",
        btnId: "8",
        label: "重置密码",
        menuName: "user_reset",
      },
      {
        id: "16",
        btnId: "9",
        label: "删除用户",
        menuName: "user_del",
      },
    ],
  },
  {
    id: "5",
    label: "角色管理",
    isMenu: true,
    display: true,
    menuName: "rolemanager",
    operations: [
      {
        id: "17",
        btnId: "10",
        label: "添加角色",
        menuName: "role_add",
      },
      {
        id: "18",
        btnId: "11",
        label: "编辑角色",
        menuName: "role_update",
      },
      {
        id: "19",
        btnId: "12",
        label: "删除角色",
        menuName: "role_delete",
      },
      {
        id: "20",
        btnId: "13",
        label: "角色授权",
        menuName: "role_modify",
      },
    ],
  },
  {
    id: "6",
    label: "自定义脚本管理",
    isMenu: true,
    display: true,
    menuName: "script",
    operations: [
      {
        id: "28",
        btnId: "21",
        label: "执行脚本",
        menuName: "run_script",
      },
      {
        id: "29",
        btnId: "22",
        label: "更新黑名单",
        menuName: "update_script_blacklist",
      },
    ],
  },
  {
    id: "7",
    label: "日志管理",
    isMenu: true,
    display: true,
    menuName: "audit",
    operations: [],
  },
  {
    id: "8",
    label: "插件管理",
    isMenu: true,
    display: true,
    menuName: "plugin",
    operations: [],
  },
  {
    id: "9",
    label: "监控告警",
    isMenu: true,
    display: false,
    menuName: "monitor",
    operations: [
      {
        id: "22",
        btnId: "15",
        label: "安装expoter",
        menuName: "expoter_install",
      },
      {
        id: "23",
        btnId: "16",
        label: "卸载expoter",
        menuName: "expoter_uninstall",
      },
    ],
  },
];
