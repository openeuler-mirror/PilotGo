<!--
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
-->
<template>
  <div class="container">
    <PGTable :data="plugins" title="插件服务列表" :total="total" v-model:page="page">
      <template v-slot:action> </template>
      <template v-slot:content>
        <el-table-column align="center" prop="serviceName" label="服务名称"> </el-table-column>
        <el-table-column align="center" prop="version" label="版本"> </el-table-column>
        <el-table-column align="center" prop="address" label="服务端地址">
          <template #default="scope">
            <span>{{ scope.row.address + ":" + scope.row.port }}</span>
          </template>
        </el-table-column>

        <el-table-column align="center" prop="enabled" label="状态">
          <template #default="scope">
            <span>{{ scope.row.status === false ? "已禁用" : "已启用" }}</span>
          </template>
        </el-table-column>
        <el-table-column align="center" label="操作">
          <template #default="scope">
            <auth-button size="small" plain auth="button/plugin_operate" @click="togglePluginState(scope.row)">
              {{ scope.row.status === true ? "禁用" : "启用" }}
            </auth-button>
          </template>
        </el-table-column>
      </template>
    </PGTable>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, watch } from "vue";
import { ElMessage } from "element-plus";

import PGTable from "@/components/PGTable.vue";
import AuthButton from "@/components/AuthButton.vue";

import { updatePlugins } from "./plugin";
import { getPluginsPaged, togglePlugin } from "@/request/plugin";
import { RespCodeOK } from "@/request/request";

const plugins = ref([]);
const total = ref(0);
const page = ref({ pageSize: 10, currentPage: 1 });

onMounted(() => {
  updatePluginList();
});

function updatePluginList() {
  getPluginsPaged({
    page: page.value.currentPage,
    size: page.value.pageSize,
  })
    .then((resp: any) => {
      if (resp.code === RespCodeOK) {
        total.value = resp.total;
        plugins.value = resp.data;
      } else {
        ElMessage.error("failed to get plugins: " + resp.msg);
      }
    })
    .catch((err: any) => {
      ElMessage.error("failed to get plugins:" + err.msg);
    });
}

// 监听分页选项的修改
watch(
  () => page.value,
  (newV) => {
    if (newV) {
      updatePluginList();
    }
  },
  { deep: true }
);

function togglePluginState(item: any) {
  togglePlugin({ serviceName: item.serviceName, enable: !item.status }).then((res: any) => {
    if (res.code === RespCodeOK) {
      ElMessage.success(res.msg);
      // 更新插件列表
      updatePluginList();
      // 更新页面插件路由、sidebar等
      updatePlugins();
    } else {
      ElMessage.error(res.msg);
    }
  });
}
</script>

<style lang="scss" scoped>
.container {
  height: 100%;
  width: 100%;
}
</style>
