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
      :data="batches"
      title="批次列表"
      :showSelect="true"
      :total="total"
      v-model:page="page"
      v-model:selectedData="selectedBatches"
    >
      <template v-slot:action>
        <el-dropdown>
          <el-button>
            操作<el-icon class="el-icon--right"><arrow-down /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item>
                <auth-button auth="button/batch_delete" link @click="batchDelete"> 删除 </auth-button>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </template>
      <template v-slot:content>
        <el-table-column align="center" label="批次名称">
          <template #default="{ row }">
            <el-link type="primary" @click="directTo(`/cluster/batch/detail/${row.ID}`)">
              {{ row.name }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column align="center" prop="manager" label="创建者"> </el-table-column>
        <el-table-column align="center" prop="DepartName" label="部门"> </el-table-column>
        <el-table-column align="center" prop="CreatedAt" label="创建时间" sortable> </el-table-column>
        <el-table-column align="center" prop="description" label="备注"> </el-table-column>
        <el-table-column align="center" prop="operation" label="操作">
          <template #default="scope">
            <auth-button auth="button/batch_update" @click="onEditBatch(scope.row.ID)"> 编辑 </auth-button>
          </template>
        </el-table-column>
      </template>
    </PGTable>

    <el-dialog title="编辑批次" v-model="showChangeBatchDialog">
      <UpdateBatch :batchID="updateBatchID" @batchUpdated="updateBatchInfo" @close="showChangeBatchDialog = false" />
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, onActivated, watch } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import AuthButton from "@/components/AuthButton.vue";
import PGTable from "@/components/PGTable.vue";
import UpdateBatch from "./components/UpdateBatch.vue";
import type { BatchItem } from "@/types/batch";
import { RespCodeOK } from "@/request/request";
import { getBatches, deleteBatch } from "@/request/batch";
import { directTo } from "@/router";

const showChangeBatchDialog = ref(false);
const updateBatchID = ref(0);

function onEditBatch(id: number) {
  updateBatchID.value = id;
  showChangeBatchDialog.value = true;
}

const batches = ref<BatchItem[]>([]);
const total = ref(0);
const page = ref({ pageSize: 10, currentPage: 1 });

onActivated(() => {
  updateBatchInfo();
});

onMounted(() => {
  updateBatchInfo();
});

function updateBatchInfo() {
  getBatches({
    page: page.value.currentPage,
    size: page.value.pageSize,
  })
    .then((resp: any) => {
      if (resp.code === RespCodeOK) {
        total.value = resp.total;
        batches.value = resp.data;
      } else {
        ElMessage.error("failed to get batch info: " + resp.msg);
      }
    })
    .catch((err: any) => {
      ElMessage.error("failed to get batch info:" + err.msg);
    });
}

const selectedBatches = ref<BatchItem[]>();

function batchDelete() {
  ElMessageBox.confirm("确定要删除该批次？", "正在删除批次", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  })
    .then(() => {
      // TODO: fix proxy object problem
      let batch_ids: number[] = [];
      let batch_names: string[] = [];
      selectedBatches.value!.forEach((item: any) => {
        batch_ids.push(item.ID);
        batch_names.push(item.name);
      });

      deleteBatch({ BatchID: batch_ids, Batches: batch_names })
        .then((resp: any) => {
          if (resp.code === RespCodeOK) {
            updateBatchInfo();
            ElMessage.success("批次删除成功");
          } else {
            ElMessage.error("failed to delete batch: " + resp.msg);
          }
        })
        .catch((err) => {
          ElMessage.error("failed to delete batch: " + err.msg);
        });
    })
    .catch(() => {
      // 取消删除批次
    });
}

// 监听分页选项的修改
watch(
  () => page.value,
  (newV) => {
    if (newV) {
      updateBatchInfo();
    }
  },
  { deep: true }
);
</script>

<style lang="scss" scoped>
.container {
  height: 100%;
  width: 100%;
}

a {
  text-decoration: none;
}
</style>
