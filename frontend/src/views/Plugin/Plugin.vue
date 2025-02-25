<!--
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
-->
<template>
  <div class="container">
    <PGTable :data="plugins" title="插件列表" :total="total" v-model:page="page">
      <template v-slot:action>
        <auth-button auth="button/plugin_operate" @click="displayDialog = true">添加插件</auth-button>
      </template>
      <template v-slot:content>
        <el-table-column align="center" prop="custom_name" label="插件名称"> </el-table-column>
        <el-table-column align="center" prop="name" label="插件类型"> </el-table-column>
        <el-table-column align="center" prop="version" label="版本"> </el-table-column>
        <el-table-column align="center" prop="description" label="概述" show-overflow-tooltip> </el-table-column>
        <el-table-column align="center" prop="url" label="服务端地址"> </el-table-column>
        <el-table-column align="center" prop="status" label="连接状态">
          <template #default="scope">
            <el-icon v-if="scope.row.status" style="color: green">
              <SuccessFilled />
            </el-icon>
            <el-icon v-if="!scope.row.status" style="color: red">
              <CircleCloseFilled />
            </el-icon>
            <span>{{ scope.row.status === true ? "连接" : "断开" }}</span>
          </template>
        </el-table-column>
        <el-table-column align="center" prop="lastheatbeat" label="上次心跳时间"> </el-table-column>
        <el-table-column align="center" prop="enabled" label="状态">
          <template #default="scope">
            <span>{{ scope.row.enabled === 0 ? "已禁用" : "已启用" }}</span>
          </template>
        </el-table-column>
        <el-table-column align="center" label="操作">
          <template #default="scope">
            <auth-button size="small" plain auth="button/plugin_operate" @click="togglePluginState(scope.row)">
              {{ scope.row.enabled === 1 ? "禁用" : "启用" }}
            </auth-button>
            <auth-button size="small" type="danger" auth="button/plugin_operate" @click="onDeletePlugin(scope.row)"
              >移除</auth-button
            >
          </template>
        </el-table-column>
      </template>
    </PGTable>
    <el-dialog title="添加插件" v-model="displayDialog" width="560px">
      <AddPlugin @pluginUpdated="updatePluginList" @close="displayDialog = false" />
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, watch } from "vue";
import { ElMessage } from "element-plus";

import PGTable from "@/components/PGTable.vue";
import AddPlugin from "./components/AddPlugin.vue";
import AuthButton from "@/components/AuthButton.vue";

import { updatePlugins } from "./plugin";
import { getPluginsPaged, togglePlugin, deletePlugins } from "@/request/plugin";
import { RespCodeOK } from "@/request/request";

const displayDialog = ref(false);

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
  let targetEnabled = item.enabled === 1 ? 0 : 1;
  togglePlugin({ uuid: item.uuid, enable: targetEnabled }).then((res: any) => {
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

function onDeletePlugin(item: any) {
  deletePlugins({ UUID: item.uuid }).then((res: any) => {
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
