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
  LastEditTime: 2022-02-25 17:40:04
  Description: provide agent log manager of pilotgo
 -->
<template>
  <div>
    <el-breadcrumb
      separator=">"
      v-if="$route.meta.breadcrumb"
    >
      <template v-for="item in $route.meta.breadcrumb">
        <el-breadcrumb-item
          :key="item.name"
          v-if="item.path"
          :to="{ path:item.path }"
        >
          {{ item.name }}
        </el-breadcrumb-item>
        <el-breadcrumb-item
          :key="item.name"
          v-else
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
};
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
.el-breadcrumb {
    font-size: 16px;
    margin-bottom: 12px;
  }
</style>