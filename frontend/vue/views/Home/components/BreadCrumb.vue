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
  Date: 2022-02-25 16:33:45
  LastEditTime: 2022-05-18 11:02:19
  Description: provide agent log manager of pilotgo
 -->
<template>
  <div>
    <el-breadcrumb separator="/">
      <el-breadcrumb-item to="/overview">
        <em class="el-icon-s-home"></em>
          <span class="el-dropdown-link">
            首页
          </span>
      </el-breadcrumb-item>
        
      <template v-for="item in $route.meta.breadcrumb" >
        <el-breadcrumb-item
          :key="item.name"
          v-if="item.path && !item.hidden && item.children"
        >
          <el-dropdown @command="handleMenuCommand">
            <span class="el-dropdown-link">
              {{ item.name }}<em class="el-icon-arrow-down el-icon--right"></em>
            </span>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item v-for="route in item.children" v-if="!route.hidden" :key="route.name" :command="route.name">
                {{ route.menuName }}
              </el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </el-breadcrumb-item>
        <el-breadcrumb-item
          :key="item.name"
          v-if="item.path && !item.hidden && !item.children"
          :to="{ path:item.path }"
        >
          {{ item.name }}
        </el-breadcrumb-item>
        <el-breadcrumb-item
          :key="item.name"
          v-if="!item.path && !item.hidden"
        >
          {{ item.name }}
        </el-breadcrumb-item>
      </template>
    </el-breadcrumb> 
  </div>
</template>

<script>
export default {
  name: "BreadCrumb",
  data() {
    return {

    }
  },
  methods: {
    handleMenuCommand(command) {
      this.$router.push({
        name: command,
      })
    }
  }
};
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
.el-breadcrumb {
  font-size: 15px;
  .el-breadcrumb__item:last-child {
    .el-breadcrumb__inner {
      color: rgb(11, 35, 117)
    }
  }
}
</style>