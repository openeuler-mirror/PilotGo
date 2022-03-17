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
  LastEditTime: 2022-03-14 20:30:59
  Description: provide agent log manager of pilotgo
 -->
<template>
  <div>
    <template v-for="item in routes">
      <router-link
        :to="item.panel"
        :key="item.header_title"
        v-if="!item.hidden"
      >
        <el-menu-item :index="item.panel">
          <em :class="item.icon_class"></em>
          <template #title>{{ item.header_title }}</template>
        </el-menu-item>
      </router-link>
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
  height: 46px;
  line-height: 46px;
  margin: 10px 0;
}
.el-menu-item.is-active {
  background: #f2f8ff !important;
  border-right: 3px solid #0076ff;
  color: #0076ff;
}
a {
  text-decoration: none;
}
 
.router-link-active {
  text-decoration: none;
}
</style>