<!--
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
-->
<template>
  <div>
    <el-breadcrumb separator="/" class="bread">
      <el-breadcrumb-item to="/overview">
        <span class="el-dropdown-link"> 首页 </span>
      </el-breadcrumb-item>

      <template v-for="item in route.meta.breadcrumb">
        <el-breadcrumb-item :key="item.name" v-if="item.path && !item.hidden && item.children">
          <el-dropdown @command="router2path">
            <span class="el-dropdown-link">
              {{ item.name }}
              <el-icon class="el-icon--right">
                <arrow-down />
              </el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item
                  v-for="route in item.children"
                  v-if="!route.hidden"
                  :key="route.name"
                  :command="route"
                >
                  {{ route.menuName }}
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </el-breadcrumb-item>

        <el-breadcrumb-item :key="item.name" v-if="!item.hidden && !item.children">
          {{ item.name }}
        </el-breadcrumb-item>
      </template>
    </el-breadcrumb>
  </div>
</template>

<script setup lang="ts">
import { watchEffect } from "vue";
import { useRoute } from "vue-router";
import { ArrowDown } from "@element-plus/icons-vue";
import { directTo } from "@/router";

const route = useRoute() as any;

const router2path = (path: any) => {
  directTo(path);
};
</script>

<style lang="scss" scoped>
.el-dropdown-link {
  cursor: pointer;
  font-weight: normal;
  color: var(--el-color-primary);
}
</style>
