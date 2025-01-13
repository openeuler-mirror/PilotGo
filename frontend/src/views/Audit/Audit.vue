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
      :data="logs"
      title="审计日志"
      :showSelect="false"
      :total="total"
      v-model:page="page"
      :isExpand="true"
      v-model:expandData="expandLog"
    >
      <template v-slot:content>
        <el-table-column type="expand">
          <div style="width: 100%; padding-left: 10%">
            <el-table
              height="200"
              :data="logChildrens"
              style="width: 60%; border: 1px solid #dedfe0; border-radius: 4px"
            >
              <el-table-column prop="action" label="子日志名称" width="220" />
              <el-table-column label="子日志进度" width="200">
                <template #default="scope">
                  <el-progress style="width: 100%" type="line" :percentage="scope.row.status == 'OK' ? 100 : 0">
                    {{ scope.row.status == "OK" ? "1/1" : "0/1" }}
                  </el-progress></template
                >
              </el-table-column>
              <el-table-column align="center" prop="operation" label="详情">
                <template #default="scope">
                  <el-button size="small" @click="handleDetail(scope.row)"> 查看 </el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-table-column>
        <el-table-column align="center" prop="action" label="日志名称"> </el-table-column>
        <el-table-column align="center" prop="user_id" label="创建者"> </el-table-column>
        <el-table-column align="center" label="进度" style="display: flex; align-items: center">
          <template #default="scope">
            <el-progress
              v-if="!scope.row.Isempty"
              style="width: 100%"
              type="line"
              :percentage="scope.row.status == 'OK' ? 100 : 0"
            >
              {{ scope.row.status == "OK" ? "1/1" : "0/1" }}
            </el-progress>
            <el-progress
              v-else
              style="width: 100%"
              type="line"
              :percentage="
                scope.row.status.split(',')[2] === '1.00' || scope.row.status.split(',')[2] === '0.00'
                  ? 100
                  : scope.row.status.split(',')[2] * 100 || 0
              "
              :status="
                scope.row.status.split(',')[2] === '0.00'
                  ? 'exception'
                  : scope.row.status.split(',')[2] === '1.00'
                  ? 'success'
                  : 'warning'
              "
            >
              {{ scope.row.status.split(",")[0] + "/" + scope.row.status.split(",")[1] }}
            </el-progress>
          </template>
        </el-table-column>
        <el-table-column align="center" prop="CreatedAt" label="创建时间" sortable>
          <template #default="scope">
            <span>{{ scope.row.CreatedAt }}</span>
          </template>
        </el-table-column>
        <el-table-column align="center" prop="operation" label="详情">
          <template #default="scope">
            <el-button size="small" @click="handleDetail(scope.row)"> 查看 </el-button>
          </template>
        </el-table-column>
      </template>
    </PGTable>
    <el-drawer title="日志详情" v-model="showDetail" :before-close="handleClose" width="40%">
      <AuditDetail class="detail" :audit="auditData" v-if="showDetail"></AuditDetail>
    </el-drawer>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, watchEffect, watch } from "vue";
import { ElMessage } from "element-plus";
import PGTable from "@/components/PGTable.vue";

import { getLogs, getLogChildrens } from "@/request/audit";
import { RespCodeOK, type RespInterface } from "@/request/request";
import type { AuditItem } from "@/types/audit";
import AuditDetail from "./AuditDetail.vue";

const logs = ref<AuditItem[]>([]);
const total = ref(0);
const auditData = ref<AuditItem>();
const showDetail = ref(false);
const page = ref({ pageSize: 10, currentPage: 1 });
onMounted(() => {
  getPageLogs();
});
function getPageLogs() {
  getLogs({
    page: page.value.currentPage,
    size: page.value.pageSize,
  })
    .then((resp: RespInterface) => {
      if (resp.code === RespCodeOK) {
        total.value = resp.total!;
        logs.value = resp.data!;
      } else {
        ElMessage.error("failed to get audit logs: " + resp.msg);
      }
    })
    .catch((err: any) => {
      ElMessage.error("failed to get audit logs:" + err.msg);
    });
}

// 获取子日志
const expandLog = ref();
const logChildrens = ref<AuditItem[]>([]);
watchEffect(() => {
  if (expandLog.value && expandLog.value.log_uuid)
    getLogChildrens({ uuid: expandLog.value.log_uuid }).then((res: RespInterface) => {
      if (res.code === RespCodeOK) {
        logChildrens.value = res.data!;
      }
    });
});

// 获取详情
const handleDetail = (auditItem: AuditItem) => {
  auditData.value = auditItem;
  showDetail.value = true;
};
// 关闭弹窗
const handleClose = () => {
  showDetail.value = false;
  auditData.value = {} as AuditItem;
};

// 监听分页选项的修改
watch(
  () => page.value,
  (newV) => {
    if (newV) {
      getPageLogs();
    }
  },
  { deep: true }
);
</script>

<style lang="scss" scoped>
.container {
  height: 100%;
  width: 100%;

  .search {
    height: 100%;
    display: flex;
    flex-direction: row;
  }

  .el-button {
    margin-left: 5px;
  }
}
</style>
