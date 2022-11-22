<!-- 
  Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
  PilotGo is licensed under the Mulan PSL v2.
  You can use this software accodring to the terms and conditions of the Mulan PSL v2.
  You may obtain a copy of Mulan PSL v2 at:
      http://license.coscl.org.cn/MulanPSL2
  THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND, 
  EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
  See the Mulan PSL v2 for more details.
  Author: zhaozhenfang
  Date: 2022-01-19 17:30:12
  LastEditTime: 2022-05-18 17:43:56
  Description: provide agent log manager of pilotgo
 -->
<template>
  <div>
    <template v-for="item in routes">
      <router-link :to="item.path" :key="item.meta.header_title" v-if="!item.meta.hidden && !item.meta.submenu">
        <el-menu-item :index="item.meta.panel">
          <em :class="item.meta.icon_class"></em>
          <template #title>
            <span>{{ item.meta.header_title }}</span>
          </template>
        </el-menu-item>
      </router-link>

      <el-submenu popper-append-to-body :index="item.meta.panel" :key="item.meta.header_title"
        v-if="!item.meta.hidden && item.meta.submenu">
        <template #title><em :class="item.meta.icon_class"></em>
          <span>{{ item.meta.header_title }}</span>
        </template>
        <router-link v-for="subItem in item.submenu" :key="subItem.menuName" :to="subItem.name">
          <el-menu-item :index="subItem.name">
            <span>{{ subItem.menuName }}</span>
          </el-menu-item>
        </router-link>
      </el-submenu>
    </template>
  </div>
</template>

<script>
export default {
  name: "SidebarItem",
  props: {
    routes: {
      type: Array,
    },
  },
};
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
em {
  margin-right: 12px;
}

.el-menu-item {
  display: flex;
  align-items: center;
  height: 46px;
  margin: 10px 0;
  color: rgb(11, 35, 117);
}

.el-menu-item:hover {
  color: rgb(241, 139, 14);
}

.el-menu-item.is-active {
  background: #f2f8ff !important;
  border-right: 3px solid rgb(241, 139, 14);
  color: rgb(241, 139, 14);
}

.el-submenu {
  .el-menu-item {
    display: flex;
    align-items: center;
    margin: 0;
    font-size: 12px;
  }
}

a {
  text-decoration: none;
}

.router-link-active {
  text-decoration: none;
}
</style>