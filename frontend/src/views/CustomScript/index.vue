<!--
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
-->
<template>
  <div class="container">
    <PGTable
      ref="refTable"
      id="exportTab"
      :data="scripts"
      title="自定义脚本"
      :showSelect="false"
      :total="total"
      v-model:page="page"
    >
      <template #action>
        <div class="search">
          <auth-button auth="button/update_script_blacklist" @click="handleBlackList"> 黑名单</auth-button>
          <el-button @click="addScript">新增</el-button>
        </div>
      </template>

      <template #content>
        <el-table-column align="center" prop="name" label="脚本名称"></el-table-column>
        <el-table-column align="center" prop="content" label="脚本内容">
          <template #default="{ row }">
            <el-popover placement="right" :width="500" trigger="hover" effect="dark" popper-class="custom-popover">
              <template #reference>
                <el-link type="primary" :underline="false">
                  {{ row.content.slice(0, 10) + "..." }}
                </el-link>
              </template>
              <pre style="color: #67c23a">{{ row.content }}</pre>
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column align="center" prop="UpdatedAt" label="更新时间"> </el-table-column>
        <el-table-column align="center" prop="description" label="描述"> </el-table-column>
        <el-table-column align="center" label="查看">
          <template #default="{ row }">
            <el-popover placement="right" :width="800" trigger="click">
              <template #reference>
                <el-button link size="small" type="warning" @click="handleHistory(row.id)">历史版本</el-button>
              </template>
              <!-- 列表table -->
              <version-list
                :ref="`versionRef${row.id}`"
                :historyList="historyList"
                v-if="showHistory === `history${row.id}`"
              />
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column align="center" label="操作" fixed="right" class="operate">
          <template #default="{ row }">
            <el-button size="small" @click="updateScript(row)">编辑</el-button>
            <auth-button size="small" auth="button/run_script" @click="runScript(row)">执行</auth-button>

            <el-popconfirm
              title="确定删除此脚本？"
              confirm-button-text="确定"
              cancel-button-text="取消"
              @confirm="onDeleteScript(row)"
              width="250"
            >
              <template #reference>
                <el-button type="danger" size="small">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </template>
    </PGTable>

    <el-dialog :title="title" v-model="display" :width="dialogWidth" destroy-on-close @close="closeDialog()">
      <AddScript
        v-if="displayDialog === 'add' || displayDialog === 'update'"
        :type="displayDialog"
        :scriptContent="scriptContent"
        @updateScript="getScriptList"
        @close="closeDialog"
      />
      <RunScript v-if="displayDialog === 'run'" :scriptContent="scriptContent" @close="closeDialog" />
      <BlackList v-if="displayDialog === 'blacklist'" @close="closeDialog" />
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, watch } from "vue";
import { ElMessage } from "element-plus";
import PGTable from "@/components/PGTable.vue";
import AddScript from "./components/addScript.vue";
import RunScript from "./components/runScript.vue";
import BlackList from "./components/BlackList.vue";
import AuthButton from "@/components/AuthButton.vue";
import { getScripts, deleteScript, getScriptHistorys } from "@/request/script";
import { RespCodeOK } from "@/request/request";

import VersionList from "./components/VersionList.vue";
const refTable = ref();
const scripts = ref([]);
const total = ref(0);
const page = ref({ pageSize: 10, currentPage: 1 });
onMounted(() => {
  getScriptList();
});

const getScriptList = () => {
  closeDialog();
  getScripts({
    page: page.value.currentPage,
    size: page.value.pageSize,
  })
    .then((resp: any) => {
      if (resp.code === RespCodeOK) {
        total.value = resp.total;
        scripts.value = resp.data;
      } else {
        ElMessage.error("failed to get scripts info: " + resp.msg);
      }
    })
    .catch((err: any) => {
      ElMessage.error("failed to get scripts info:" + err.msg);
    });
};

const dialogWidth = ref<number | string>("560px");
const display = ref(false);
const displayDialog = ref("");
const title = ref("");

const closeDialog = () => {
  display.value = false;
  title.value = "";
  displayDialog.value = "";
  scriptContent.value = {};
};
// 查看历史版本列表
const showHistory = ref("");
const historyList = ref<any[]>();
const handleHistory = (id: number) => {
  getScriptHistorys({ script_id: id })
    .then((resp: any) => {
      if (resp.code === RespCodeOK) {
        historyList.value = resp.data;
        showHistory.value = "history" + id;
        ElMessage.success("success to get history list:" + resp.msg);
      } else {
        ElMessage.error("failed to history list:" + resp.msg);
      }
    })
    .catch((err: any) => {
      ElMessage.error("failed to history lists:" + err.msg);
    });
};

// 新增脚本
const addScript = () => {
  title.value = "添加自定义脚本";
  displayDialog.value = "add";
  display.value = true;
  scriptContent.value = {};
};

// 执行脚本
const runScript = (row: any) => {
  title.value = "执行脚本";
  displayDialog.value = "run";
  display.value = true;
  dialogWidth.value = "720px";
  scriptContent.value = row;
};

// 删除脚本
const onDeleteScript = (row: any) => {
  deleteScript({
    script_id: row.id,
  })
    .then((resp: any) => {
      if (resp.code === RespCodeOK) {
        getScriptList();
        ElMessage.success("success to delete script:" + resp.msg);
      } else {
        ElMessage.error("failed to delete script:" + resp.msg);
      }
    })
    .catch((err: any) => {
      ElMessage.error("failed to delete script:" + err.msg);
    });
};

// 编辑脚本
const scriptContent = ref();
const updateScript = (row: any) => {
  scriptContent.value = row;
  title.value = "编辑脚本";
  displayDialog.value = "update";
  display.value = true;
};

// 更新黑名单
const handleBlackList = () => {
  title.value = "更新黑名单";
  displayDialog.value = "blacklist";
  display.value = true;
  dialogWidth.value = "720px";
};

// 监听分页选项的修改
watch(
  () => page.value,
  (newV) => {
    if (newV) {
      getScriptList();
    }
  },
  { deep: true }
);
</script>

<style lang="scss" scoped>
.container {
  height: 100%;
  width: 100%;
  :deep(.el-popper.is-light) {
    background-color: #303133 !important;
  }

  .search {
    height: 100%;
    display: flex;
    flex-direction: row;
    align-items: center;
    &_input {
      width: 300px;
    }
  }

  .el-button {
    margin-left: 5px;
  }
}
</style>
