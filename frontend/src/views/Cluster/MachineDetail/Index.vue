<!--
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
-->
<template>
  <div class="tab_content">
    <el-tabs class="tab" v-model="activePane">
      <el-tab-pane v-for="tab in tabs" :key="tab.name" :name="tab.name" :label="tab.label">
        <component :is="tab.component" v-if="activePane == tab.name" style="padding: 30px;">
        </component>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, markRaw, ref, type Component } from "vue";

import Base from "./Base.vue";
import User from "./User.vue";
import Service from "./Service.vue";
import Network from "./Network.vue";
import Sysctl from "./Sysctl.vue";
import Package from "./Package.vue";
import Terminal from "./Terminal.vue";

const activePane = ref("base")
interface Tab {
  name: string;
  label: string;
  component: Component;
}
let tabs = ref<Tab[]>([]);
onMounted(() => {
  tabs.value = [
    { name: 'base', label: '机器信息', component: markRaw(Base) },
    { name: 'user', label: '用户信息', component: markRaw(User) },
    { name: 'service', label: '服务信息', component: markRaw(Service) },
    { name: 'network', label: '网络配置', component: markRaw(Network) },
    { name: 'systl', label: '内核参数', component: markRaw(Sysctl) },
    { name: 'package', label: '软件包', component: markRaw(Package) },
    // { name: 'terminal', label: '远程终端', component: markRaw(Terminal) },
  ]
})

</script>

<style lang="scss" scoped>
.tab_content {
  width: 100%;
  height: 100%;

  .tab {
    width: 100%;
    height: 100%;
    display: flex;

    :deep(.el-tabs__content) {
      flex: 1;

      .el-tab-pane {
        height: 100%;
      }
    }
  }
}
</style>