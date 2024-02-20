<template>
  <div class="container">
    <PGTable :data="plugins" title="插件列表" :total="total" :currentPage="currentPage">
      <template v-slot:action>
        <el-button type="primary" @click="displayDialog = true">添加插件</el-button>
      </template>
      <template v-slot:content>
        <el-table-column align="center" prop="name" label="名称" width="150">
        </el-table-column>
        <el-table-column align="center" prop="version" label="版本" width="150">
        </el-table-column>
        <el-table-column align="center" prop="description" label="概述" show-overflow-tooltip>
        </el-table-column>
        <el-table-column align="center" prop="url" label="服务端地址" width="250">
        </el-table-column>
        <el-table-column align="center" prop="status" label="连接状态" width="150">
          <template #default="scope">
            <el-icon v-if="scope.row.status" style="color: green;">
              <SuccessFilled />
            </el-icon>
            <el-icon v-if="!scope.row.status" style="color: red;">
              <CircleCloseFilled />
            </el-icon>
            <span>{{ scope.row.status === true ? '连接' : '断开' }}</span>
          </template>
        </el-table-column>
        <el-table-column align="center" prop="lastheatbeat" label="上次心跳时间" width="250">
        </el-table-column>
        <el-table-column align="center" prop="enabled" label="状态" width="150">
          <template #default="scope">
            <span>{{ scope.row.enabled === 0 ? '已禁用' : '已启用' }}</span>
          </template>
        </el-table-column>
        <el-table-column align="center" label="操作" width="160">
          <template #default="scope">
            <el-button type="primary" plain name="default_all" @click="togglePluginState(scope.row)">
              {{ scope.row.enabled === 1 ? '禁用' : '启用' }}
            </el-button>
            <el-button type="danger" @click="onDeletePlugin(scope.row)">移除</el-button>
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
import { ref, onMounted } from "vue";
import { ElMessage } from 'element-plus';

import PGTable from "@/components/PGTable.vue";
import AddPlugin from "./components/AddPlugin.vue";

import { updatePlugins } from "./plugin";
import { getPluginsPaged, togglePlugin, deletePlugins } from "@/request/plugin";
import { RespCodeOK } from "@/request/request";
import type { Extention } from '@/types/plugin';
import { usePluginStore } from "@/stores/plugin";
const displayDialog = ref(false)

const plugins = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

onMounted(() => {
  updatePluginList()
})

function updatePluginList() {
  getPluginsPaged({
    page: currentPage.value,
    size: pageSize.value,
  }).then((resp: any) => {
    if (resp.code === RespCodeOK) {
      total.value = resp.total
      currentPage.value = resp.page
      pageSize.value = resp.size
      plugins.value = resp.data
    } else {
      ElMessage.error("failed to get plugins: " + resp.msg)
    }
  }).catch((err: any) => {
    ElMessage.error("failed to get plugins:" + err.msg)
  })
}

function togglePluginState(item: any) {
  let targetEnabled = item.enabled === 1 ? 0 : 1
  togglePlugin({ uuid: item.uuid, enable: targetEnabled }).then((res: any) => {
    if (res.code === RespCodeOK) {
      ElMessage.success(res.msg);
      // 更新插件列表
      updatePluginList();
      // 更新页面插件路由、sidebar等
      updatePlugins();
      // 删除插件tagview、增删全局扩展点
      let pluginExt = res.data.filter((item: Extention) => item.type === 'machine');
      if (targetEnabled === 0) {
        clearTagview(item);
        usePluginStore().delExtention(pluginExt)
      } else {
        usePluginStore().addExtention(pluginExt)
      }
    } else {
      ElMessage.error(res.msg);
    }
  })
}

function onDeletePlugin(item: any) {
  deletePlugins({ UUID: item.uuid }).then((res: any) => {
    if (res.code === RespCodeOK) {
      ElMessage.success(res.msg);
      // 更新插件列表
      updatePluginList();
      // 更新页面插件路由、sidebar等
      updatePlugins();
      // 删除插件tagview
      clearTagview(item);
    } else {
      ElMessage.error(res.msg);
    }
  })
}

import { tagviewStore } from '@/stores/tagview';
function clearTagview(item: any) {
  for (let i = 0; i < tagviewStore().taginfos.length; i++) {
    if (tagviewStore().taginfos[i].path === "/plugin-" + item.name) {
      tagviewStore().taginfos.splice(i, 1)
    }
  }
}




</script>

<style lang="scss" scoped>
.container {
  height: 100%;
  width: 100%;
}
</style>